package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "8080"
	SERVER_TYPE = "tcp"
	ADDRESS     = SERVER_HOST + ":" + SERVER_PORT
)

func main() {
	connection := establishConnection()
	register(connection)
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

func register(connection net.Conn) {
	fmt.Print("What is your name? ==> ")
	name := getMessageFromUser()
	sendMessage(name, connection)
}

func chatLoop(connection net.Conn) {
	var message string
	go waitForUpdates(connection)
	fmt.Print("Write your message here (f to finish): ")
	message = getMessageFromUser()
	sendMessage(message, connection)
	for message != "f" {
		message = getMessageFromUser()
		sendMessage(message, connection)
	}
}

func getMessageFromUser() string {
	reader := bufio.NewReader(os.Stdin)
	message, err := reader.ReadString('\n')
	message = message[:len(message)-1]
	for err != nil {
		fmt.Println("Error with the message: ", err.Error())
		_, err = fmt.Scanln(&message)
	}
	return message
}

func sendMessage(message string, connection net.Conn) {
	connection.Write([]byte(message))

}

func waitForUpdates(connection net.Conn) {
	for {
		buffer := make([]byte, 1024)
		mLen, err := connection.Read(buffer)
		if err != nil {
			fmt.Println("Error reading: ", err.Error())
		}
		clearScreen()
		fmt.Println(string(buffer[:mLen]))
		fmt.Print("Write your message here (f to finish): ")
	}
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
