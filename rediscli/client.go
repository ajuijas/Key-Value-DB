package rediscli

import (
	"net"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// TODO: Clear the below thoughts after implementing cli

// Currently, I am using telnet client to send commands to my server
// This is enough for now, but I need more control over the client when I and implementing quorum based replication

// For now I will create a simple client which will create a connection, send and receive data from the server


type Server struct {
	id  uuid.UUID
	conn *net.Conn
}

type Client struct {
	servers []Server
}

func (c *Client) Set (key string, value string) string {
	// I will send the command to all the servers
	// For now I will wait for the response from one server only
	// TODO: Implement quorum based replication

	return sendDBCommand("set "+key+" "+value, c.servers, c.servers[0])
}

func (c *Client) Get (key string) string {
	// I will send the command to all the servers
	// For now I will wait for the response from one server only
	// TODO: Implement quorum based replication
	return sendDBCommand("get "+key, c.servers, c.servers[0])
}

func (c *Client) Del (keys []string) string {
	// TODO: test with keys having spaces
	return sendDBCommand("del "+strings.Join(keys, " "), c.servers, c.servers[0])
}

func (c *Client) Incr (key string) string {
	return sendDBCommand("incr "+key, c.servers, c.servers[0])
}

func (c *Client) Incrby (key string, value int) string {
	return sendDBCommand("incrby "+key+" "+strconv.Itoa(value), c.servers, c.servers[0])
}

func (c *Client) Close () string {
	return sendDBCommand("exit", c.servers, c.servers[0])
}

func sendDBCommand(command string, servers []Server, readServer Server) string {
	// I will send command to all the servers
	// Currently waiting for the response from one server only
	for _, server := range servers {
		conn := *server.conn
		_, err := conn.Write([]byte(command + "\n"))
		if err != nil {
			panic(err)
		}
	}

	// Read the response from the server
	buffer := make([]byte, 4096)
	readConn := *readServer.conn
	n, err := readConn.Read(buffer)
	if err != nil {
		panic(err)
	}
	return string(buffer[:n])
}

func createConnection (host string, port string) *net.Conn {
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		panic(err)
	}
	return &conn
}	

func NewClient(host, port string) *Client {
		servers := []Server{
			{	
				id : uuid.New(),  // I should get this id from server.
				conn: createConnection(host, port),
			},
		}
	return &Client{
		servers: servers,
	}
}

