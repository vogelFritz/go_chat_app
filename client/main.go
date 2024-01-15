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

func establishConnection() net.Conn {
	connection, err := net.Dial(SERVER_TYPE, ADDRESS)
	if err != nil {
		panic(err)
	}
	return connection
}

func chatLoop(connection net.Conn) {
	var message string
	message = getMessageFromUser()
	sendMessage(message, connection)
	for message != "f" {
		message = getMessageFromUser()
		sendMessage(message, connection)
	}
}

func getMessageFromUser() string {
	var message string
	fmt.Print("Write your message here (f to finish): ")
	_, err := fmt.Scanln(&message)
	for err != nil {
		fmt.Println("Error with the message: ", err.Error())
		_, err = fmt.Scanln(&message)
	}
	return message
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
