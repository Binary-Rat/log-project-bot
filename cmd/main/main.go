package main

import (
	"log"
	"log-proj/internal/events/tg"
	"log-proj/internal/source/ati"
	tgClient "log-proj/pkg/clients/tg"
	event_consumer "log-proj/pkg/consumer/event-consumer"
	"log-proj/pkg/db/array"
	"log-proj/pkg/fsm/redis"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

const (
	tgBotHost = "api.telegram.org"
	batchSize = 100
)

//TODO Add config

// main starts bot. It takes one flag -t which is a token for Telegram API.
// If token is empty, the bot will panic.
func main() {
	tgBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if tgBotToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is empty!")
	}

	client := tgClient.New(tgBotHost, tgBotToken)
	//setup the db need to do it via config
	db := array.New()
	//setup the fsm
	fsm := redis.New()

	atiToken := os.Getenv("ATISU_TOKEN")
	if atiToken == "" {
		log.Fatal("ATISU_TOKEN is empty!")
	}

	ati, err := ati.New(atiToken, true)
	if err != nil {
		log.Fatal(err)
	}

	eventProc := tg.New(client, db, fsm, false, ati)

	log.Printf("service started")

	consumer := event_consumer.New(eventProc, eventProc, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}
