package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	connstr := "amqp://guest:guest@localhost:5672/"

	Conn, err := amqp.Dial(connstr)
	if err != nil {
		log.Fatalf("Error connecting %v", err)
	}
	defer Conn.Close()
	fmt.Println("Connection successfully:")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
