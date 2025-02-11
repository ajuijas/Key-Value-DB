package cmd

import (
	"net"
	"testing"
	"time"
)

func Test_redis_commands(t *testing.T) {
	tests := []struct {
		command string
		expacted string
	}{
		{"SET key value1", "OK\n"},
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