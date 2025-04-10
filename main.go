package main

import (
	tgClient "SPBHistoryBot/clients/telegram"
	"SPBHistoryBot/consumer/event-consumer"
	"SPBHistoryBot/events/telegram"
	"flag"
	"log"
)

const (
	tgBotHost = "api.telegram.org"
	batchSize = 100
)

func main() {
	eventsProcessor := telegram.NewProcessor(tgClient.NewClient(tgBotHost, mustToken()))
	log.Print("service started")

	consumer := event_consumer.NewConsumer(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service stopped", err)
	}
}

func mustToken() string {
	token := flag.String(
		"t",
		"",
		"Telegram bot access token")

	flag.Parse()

	if *token == "" {
		log.Fatal("Token is required")
	}
	return *token
}
