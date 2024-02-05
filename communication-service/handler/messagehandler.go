package handler

import (
	"encoding/json"
	common "github.com/matnich89/benefex/common/model"
	"github.com/matnich89/benefex/communcation/model"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
	"time"
)

const maxConcurrentSends = 200

type FanBaseStore interface {
	GetFansForArtist(artist string) ([]model.Fan, error)
}

type EmailSender interface {
	SendEmail(email model.FanEmail)
}

type MessageHandler struct {
	ErrC        chan<- error
	FanBaseDb   FanBaseStore
	EmailClient EmailSender
}

func NewMessageHandler(errC chan<- error, fanBaseDb FanBaseStore, emailClient EmailSender) *MessageHandler {
	return &MessageHandler{
		ErrC:        errC,
		FanBaseDb:   fanBaseDb,
		EmailClient: emailClient,
	}
}

func (h *MessageHandler) HandleMessage(d amqp.Delivery) {
	log.Printf("received a message: %s", d.Body)

	var release common.Release
	if err := json.Unmarshal(d.Body, &release); err != nil {
		log.Printf("failed to unmarshal message: %v", err)
		h.ErrC <- err
		return
	}

	if release.ReleaseDate.After(time.Now()) {
		fans, err := h.FanBaseDb.GetFansForArtist(release.Artist)
		if err != nil {
			log.Printf("error fetching fans for artist %s: %v", release.Artist, err)
			h.ErrC <- err
			return
		}

		sem := make(chan struct{}, maxConcurrentSends)

		var wg sync.WaitGroup
		for _, fan := range fans {
			wg.Add(1)
			sem <- struct{}{}

			go func(fan model.Fan) {
				defer wg.Done()
				email := model.FanEmail{
					Fan:         fan,
					Artist:      release.Artist,
					Title:       release.Title,
					ReleaseDate: release.ReleaseDate,
				}

				h.EmailClient.SendEmail(email)
				<-sem
			}(fan)
		}
	} else {
		log.Println("Release is already out, no need to notify the fanbase.")
	}
}
