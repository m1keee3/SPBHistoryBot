package main

import (
	"flag"
	"log"
)

func main() {
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
