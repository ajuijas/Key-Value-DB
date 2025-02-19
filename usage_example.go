package main

import (
	"fmt"
	"redis/redisclient"
)

// Start the server using the binary
// go build -o my-redis redis/main.go
// ./my-redis -H localhost -p 8080


func main() {
	client := redisclient.NewClient("localhost", "8080")
	fmt.Println(client.Set("key", "value"))
	fmt.Println(client.Get("key"))
	fmt.Println(client.Del([]string{"key", "key2"}))
	fmt.Println(client.Incr("key"))
	fmt.Println(client.Incrby("key", 3))
	fmt.Println(client.Close())
}