package client

import (
	"github.com/matnich89/benefex/communcation/model"
	"log"
)

type StubEmailClient struct {
}

func (e *StubEmailClient) SendEmail(email model.FanEmail) {
	log.Printf("sending email mesage: %s  to emailaddress %s", email.Message(), email.EmailAddress)
}
