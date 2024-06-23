package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/quackdiscord/events-api/events"
	"github.com/quackdiscord/events-api/services"
	"github.com/quackdiscord/events-api/structs"
	log "github.com/sirupsen/logrus"
)

func init() {
	// load .env file
	if err := godotenv.Load(".env"); err != nil {
		return
	}

	// set the environment
	env := os.Getenv("ENVIORNMENT")

	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	if env == "dev" {
		log.Warn("Running in development mode")
	}
}

func main() {
	// convert the map to a slice
	var eventSlice []*structs.Event
	for _, event := range events.Events {
		eventSlice = append(eventSlice, (*structs.Event)(event))
	}

	services.ConnectKafka(eventSlice)
	services.ConnectRedis()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Info("Press Ctrl+C to exit")

	// start consuming
	go services.KafkaReader.Consume(context.Background())

	// handle shutdown
	<-stop
	log.Warn("Shutting down")
	services.DisconnectKafka()
	services.DisconnectRedis()

	log.Info("Goodbye!")

}
