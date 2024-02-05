package main

import (
	"context"
	"fmt"
	"github.com/matnich89/benefex/common/rabbitmq"
	"github.com/matnich89/benefex/communcation/handler"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type app struct {
	conn       *amqp.Connection
	ch         *amqp.Channel
	queueName  string
	msgHandler *handler.MessageHandler
	errC       chan error
	sigC       chan os.Signal
}

func newApp(rabbitUrl string, msgHandler *handler.MessageHandler, queueName string, errC chan error) (*app, error) {
	conn, err := rabbitmq.OpenConnection(rabbitUrl)

	if err != nil {
		return nil, fmt.Errorf("could not conenct to rabbit mq %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("could not create channel %s", err)
	}

	return &app{conn: conn, msgHandler: msgHandler, ch: ch, queueName: queueName, errC: errC, sigC: make(chan os.Signal, 1)}, nil
}

func (a *app) run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := rabbitmq.DeclareQueue(a.ch, a.queueName); err != nil {
		return fmt.Errorf("could not declare queue %s", err)
	}

	messageHandler := func(d amqp.Delivery) {
		a.msgHandler.HandleMessage(d)
	}

	err := rabbitmq.ConsumeMessage(a.conn, a.queueName, messageHandler)

	if err != nil {
		return fmt.Errorf("could not being consuming messages %s queue %s", a.queueName, err)
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

	return nil
}
