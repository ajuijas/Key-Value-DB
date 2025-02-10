package cmd

import (
	"bytes"
	"net"
	"os"
	"testing"
	"time"
)

func captureStandardOutput() (func () string) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	rOut, wOut, _ := os.Pipe()
	os.Stdout = wOut
	os.Stderr = wOut

	return func() string {
		wOut.Close()

		var buf bytes.Buffer
		_, _ = buf.ReadFrom(rOut)

		rOut.Close()

		os.Stdout = oldStdout
		os.Stderr = oldStderr

		return buf.String()
	}
}

func Test_redis_commands(t *testing.T) {
	tests := []struct {
		command string
		expacted string
	}{
		{"SET key value", "OK\n"},
		{"SET key value", "OK\n"},
		{"GET key", "\"value\"\n"},
		{"DEL key", "(integer) 1\n"},
		{"GET key", "nil\n"},
	}

	host, port, db := "localhost", "8080", "0"

	// start the server 
	rootCmd.SetArgs([]string{"-h", host, "-p", port, "-d", db})
	go rootCmd.Execute()

	for _, test := range tests{
		// captureStandardOutput := captureStandardOutput()

		conn, err := net.DialTimeout("tcp", host + ":" + port, 5*time.Second)
		if err!=nil {
			t.Fatalf("Unable to connect to test server")
		}
		defer conn.Close()

		request := test.command
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