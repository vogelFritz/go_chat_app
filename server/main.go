package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var connections [MAX_CLIENTS]net.Conn

const (
	SERVER_TYPE = "tcp"
	SERVER_HOST = "localhost"
	SERVER_PORT = "8080"
	ADDRESS     = SERVER_HOST + ":" + SERVER_PORT
	MAX_CLIENTS = 10
)

func main() {
	dbInit()
	defer db.Close()
	server := startServer()
	defer server.Close()
	fmt.Println("Listening on " + ADDRESS)
	fmt.Println("Waiting for client...")
	waitForClients(server)
}

func dbInit() {
	os.Remove("sqlite-database.db")
	file, err := os.Create("sqlite-database.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	db, _ = sql.Open("sqlite3", "sqlite-database.db")

	createTable(db)

	err = db.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func createTable(db *sql.DB) {
	createMessageTableSql := `CREATE TABLE messages (
		"id" integer NOT NULL PRIMARY KEY autoincrement,
		"emisorName" CHAR(50) NOT NULL,
		"message" CHAR(255)
	);`
	statement, err := db.Prepare(createMessageTableSql)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
}

func insertMessage(emisorName string, message string) {
	insertMessageSql := `INSERT INTO messages (emisorName, message) VALUES(?, ?)`
	statement, err := db.Prepare(insertMessageSql)
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = statement.Exec(emisorName, message)
	if err != nil {
		log.Fatal(err.Error())
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

func waitForClients(server net.Listener) {
	var connIndex int = -1
	for {
		connection, err := server.Accept()
		connIndex++
		connections[connIndex] = connection
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("client connected")
		go listenToClient(connIndex)
	}
}

func listenToClient(clientIndex int) {
	sendRefreshedChat()
	receivedMessage := readMessage(connections[clientIndex])
	fmt.Println("Received: ", receivedMessage)
	for receivedMessage != "f" {
		insertMessage("vogel", receivedMessage)
		sendRefreshedChat()
		receivedMessage = readMessage(connections[clientIndex])
		fmt.Println("Received: ", receivedMessage)
	}
	connections[clientIndex].Write([]byte("Aufwiedersehen"))
	connections[clientIndex].Close()
}

func sendRefreshedChat() {
	var chat string
	row, err := db.Query("SELECT emisorName, message FROM messages")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer row.Close()
	for row.Next() {
		var emisorName string
		var message string
		row.Scan(&emisorName, &message)
		chat += emisorName + ": " + message + "\n"
	}
	for i := 0; connections[i] != nil; i++ {
		connections[i].Write([]byte(chat))
	}
}

func readMessage(connection net.Conn) string {
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
	}
	receivedMessage := string(buffer[:mLen])
	return receivedMessage
}
