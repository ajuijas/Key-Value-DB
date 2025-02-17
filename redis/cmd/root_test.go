package cmd

import (
	"io"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_redis_commands(t *testing.T) {
	tests := []struct {
		command  string
		expacted string
	}{
		{"", "\n"},
		{"SET key \"value with space\"", "OK\n"},
		{"SET \"key with space\" \"value with space\"", "OK\n"},
		{"SET 'key with space' 'value with space'", "OK\n"},
		{"GET 'key with space'", "\"value with space\"\n"},
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
		{"EXEC", "\n"},

		{"MULTI", "OK\n"},
		{"INCR foo", "QUEUED\n"},
		{"SET bar 1", "QUEUED\n"},
		{"EXEC", "\n1) (integer) 1\n2) OK\n"},

		{"MULTI", "OK\n"},
		{"INCR foo", "QUEUED\n"},
		{"SET key1 1", "QUEUED\n"},
		{"DISCARD", "OK\n"},
		{"GET key1", "(nil)\n"},
		// TODO: Add testcases where error occured while using multi ops

		{"exit", "Bye!\n"},
	}

	host, port := "localhost", "8081"

	// start the server
	dbFile := uuid.New().String() + "/"
	_ = os.MkdirAll(dbFile, os.ModePerm)
	defer os.RemoveAll(dbFile)
	rootCmd.SetArgs([]string{"-H", host, "-p", port, "-s", dbFile})
	go rootCmd.Execute()

	conn, err := net.DialTimeout("tcp", host+":"+port, 5*time.Second)
	if err != nil {
		t.Fatalf("Unable to connect to test server")
	}
	defer conn.Close()

	for _, test := range tests {

		request := test.command + "\n"
		_, err = conn.Write([]byte(request))
		if err != nil {
			t.Fatalf("Error while write to connection")
		}

		buffer := make([]byte, 4096)

		n, err := conn.Read(buffer)

		if err != nil {
			t.Fatalf("Error reading from connection")
		}

		got := string(buffer[:n])

		if got != test.expacted {
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
	dbFile := uuid.New().String() + "/"
	_ = os.MkdirAll(dbFile, os.ModePerm)
	defer os.RemoveAll(dbFile)
	rootCmd.SetArgs([]string{"-H", host, "-p", port, "-s", dbFile})
	go rootCmd.Execute()

	conn1, _ := net.DialTimeout("tcp", host+":"+port, 5*time.Second)
	conn2, _ := net.DialTimeout("tcp", host+":"+port, 5*time.Second)
	conn3, _ := net.DialTimeout("tcp", host+":"+port, 5*time.Second)
	conn4, _ := net.DialTimeout("tcp", host+":"+port, 5*time.Second)
	conn5, _ := net.DialTimeout("tcp", host+":"+port, 5*time.Second)

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
		for _, connection := range connections {

			if val >= count {
				break
			}
			_, _ = connection.Write([]byte("incr " + key + "\n"))
			_, _ = connection.Write([]byte("incrby " + key + " 1\n"))
			val += 2
		}
	}

	time.Sleep(2 * time.Second) // Assumes that all the operations are completed after 3 seconds

	conn6, _ := net.DialTimeout("tcp", host+":"+port, 5*time.Second)
	defer conn6.Close()

	_, err := conn6.Write([]byte("get " + key + "\n"))
	if err != nil {
		t.Fatalf("Error while write to connection")
	}

	buffer := make([]byte, 4096)

	n, err := conn6.Read(buffer)

	if err != nil {
		t.Fatalf("Error reading from connection")
	}

	got := string(buffer[:n])

	if got != strconv.Itoa(count) {
		t.Errorf("Expected <<%v>> Got <<%v>>", strconv.Itoa(count), got)
	}
}

func sendDBCommand(cmd string, conn net.Conn) string {
	_, err := conn.Write([]byte(cmd + "\n"))
	if err != nil {
		panic("Error while write to connection")
	}

	buffer := make([]byte, 4096)

	n, err := conn.Read(buffer)

	if err != nil {
		panic("Error reading from connection")
	}
	got := string(buffer[:n])
	return got
}

func Test_atominc_multi_ops(t *testing.T) {
	// Testing the multi operations are atomic or not.
	// I will incr the value of a key 1000 times and try to read the value from another client.
	// I should be able to read the final value only.

	host, port := "localhost", "8083"

	// start the server
	dbFile := uuid.New().String() + "/"
	_ = os.MkdirAll(dbFile, os.ModePerm)
	defer os.RemoveAll(dbFile)
	rootCmd.SetArgs([]string{"-H", host, "-p", port, "-s", dbFile})
	go rootCmd.Execute()

	n := 10000

	conn1, _ := net.DialTimeout("tcp", host+":"+port, 5*time.Second)
	conn2, _ := net.DialTimeout("tcp", host+":"+port, 5*time.Second)
	conn3, _ := net.DialTimeout("tcp", host+":"+port, 5*time.Second)

	defer conn1.Close()
	defer conn2.Close()
	defer conn3.Close()

	sendDBCommand("multi", conn1)
	sendDBCommand("multi", conn2)

	for i := 0; i < n; i++ {
		_, _ = conn1.Write([]byte("incr key\n"))
		_, _ = conn2.Write([]byte("incr key\n"))
	}

	_, _ = conn1.Write([]byte("exec\n"))
	_, _ = conn2.Write([]byte("exec\n"))

	value := sendDBCommand("get key", conn3)

	if value != "(nil)\n" && value != "10000" && value != "20000"{
		// The test is failed when the 'got value' is not 10000 or 20000 or (nil), since the operations are atomic
		t.Errorf("Expected <<%v>> Got <<%v>>", "any of (10000, 20000, (nil))", value)
	}
}

func Test_rdbFile(t *testing.T) {

	host, port := "localhost", "8084"

	// start the server
	dbFile := uuid.New().String() + "/"
	_ = os.MkdirAll(dbFile, os.ModePerm)
	defer os.RemoveAll(dbFile)
	rootCmd.SetArgs([]string{"-H", host, "-p", port, "-s", dbFile})
	go rootCmd.Execute()

	conn, _ := net.DialTimeout("tcp", host+":"+port, 5*time.Second)
	defer conn.Close()

	_ = sendDBCommand("set key 55", conn)
	_ = sendDBCommand("get key", conn)
	_ = sendDBCommand("del key", conn)
	_ = sendDBCommand("set key 45", conn)
	_ = sendDBCommand("incr key", conn)
	_ = sendDBCommand("incrby key 5", conn)
	_ = sendDBCommand("exit", conn)

	time.Sleep(2 * time.Second)

	file, err := os.Open("./" + dbFile + "dump.rdb")
	if err != nil {
		t.Fatalf("Error while reading from file")
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("Error while reading from file")
	}

	// I am expecting the set, incr, incrby, and del commands to be in the file.
	// get, exit and multi commands are not expected in the file
	expectedContent := "set key 55\ndel key\nset key 45\nincr key\nincrby key 5\n"
	got := string(data)
	if got != expectedContent {
		t.Errorf("Expected <<%v>> Got <<%v>>", expectedContent, got)
	}
}


func Test_db_loaded_from_rdbFile(t *testing.T) {
	//TODO: This is a very basic test. I need to test more cases

	host, port := "localhost", "8085"

	// start the server
	dbFile := uuid.New().String() + "/"
	_ = os.MkdirAll(dbFile, os.ModePerm)
	defer os.RemoveAll(dbFile)

	// Write some data to the rdb file
	file, err := os.OpenFile(dbFile+"dump.rdb", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write([]byte("set key 55\ndel key\nset key 45\nincr key\nincrby key 5\n"))
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	rootCmd.SetArgs([]string{"-H", host, "-p", port, "-s", dbFile})
	go rootCmd.Execute()

	conn, _ := net.DialTimeout("tcp", host+":"+port, 5*time.Second)
	defer conn.Close()

	got := sendDBCommand("get key", conn)

	if got != "51" {
		t.Errorf("Expected <<51>> Got <<%v>>", got)
	}
}