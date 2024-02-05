package handler

import (
	"encoding/json"
	"errors"
	common "github.com/matnich89/benefex/common/model"
	"github.com/matnich89/benefex/distribution/model"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

const (
	dateFormat = "2006-01-02"
)

type VinylSender interface {
	SendVinylOrder(order model.VinylOrder) error
}

type CDSender interface {
	SendCDOrder(order model.CDOrder) error
}

type MessageHandler struct {
	ErrC        chan<- error
	vinylClient VinylSender
	cdClient    CDSender
}

func NewMessageHandler(vinylClient VinylSender, cdClient CDSender, errC chan<- error) *MessageHandler {
	return &MessageHandler{
		ErrC:        errC,
		vinylClient: vinylClient,
		cdClient:    cdClient,
	}
}

func (h *MessageHandler) HandleMessage(d amqp.Delivery) {
	log.Printf("recieved a message: %s", d.Body)

	var release common.Release

	if err := json.Unmarshal(d.Body, &release); err != nil {
		log.Printf("failed to unmarshal message: %v", err)
		return
	}

	for _, val := range release.Distribution {

		if val.Type != "vinyl" && val.Type != "cd" {
			log.Printf("unknown distribution type: %s", val.Type)
			h.ErrC <- errors.New("unknown distribution type received: %s")
			continue
		}

		if val.Type == "vinyl" {
			order := model.VinylOrder{
				Artist:        release.Artist,
				Title:         release.Title,
				Quantity:      val.Qty,
				DateOfRelease: release.ReleaseDate.Format(dateFormat),
			}
			if err := h.vinylClient.SendVinylOrder(order); err != nil {
				log.Printf("failed to send vinyl order: %v", err)
				h.ErrC <- err
				return
			}
		}
		if val.Type == "cd" {
			order := model.CDOrder{
				Artist:      release.Artist,
				Album:       release.Title,
				Quantity:    val.Qty,
				ReleaseDate: release.ReleaseDate.Format(dateFormat),
			}
			if err := h.cdClient.SendCDOrder(order); err != nil {
				log.Printf("failed to send CD order: %v", err)
				h.ErrC <- err
				return
			}
		}
	}

	if err := d.Ack(false); err != nil {
		h.ErrC <- errors.New("failed to acknowledge message")
	}
}
