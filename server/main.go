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

const (
	SERVER_TYPE = "tcp"
	SERVER_HOST = "localhost"
	SERVER_PORT = "8080"
	ADDRESS     = SERVER_HOST + ":" + SERVER_PORT
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
	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("client connected")
		go listenToClient(connection)
	}
}

func listenToClient(connection net.Conn) {
	receivedMessage := readMessage(connection)
	fmt.Println("Received: ", receivedMessage)
	for receivedMessage != "f" {
		connection.Write([]byte("Thanks! Got your message:" + receivedMessage))
		insertMessage("vogel", receivedMessage)
		receivedMessage = readMessage(connection)
	}
	connection.Write([]byte("Aufwiedersehen"))
	connection.Close()
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
