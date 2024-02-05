package model

import (
	"fmt"
	"time"
)

const (
	releaseDateFormat = "2006-01-02"
)

type Fan struct {
	Name         string
	EmailAddress string
}

type FanEmail struct {
	Fan
	Artist      string    `json:"artist"`
	Title       string    `json:"album"`
	ReleaseDate time.Time `json:"release_date"`
}

func (f FanEmail) Message() string {
	releaseDateStr := f.ReleaseDate.Format(releaseDateFormat)
	return fmt.Sprintf("Dear %s one of your favourite artists %s is releasing their new title %s, it will be available %s. Rock and Roll!! 🤘", f.Fan.Name, f.Artist, f.Title, releaseDateStr)
}
