package main

import (
	"context"
	"github.com/matnich89/benefex/common/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type app struct {
	conn      *amqp.Connection
	ch        *amqp.Channel
	queueName string
	errC      chan error
	sigC      chan os.Signal
}

func newApp(rabbitUrl, queueName string, errC chan error) (*app, error) {
	conn, err := rabbitmq.OpenConnection(rabbitUrl)

	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &app{conn: conn, ch: ch, queueName: queueName, errC: errC, sigC: make(chan os.Signal, 1)}, nil
}

func (a *app) run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := rabbitmq.DeclareQueue(a.ch, a.queueName); err != nil {
		log.Fatalln(err)
	}

	messageHandler := func(d amqp.Delivery) {
		log.Println("message received will distribute to customers")
	}

	err := rabbitmq.ConsumeMessage(a.conn, a.queueName, messageHandler)

	if err != nil {
		log.Fatalln(err)
	}

	signal.Notify(a.sigC, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-a.sigC
		log.Println("cleaning up...")
		if err := a.conn.Close(); err != nil {
			log.Printf("error closing rabbitmq connection: %v", err)
		}
		cancel()
		close(a.errC)
	}()

	go func() {
		for err := range a.errC {
			log.Printf("error received: %v, sending alert!!!!", err)
		}
	}()

	<-ctx.Done()
	log.Println("application stopping....")
}
