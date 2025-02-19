package redisclient

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"testing"
	"time"
)


func runMockedServer (host, port string) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Mocked server listening to ", host, ":", port)
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleMockedRequest(conn)
	}
}

func handleMockedRequest (conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			conn.Close()
			return
		}
		_, _ = conn.Write([]byte(message))
	}
}

func Test_mockedServer(t *testing.T) {

	// In for this test I will test my client with a mocked tcp servers.
	// The servers will always return the same command back to the client.

	host, port := "localhost", "8085"
	go runMockedServer(host, port)

	conn, _ := net.DialTimeout("tcp", host+":"+port, 5*time.Second)
	defer conn.Close()

	client := NewClient(host, port)
	
	tests := []struct {
		funcCall, expected string
	}{
		{client.Set("key", "value"), "set key value\n"},
		{client.Get("key"), "get key\n"},
		{client.Del([]string{"key1", "key2"}), "del key1 key2\n"},
		{client.Incr("key"), "incr key\n"},
		{client.Incrby("key", 3), "incrby key 3\n"},
		{client.Close(), "exit\n"},
	}

	for _, test := range tests {
		if test.funcCall != test.expected{
			t.Errorf("Expected <<%v>> Got <<%v>>", test.expected, test.funcCall)
		}
	}
}
