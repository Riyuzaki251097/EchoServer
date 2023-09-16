package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	go receiveMessages(conn)

	// Ask the user for their name
	fmt.Print("Enter your name: ")
	inputReader := bufio.NewReader(os.Stdin)
	name, _ := inputReader.ReadString('\n')
	name = strings.TrimSpace(name)

	// Read and send messages to the server
	for {
		fmt.Print("Send a message (or 'exit' to quit): ")
		message, _ := inputReader.ReadString('\n')
		message = strings.TrimSpace(message)

		if message == "exit" {
			break
		}

		// Format the message with the sender's name
		messageToSend := name + ": " + message

		_, err := conn.Write([]byte(messageToSend + "\n"))
		if err != nil {
			fmt.Println("Error sending message to the server:", err)
			break
		}
	}
}

func receiveMessages(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Server closed the connection.")
			} else {
				fmt.Println("Error reading message from the server:", err)
			}
			break
		}
		fmt.Print("Received message from server: ", message)
	}
}
