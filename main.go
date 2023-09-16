package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"
)

var (
	clients   = make(map[net.Conn]struct{})
	clientsMu sync.Mutex
)

func main() {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		clientsMu.Lock()
		clients[conn] = struct{}{}
		clientsMu.Unlock()

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer func() {
		conn.Close()
		clientsMu.Lock()
		delete(clients, conn)
		clientsMu.Unlock()
	}()

	fmt.Printf("Client %s connected.\n", conn.RemoteAddr())

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Client %s disconnected.\n", conn.RemoteAddr())
				return
			}
			fmt.Printf("Error reading from client %s: %v\n", conn.RemoteAddr(), err)
			return
		}

		// Broadcast the message to all other clients
		broadcastMessage(message, conn)
	}
}

func broadcastMessage(message string, sender net.Conn) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for client := range clients {
		if client != sender {
			_, err := io.WriteString(client, message)
			if err != nil {
				fmt.Printf("Error broadcasting message to client %s: %v\n", client.RemoteAddr(), err)
				// Handle the error (e.g., remove the disconnected client).
			}
		}
	}
}
