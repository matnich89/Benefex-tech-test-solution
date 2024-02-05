package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/matnich89/benefex/common/model"
	"github.com/matnich89/benefex/common/rabbitmq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type queueClient struct {
	conn               *amqp.Connection
	communicationQueue string
	distributionQueue  string
	errCh              chan<- error
}

func (q *queueClient) Send(ctx context.Context, r model.Release) {
	var wg sync.WaitGroup

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	if err := encoder.Encode(r); err != nil {
		q.errCh <- err
		return
	}

	// We are always sending to comms channel
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := q.sendToQueue(ctx, q.communicationQueue, buf.Bytes()); err != nil {
			q.errCh <- err
			return
		}
	}()

	// We only send to distribution if the info is available
	if len(r.Distribution) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := q.sendToQueue(ctx, q.distributionQueue, buf.Bytes()); err != nil {
				q.errCh <- err
			}
		}()
	}

	wg.Wait()
}

func (q *queueClient) sendToQueue(ctx context.Context, queueName string, body []byte) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := q.publishMessage(ctx, queueName, body)

	if err != nil {
		return err
	}

	return nil
}

func (q *queueClient) publishMessage(ctx context.Context, queueName string, body []byte) error {
	ch, err := q.conn.Channel()

	if err != nil {
		return errors.New(fmt.Sprintf("error occured trying to obtain channel for queue %s", queueName))
	}

	err = rabbitmq.DeclareQueue(ch, queueName)

	if err != nil {
		return errors.New(fmt.Sprintf("unable to declare queue for queue name %s", queueName))
	}

	err = ch.PublishWithContext(ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	select {
	case <-ctx.Done():
		return fmt.Errorf("sending message to %s timed out", queueName)
	default:
		if err != nil {
			return fmt.Errorf("error publishing message to %s: %v", queueName, err)
		}
	}

	return nil
}

func newConnectedQueueClient(rabbitMqUrl string, errCh chan<- error) (*queueClient, error) {

	conn, err := rabbitmq.OpenConnection(rabbitMqUrl)

	if err != nil {
		return nil, fmt.Errorf("could not connect to rabbit %s", err)
	}

	return &queueClient{conn, "communication", "distribution", errCh}, nil
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatalf("run failed: %s", err.Error())
	}
}

func run(ctx context.Context) error {

	rabbitMqUrl, ok := os.LookupEnv("RABBITMQ_SERVER_URL")

	if !ok {
		return errors.New("could not find RABBITMQ_SERVER_URL env var")
	}

	errC := make(chan error)
	sigC := make(chan os.Signal)
	relC := make(chan model.Release)

	q, err := newConnectedQueueClient(rabbitMqUrl, errC)
	if err != nil {
		return fmt.Errorf("error connecting queue client: %w", err)
	}

	signal.Notify(sigC, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := runServer(":8080", newHandler(relC)); err != nil {
			errC <- err
			return
		}
	}()

	for {
		select {
		case rel := <-relC:
			q.Send(ctx, rel)
		case err := <-errC:
			return fmt.Errorf("error received: %w", err)
		case <-sigC:
			log.Println("cleaning up...")
			close(errC)
			close(relC)
			if err := q.conn.Close(); err != nil {
				log.Printf("error closing rabbitmq connection: %v", err)
			}
			return nil
		}
	}
}

func runServer(addr string, hnd http.Handler) error {
	log.Printf("listening on %s...", addr)

	if err := http.ListenAndServe(addr, hnd); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}

func newHandler(relC chan<- model.Release) http.Handler {
	hnd := http.NewServeMux()
	hnd.HandleFunc("/releases", releasesHandler(relC))
	return hnd
}

func releasesHandler(relC chan<- model.Release) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var inp []model.Release

		if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
			writeError(w, fmt.Errorf("cannot decode request: %w", err), http.StatusBadRequest)
			return
		}

		for _, rel := range inp {
			relC <- rel
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func writeError(w http.ResponseWriter, err error, status int) {
	type errBody struct {
		Message string `json:"message"`
	}

	resp := errBody{Message: err.Error()}
	b := &bytes.Buffer{}

	if err := json.NewEncoder(b).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("cannot encode response: %s", err.Error())
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(b.Bytes()); err != nil {
		log.Printf("cannot write response: %s", err.Error())
	}
}
