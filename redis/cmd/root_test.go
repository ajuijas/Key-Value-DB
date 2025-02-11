package cmd

import (
	"net"
	"strconv"
	"testing"
	"time"
)

func Test_redis_commands(t *testing.T) {
	tests := []struct {
		command string
		expacted string
	}{
		{"", "\n"},
		{"SET key value1", "OK\n"}, // TODO: test case for lower letter commands
		{"SET key value", "OK\n"},
		{"GET key", "\"value\"\n"},

		{"DEL key", "(integer) 1\n"},

		{"SET key value", "OK\n"},
		{"SET key1 value", "OK\n"},
		{"DEL key key1 key2", "(integer) 2\n"},

		{"GET key", "(nil)\n"},

		{"SET key value invalid", "(error) ERR syntax error\n"},
		{"SET key", "(error) ERR wrong number of arguments for 'set' command\n"},
		{"GET key key2", "(error) ERR wrong number of arguments for 'get' command\n"},

		{"HI abc defg", "(error) ERR unknown command 'HI', with args beginning with: 'abc' 'defg'\n"},

		{"INCR key", "(integer) 1\n"},
		{"INCR key", "(integer) 2\n"},
		{"INCR key", "(integer) 3\n"},
	
		{"INCRBY key 3", "(integer) 6\n"},
		{"INCRBY key -2", "(integer) 4\n"},
		{"INCRBY key3 -2", "(integer) -2\n"},

		{"INCR key key1", "(error) ERR wrong number of arguments for 'incr' command\n"},
		{"INCRBY key", "(error) ERR wrong number of arguments for 'incrby' command\n"},
		{"INCRBY key 5 key2", "(error) ERR wrong number of arguments for 'incrby' command\n"},
		{"INCRBY key value", "(error) ERR value is not an integer or out of range\n"},

		{"SET key notanint", "OK\n"},
		{"INCR key", "(error) ERR value is not an integer or out of range\n"},
		{"INCRBY key 1", "(error) ERR value is not an integer or out of range\n"},

		{"MULTI", "OK\n"},
		{"INCR foo", "QUEUED\n"},
		{"SET bar 1", "QUEUED\n"},
		{"EXEC", "1) (integer) 1\n2) OK\n"},

		{"MULTI", "OK\n"},
		{"INCR foo", "QUEUED\n"},
		{"SET key1 1", "QUEUED\n"},
		{"DISCARD", "OK\n"},
		{"GET key1", "(nil)\n"},
		// TODO: Add testcases where error occured while using multi ops

		{"exit", "\n"},
	}

	host, port := "localhost", "8081"

	// start the server 
	rootCmd.SetArgs([]string{"-H", host, "-p", port})
	go rootCmd.Execute()

	conn, err := net.DialTimeout("tcp", host + ":" + port, 5*time.Second)
	if err!=nil {
		t.Fatalf("Unable to connect to test server")
	}
	defer conn.Close()

	for _, test := range tests{

		request := test.command + "\n"
		_, err = conn.Write([]byte(request))
		if err!=nil {
			t.Fatalf("Error while write to connection")
		}

		buffer := make([]byte, 4096)

		n, err := conn.Read(buffer)

		if err!=nil {
			t.Fatalf("Error reading from connection")
		}

		got := string(buffer[:n])

		if got!=test.expacted {
			t.Errorf("Expected <<%v>> Got <<%v>>", test.expacted, got)
		}
	}
}

func Test_atomic_operations(t *testing.T) {
	// Testing the incr operations are atomic or not.
	// I will create 5 clients, from each client I will increase value of a key 100 times.
	/* trunk-ignore(git-diff-check/error) */
	// The final value of the key will be 100 if all the operations where happend atomically.
	// Now I have no idea if this test will work or not. Lets see it

	host, port := "localhost", "8082"

	// start the server
	rootCmd.SetArgs([]string{"-H", host, "-p", port})
	go rootCmd.Execute()

	conn1, _ := net.DialTimeout("tcp", host + ":" + port, 5*time.Second)
	conn2, _ := net.DialTimeout("tcp", host + ":" + port, 5*time.Second)
	conn3, _ := net.DialTimeout("tcp", host + ":" + port, 5*time.Second)
	conn4, _ := net.DialTimeout("tcp", host + ":" + port, 5*time.Second)
	conn5, _ := net.DialTimeout("tcp", host + ":" + port, 5*time.Second)

	defer conn1.Close()
	defer conn2.Close()
	defer conn3.Close()
	defer conn4.Close()
	defer conn5.Close()

	connections := []net.Conn{conn1, conn2, conn3, conn4, conn5}
	key := "test_atomic"
	val := 0
	count := 500
	for {
		if val >= count {
			break
		}
		for _, connection := range connections{

			if val >= count {
				break
			}
			_, _ = connection.Write([]byte("incr " + key + "\n"))
			_, _ = connection.Write([]byte("incrby " + key + " 1\n"))
			val += 2
		}
	}

		time.Sleep(2 * time.Second) // Assumes that all the operations are completed after 3 seconds

		conn6, _ := net.DialTimeout("tcp", host + ":" + port, 5*time.Second)
		defer conn6.Close()

		_, err := conn6.Write([]byte("get " + key + "\n"))
		if err!=nil {
			t.Fatalf("Error while write to connection")
		}

		buffer := make([]byte, 4096)

		n, err := conn6.Read(buffer)

		if err!=nil {
			t.Fatalf("Error reading from connection")
		}

		got := string(buffer[:n])

		if got!=strconv.Itoa(count) {
			t.Errorf("Expected <<%v>> Got <<%v>>", strconv.Itoa(count), got)
		}
}

func Test_atominc_multi_ops(t *testing.T) {
	// Testing the multi operations are atomic or not.
	// I will incr the value of a key 1000 times and try to read the value from another client.
	// I should be able to read the final value only.
	host, port := "localhost", "8083"

	// start the server
	rootCmd.SetArgs([]string{"-H", host, "-p", port})
	go rootCmd.Execute()

	conn1, _ := net.DialTimeout("tcp", host + ":" + port, 5*time.Second)
	conn2, _ := net.DialTimeout("tcp", host + ":" + port, 5*time.Second)

	defer conn1.Close()
	defer conn2.Close()

	key := "test_atomic"

	_, _ = conn1.Write([]byte("multi\n"))
	_, _ = conn1.Write([]byte("incr " + key + "\n"))

	_, err := conn2.Write([]byte("get " + key + "\n"))
	if err!=nil {
		t.Fatalf("Error while write to connection")
	}

	buffer := make([]byte, 4096)

	n, err := conn2.Read(buffer)

	if err!=nil {
		t.Fatalf("Error reading from connection")
	}

	got := string(buffer[:n])

	if got!="(nil)\n" {
		t.Errorf("Expected <<%v>> Got <<%v>>", "(nil)\n", got)
	}

	_, _ = conn1.Write([]byte("exec\n"))

	_, err = conn2.Write([]byte("get " + key + "\n"))
	if err!=nil {
		t.Fatalf("Error while write to connection")
	}

	buffer = make([]byte, 4096)

	n, err = conn2.Read(buffer)

	if err!=nil {
		t.Fatalf("Error reading from connection")
	}

	got = string(buffer[:n])

	if got!="1" {
		t.Errorf("Expected <<%v>> Got <<%v>>", "1", got)
	}
}