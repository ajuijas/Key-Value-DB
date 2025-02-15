package cmd

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
)

// Server ...
type Server struct {
	host    string
	port    string
	storage *Storage
}

// Client ...
type Client struct {
	conn    net.Conn
	reader  *bufio.Reader
	storage *Storage
	log *Logger
}

// Config ...
type Config struct {
	Host        string
	Port        string
	StorageFile string
}

// New ...
func New(config *Config) *Server {
	server := &Server{
		host:    config.Host,
		port:    config.Port,
		storage: getStorage(),
	}
	return server
}

// Run ...
func (server *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Server listening to %s:%s \n", server.host, server.port)
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		client := &Client{
			conn:    conn,
			storage: server.storage,
			log:     NewLogger(),
		}

		go client.handleRequest()
	}
}

func parseCommand(input string) []string {
	re := regexp.MustCompile(`'([^']*)'|"([^"]*)"|(\S+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	var result []string
	for _, match := range matches {
		if match[1] != "" {
			result = append(result, match[1])
		} else if match[2] != "" {
			result = append(result, match[2])
		} else {
			result = append(result, match[3])
		}
	}
	return result
}

func (client *Client) close () {
	client.conn.Close()
	client.log.Close()
}

func (client *Client) handleRequest() {
	client.reader = bufio.NewReader(client.conn)
	for {
		message, err := client.reader.ReadString('\n')
		if err != nil {
			client.close()
			return
		}

		cmd := parseCommand(string(message))

		if len(cmd) == 0 {
			_, _ = client.conn.Write([]byte("\n"))
			continue
		}

		var msg string

		switch strings.ToLower(cmd[0]) {
		case "exit":
			msg = "\n"
		case "multi":
			msg = client.handleMulti(cmd)
		default:
			msg = client.executeCmd(cmd, false)
		}
		_, _ = client.conn.Write([]byte(msg))
	}
}

func (client *Client) handleMulti(args []string) string {

	if len(args) != 1 {
		return "(error) ERR wrong number of arguments for 'multi' command"
	}
	_, _ = client.conn.Write([]byte("OK\n")) // TODO: Move client.conn.Write to a function with proper error handling
	var cmdList [][]string

	for {
		message, err := client.reader.ReadString('\n')
		if err != nil {
			client.conn.Close()
			return "\n" // TODO: correct error handling
		}

		cmd := strings.Fields(string(message))
		if strings.ToLower(cmd[0]) == "exec" {
			msg := client.executeMulti(cmdList)
			_, _ =client.conn.Write([]byte("\n"))
			return msg
		} else if strings.ToLower(cmd[0]) == "discard" {
			return "OK\n"
		} else {
			cmdList = append(cmdList, cmd)
			_, _ = client.conn.Write([]byte("QUEUED\n"))
		}
	}
}

func (client *Client) executeMulti(cmdList [][]string) string {

	client.storage.mutex.Lock()
	defer client.storage.mutex.Unlock()

	var msg string
	for i, cmd := range cmdList {
		resp := client.executeCmd(cmd, true)
		msg += fmt.Sprintf("%v) %v", i+1, resp)
	}
	return msg
}

func (client *Client) executeCmd(cmd []string, isMulti bool) (string) {
	var msg string

	if !isMulti { // The goroutine is already locked for multi.
		client.storage.mutex.Lock()
		defer client.storage.mutex.Unlock()
	}

	isValueChanged := false

	switch strings.ToLower(cmd[0]) { // TODO: Is there any better method than switch for this?
	case "set":
		isValueChanged, msg = client.storage.set(cmd[1:])
	case "get":
		_, msg = client.storage.get(cmd[1:])
	case "del":
		isValueChanged, msg = client.storage.del(cmd[1:])
	case "incr":
		isValueChanged, msg = client.storage.incr(cmd[1:])
	case "incrby":
		isValueChanged, msg = client.storage.incrby(cmd[1:])
	default:
		msg = "(error) ERR unknown command '" + cmd[0] + "', with args beginning with: '" + strings.Join(cmd[1:], "' '") + "'\n"
	}

	if isValueChanged {
		// log the cmd to backup file
		client.log.log.Println(strings.Join(cmd, " "))
		// fmt.Println(strings.Join(cmd, " "))
	}

	return msg
}
