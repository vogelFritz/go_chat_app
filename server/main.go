package main

import (
	"fmt"
	"net"
	"os"
)

const (
	SERVER_TYPE = "tcp"
	SERVER_HOST = "localhost"
	SERVER_PORT = "8080"
	ADDRESS     = SERVER_HOST + ":" + SERVER_PORT
)

func main() {
	server := startServer()
	defer server.Close()
	fmt.Println("Listening on " + ADDRESS)
	fmt.Println("Waiting for client...")
	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("client connected")
		go processClient(connection)
	}

}

func startServer() net.Listener {
	server, err := net.Listen(SERVER_TYPE, ADDRESS)
	if err != nil {
		fmt.Println("Error listening: ", err.Error())
		os.Exit(1)
	}
	return server
}

func processClient(connection net.Conn) {
	receivedMessage := readMessage(connection)
	for receivedMessage != "f" {
		connection.Write([]byte("Thanks! Got your message:" + receivedMessage))
		receivedMessage = readMessage(connection)
	}
	connection.Close()
}

func readMessage(connection net.Conn) string {
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
	}
	receivedMessage := string(buffer[:mLen])
	fmt.Println("Received: ", receivedMessage)
	return receivedMessage
}
