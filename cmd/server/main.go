package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Tony-Ledoux/peril/internal/pubsub/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	// Connection string for RabbitMQ
	connStr := "amqp://guest:guest@localhost:5672/"

	// Connect to RabbitMQ
	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %v", err)
	}
	defer func() {
		fmt.Println("Closing RabbitMQ connection...")
		_ = conn.Close()
	}()

	fmt.Println("Connection to RabbitMQ successful.")

	// create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error creating channel: %v", err)
	}
	defer ch.Close()

	fmt.Println("Channel created successfully.")

	// publish a pause message
	// Publish a pause message
	err = pubsub.PublishJSON(
		ch,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		routing.PlayingState{IsPaused: true},
	)
	if err != nil {
		log.Printf("Error publishing pause message: %v", err)
	} else {
		fmt.Println("Pause message published successfully.")
	}

	// Wait for Ctrl+C (SIGINT)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	fmt.Println("Press Ctrl+C to exit.")
	<-signalChan

	fmt.Println("Signal received. Shutting down server...")
}
