package main

import (
	"fmt"
	"net"
	"os"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "8080"
	SERVER_TYPE = "tcp"
	ADDRESS     = SERVER_HOST + ":" + SERVER_PORT
)

func main() {
	connection := establishConnection()
	chatLoop(connection)
	defer connection.Close()
}

func chatLoop(connection net.Conn) {
	var message string
	for {
		fmt.Print("Write your message here (f to finish): ")
		fmt.Scanf("%v", message)
		if message == "f" {
			os.Exit(1)
		}
		sendMessage(message, connection)
	}
}

func sendMessage(message string, connection net.Conn) {
	_, err := connection.Write([]byte(message))
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
	}
	fmt.Println("Received: ", string(buffer[:mLen]))
}

func establishConnection() net.Conn {
	connection, err := net.Dial(SERVER_TYPE, ADDRESS)
	if err != nil {
		panic(err)
	}
	return connection
}
