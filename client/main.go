package main

import (
	"fmt"
	"net"
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
		_, err := fmt.Scanf("%s", &message)
		for err != nil {
			fmt.Println("Error with the message: ", err.Error())
			_, err = fmt.Scanf("%s", &message)
		}
		if message == "f" {
			break
		}
		sendMessage(message, connection)
	}
}

func sendMessage(message string, connection net.Conn) {
	connection.Write([]byte(message))
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
