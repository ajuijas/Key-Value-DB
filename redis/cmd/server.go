package cmd

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

// Server ...
type Server struct {
	host string
	port string
	storage *Storage
}

// Client ...
type Client struct {
	conn net.Conn
	reader *bufio.Reader
	storage *Storage
}

// Config ...
type Config struct {
	Host string
	Port string
}

// New ...
func New(config *Config) *Server {
	return &Server{
		host: config.Host,
		port: config.Port,
		storage: getStorage(),
	}
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
			conn: conn,
			storage: server.storage,
		}
		go client.handleRequest()
	}
}

func (client *Client) handleRequest() {
	client.reader = bufio.NewReader(client.conn)
	for {
		message, err := client.reader.ReadString('\n')
		if err != nil {
			client.conn.Close()
			return
		}

		cmd := strings.Fields(string(message))

		var msg string

		switch strings.ToLower(cmd[0]){
		case "multi" : msg = client.handleMulti(cmd)
		default : msg = client.executeCmd(cmd)
		}
		client.conn.Write([]byte(msg))
	}
}

func (client *Client) handleMulti (args []string) string {

	if len(args) != 1 {
		return "(error) ERR wrong number of arguments for 'multi' command"
	}

	client.conn.Write([]byte("OK\n"))
	var cmdList [][]string


	for {
		message, err := client.reader.ReadString('\n')

		if err != nil {
			client.conn.Close()
			return "" // TODO: correct error handling
		}

		cmd := strings.Fields(string(message))
		if strings.ToLower(cmd[0]) == "exec"{
			break
		}else if strings.ToLower(cmd[0]) == "discard"{
			return "OK\n"
		}else {
			cmdList = append(cmdList, cmd)
			client.conn.Write([]byte("QUEUED\n"))
		}
	}

	var msg string
	// TODO: make operations in the list atomic
	for i, cmd := range cmdList {
		resp := client.executeCmd(cmd)
		msg += fmt.Sprintf("%v) %v", i+1, resp)
	}
	return msg
}

func (client *Client) executeCmd (cmd []string) string {
		var msg string
		switch strings.ToLower(cmd[0]){  // TODO: Is there any better method than switch for this?
		case "set":
			msg = client.storage.set(cmd[1:])
		case "get":
			msg = client.storage.get(cmd[1:])
		case "del":
			msg = client.storage.del(cmd[1:])
		case "incr":
			msg = client.storage.incr(cmd[1:])
		case "incrby":
			msg = client.storage.incrby(cmd[1:])
		default:
			msg = "(error) ERR unknown command '" + cmd[0] + "', with args beginning with: '" + strings.Join(cmd[1:], "' '") + "'\n"
		}
		return msg
}