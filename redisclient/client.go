package redisclient

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
	n, m int
}

func (c *Client) Set (key string, value string) string {
	// I will send the command to all the servers
	// For now I will wait for the response from one server only
	// TODO: Implement quorum based replication

	return c.sendDBCommand("set "+key+" "+value, "w")
}

func (c *Client) Get (key string) string {
	// I will send the command to all the servers
	// For now I will wait for the response from one server only
	// TODO: Implement quorum based replication
	return c.sendDBCommand("get "+key, "r")
}

func (c *Client) Del (keys []string) string {
	// TODO: test with keys having spaces
	return c.sendDBCommand("del "+strings.Join(keys, " "), "w")
}

func (c *Client) Incr (key string) string {
	return c.sendDBCommand("incr "+key, "r")
}

func (c *Client) Incrby (key string, value int) string {
	return c.sendDBCommand("incrby "+key+" "+strconv.Itoa(value), "r")
}

func (c *Client) Close () string {
	return c.sendDBCommand("exit", "w")
}

func (c *Client) sendDBCommand (command string, readOrWrite string) string {
	// I will send command to all the servers
	// Currently waiting for the response from one server only
	for _, server := range c.servers {
		conn := *server.conn
		_, err := conn.Write([]byte(command + "\n"))
		if err != nil {
			panic(err)
		}
	}
	confirmationCount := 1

	if readOrWrite == "w" {
		confirmationCount = c.n
	}else if readOrWrite == "r" {
		confirmationCount = c.m
	}

	responses := make([]string, confirmationCount)

	for i, server := range c.servers {
		buffer := make([]byte, 4096)
		readConn := *server.conn
		if i >= confirmationCount {
			go readConn.Read(buffer)
			// Not waiting for responses
			continue
		}else {
			n, err := readConn.Read(buffer) // I should't be waiting here
			if err != nil {
				panic(err)
			}
			got := string(buffer[:n])
			responses[i] = got
		}
	}
	// Ideally all the responses should be the same.
	// TODO: Add a mechanism to sync the servers incase of different responses

	return responses[0]
}

func createConnection (hostUrl string) *net.Conn {
	conn, err := net.Dial("tcp", hostUrl)
	if err != nil {
		panic(err)
	}
	return &conn
}	

func NewClient(hostUrls []string, n, m int) *Client {
		servers := make([]Server, len(hostUrls))
		for i, hostUrl := range hostUrls {
			conn := createConnection(hostUrl)
			servers[i] = Server{
				id: uuid.New(),
				conn: conn,
			}
		}
	return &Client{
		servers: servers,
		n: n,
		m: m,
	}
}

