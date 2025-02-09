package cmd

import (
	"bytes"
	"os"
	"testing"
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

func Test_cli_shell(t *testing.T){
	tests := []struct {
		args []string
		expectedResponse string
	}{
		{
			args: []string{},
			expectedResponse: "Clinet connected to localhost:8080. type exit to quit",
		},
		{
			args: []string{"-p", "8081"},
			expectedResponse: "Clinet connected to localhost:8081. type exit to quit",
		},
	}
	
	for _, test := range tests {
		captureStandardOutput := captureStandardOutput()
		rootCmd.SetArgs(test.args)
		_ = rootCmd.Execute()
		response := captureStandardOutput()
		if response != test.expectedResponse {
			t.Errorf("expected <<%s>>, got <<%s>>", test.expectedResponse, response)
		}
	}
}

func Test_commands_usage(t *testing.T) {
	tests := []struct {
		command string
		expectedResponse string
	}{
		{
			command: "set",
			expectedResponse: "(error) ERR wrong number of arguments for 'set' command",
		},
		{
			command: "set test_ket test_value",
			expectedResponse: "ok",
		},
		{
			command: "set test_key test_value invalid",
			expectedResponse: "(error) ERR syntax error",
		},
		{
			command: "get",
			expectedResponse: "(error) ERR wrong number of arguments for 'get' command",
		},
		{
			command: "get test_key_not_exist",
			expectedResponse: "(nil)",
		},
		{
			command: "get test_key",
			expectedResponse: "\"test_value\"",
		},
		{
			command: "get test_key invalid",
			expectedResponse: "(error) ERR wrong number of arguments for 'get' command",
		},
	}

	for _, test := range tests {
		captureStandardOutput := captureStandardOutput()

		oldStdin := os.Stdin
		r, w, _ := os.Pipe()
		os.Stdin = r

		defer func() {
			os.Stdin = oldStdin
			r.Close()
			w.Close()
		}()

		// Writing to the stdin pipe
		go func() {
			defer w.Close()
			_, _ = w.Write([]byte(test.command))
		}()
		// TODO: move this to a function if needed

		rootCmd.SetArgs([]string{})
		_ = rootCmd.Execute()
		response := captureStandardOutput()
		if response != test.expectedResponse {
			t.Errorf("expected <<%s>>, got <<%s>>", test.expectedResponse, response)
		}
	}
}