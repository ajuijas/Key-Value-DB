# Simple Key-Value Database in Go

This is a simple Redis-like key-value database server implemented in Go. It allows you to store and retrieve data using a simple command-line interface through telnet.

## Features

- In-memory key-value storage
- Support for basic Redis-like commands
- Multiple client connections
- Data persistence using RDB file
- Transaction support using MULTI commands

## Getting Started

### Prerequisites

- Go 1.x or higher installed on your system
- Telnet client for connecting to the server

### Building the Binary

Alternatively, you can build a binary file for the server. Run the following command from the project root:

```bash
go build -o my-redis redis/main.go
```

Then, start the server using the binary:

```bash
./my-redis
```

### Running the Server

The server runs using the binary file. By default, it listens on:
- Host: localhost
- Port: 8080
- Storage path: /usr/local/var/db/my_redis/

You can override settings using flags:

```bash
./my-redis -H <host> -p <port> -s <storage-path>
```

### Connecting to the Server

Use telnet to connect to the server:
```bash
telnet localhost 8080
```

## Available Commands

1. `SET key value` - Store a key-value pair
   ```
   SET name john
   ```

2. `GET key` - Retrieve a value by key
   ```
   GET name
   ```

3. `DEL key [key ...]` - Delete one or more keys
   ```
   DEL name age
   ```

4. `INCR key` - Increment the integer value of a key by 1
   ```
   INCR counter
   ```

5. `INCRBY key increment` - Increment the integer value of a key by the given amount
   ```
   INCRBY counter 5
   ```

6. Transaction Commands:
   - `MULTI` - Start a transaction
   - `EXEC` - Execute all commands in the transaction
   - `DISCARD` - Discard the transaction

   Example:
   ```
   MULTI
   SET user1 john
   SET user2 jane
   EXEC
   ```

7. `exit` - Close the connection

## Data Persistence

The server automatically saves data to the specified storage path in an RDB file. This data is loaded when the server starts up again.

## Notes

- All commands are case-insensitive
- String values containing spaces should be wrapped in quotes
- The server supports multiple concurrent client connections

### Running Tests

To run tests for the project, execute the following command from the project root:

```bash
go test ./...
```
