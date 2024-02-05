package main

import (
	"log"
	"os"
)

func main() {
	rabbitUrl, ok := os.LookupEnv("RABBITMQ_SERVER_URL")

	if !ok {
		log.Fatalln("could not find RABBITMQ_SERVER_URL env var")
	}

	errCh := make(chan error)

	app, err := newApp(rabbitUrl, "distribution", errCh)

	if err != nil {
		log.Fatalf("error creating app: %v", err)
	}

	app.run()

}
