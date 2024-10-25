package main

import (
	"flag"
	"log"
	"log-proj/internal/events/tg"
	tgClient "log-proj/pkg/clients/tg"
	event_consumer "log-proj/pkg/consumer/event-consumer"
	"log-proj/pkg/db/array"
	"log-proj/pkg/fsm/redis"
)

const (
	tgBotHost = "api.telegram.org"
	batchSize = 100
)

// main starts bot. It takes one flag -t which is a token for Telegram API.
// If token is empty, the bot will panic.
func main() {
	client := tgClient.New(tgBotHost, mustToken())
	//setup the db
	db := array.New()
	//setup the fsm
	fsm := redis.New()
	eventProc := tg.New(client, db, fsm)

	log.Printf("service started")

	consumer := event_consumer.New(eventProc, eventProc, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

// mustToken returns token for Telegram API.
// If token is empty, the function will fatal with appropriate error message.
func mustToken() string {
	t := flag.String("t", "", "Token for Telegram API")

	flag.Parse()

	if *t == "" {
		log.Fatal("Token is empty!")
	}
	return *t
}
