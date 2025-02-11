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
	reader := bufio.NewReader(client.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			client.conn.Close()
			return
		}
		cmd := strings.Fields(string(message))
		var msg string
		switch strings.ToLower(cmd[0]){  // TODO: Is there better method than switch for this?
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
	
		client.conn.Write([]byte(msg))
	}
}