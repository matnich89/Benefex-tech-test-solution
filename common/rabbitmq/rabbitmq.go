package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

/*
In a prod system these constants would
be configurable
*/
const (
	maxRetries = 10
	baseDelay  = 100
)

func DeclareQueue(ch *amqp.Channel, queueName string) error {
	_, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	return nil
}

func OpenConnection(url string) (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error

	for i := 0; i < maxRetries; i++ {
		conn, err = attemptConnection(url)

		/*
		 exponential back off
		*/
		if err != nil {
			delay := time.Duration(baseDelay<<i) * time.Millisecond // '<<' is bit-shifting to left so multiples the basedelay by 2^i
			log.Printf("Attempt %d: Could not connect to rabbitMq, will retry in %v. Error: %v", i+1, delay, err)
			time.Sleep(delay)
		} else {
			log.Println("successfully connected to rabbitMq....")
			return conn, nil
		}
	}

	return nil, err
}

func ConsumeMessage(conn *amqp.Connection, queueName string, handler func(delivery amqp.Delivery)) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			handler(msg)
		}
	}()

	return nil
}

func attemptConnection(rabbitUrl string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(rabbitUrl)

	if err != nil {
		return nil, err
	}

	return conn, nil
}
