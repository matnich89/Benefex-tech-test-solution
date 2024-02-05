package main

import (
	"github.com/matnich89/benefex/communcation/client"
	"github.com/matnich89/benefex/communcation/db"
	"github.com/matnich89/benefex/communcation/handler"
	"log"
	"os"
)

func main() {
	rabbitUrl, ok := os.LookupEnv("RABBITMQ_SERVER_URL")

	if !ok {
		log.Fatalln("could not find RABBITMQ_SERVER_URL env var")
	}

	errCh := make(chan error)

	database := db.NewFanbaseDB()

	emailClient := client.StubEmailClient{}

	msgHandler := handler.NewMessageHandler(errCh, database, &emailClient)

	app, err := newApp(rabbitUrl, msgHandler, "communication", errCh)

	if err != nil {
		log.Fatalf("error creating app: %v", err)
	}

	err = app.run()

	if err != nil {
		log.Fatalf("error running app: %v", err)
	}
}
