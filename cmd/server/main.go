package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	connstr := "amqp://guest:guest@localhost:5672/"

	Conn, err := amqp.Dial(connstr)
	if err != nil {
		log.Fatalf("Error connecting %v", err)
	}
	defer Conn.Close()
	fmt.Println("Connection successfully:")

	newChan, err := Conn.Channel()
	if err != nil {
		log.Fatalf("Err creating channel: %v", err)
	}

	err = pubsub.SubscribeGob(
		Conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.SimpleQueueDurable,
		handlerLogs(),
	)
	if err != nil {
		log.Fatalf("could not starting consuming logs: %v", err)
	}

	gamelogic.PrintServerHelp()
	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}

		switch input[0] {
		case "pause":
			fmt.Println("Publishing pause game state")
			err = pubsub.PublishJSON(newChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: true,
			})
			if err != nil {
				log.Fatalf("error publishing pause message: %v", err)
			}
		case "resume":
			fmt.Println("Publishing resume game state")
			err = pubsub.PublishJSON(newChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: false,
			})
			if err != nil {
				log.Fatalf("error publishing resume message: %v", err)
			}
		case "quit":
			fmt.Println("Quiting the game")
			return
		default:
			fmt.Println("Unknown command")
		}

	}
}
