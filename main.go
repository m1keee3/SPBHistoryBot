package main

import (
	tgClient "SPBHistoryBot/clients/telegram"
	event_consumer "SPBHistoryBot/consumer/event-consumer"
	"SPBHistoryBot/events/telegram"
	"SPBHistoryBot/lib/e"
	"SPBHistoryBot/lib/storage"
	"flag"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

const (
	tgBotHost = "api.telegram.org"
	dsn       = "host=localhost user=spb_user password=spb_pass dbname=spb_history port=5432 sslmode=disable"
	batchSize = 100
)

var (
	token = flag.String(
		"t",
		"",
		"Telegram bot access token")
)

func main() {

	flag.Parse()

	db, err := storage.NewDBStorage(dsn)
	if err != nil {
		log.Fatal(e.Wrap("can't connect to database", err))
	}

	if err := storage.Init(db); err != nil {
		log.Fatal(err)
	}

	if err := storage.SeedifEmpty(db); err != nil {
		log.Print(err)
	}

	eventsProcessor := telegram.NewProcessor(tgClient.NewClient(tgBotHost, mustToken()), db)
	log.Print("service started")

	consumer := event_consumer.NewConsumer(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service stopped", err)
	}
}

func mustToken() string {
	if *token == "" {
		log.Fatal("Token is required")
	}
	return *token
}
