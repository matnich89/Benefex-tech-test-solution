package model

import (
	"fmt"
	"github.com/matnich89/benefex/common/model"
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
	model.Release
}

func (f FanEmail) Message() string {
	releaseDateStr := f.ReleaseDate.Format(releaseDateFormat)
	return fmt.Sprintf("Dear %s one of your favourite artists %s is releasing their new title %s, it will be available %s. Rock and Roll!! 🤘", f.Fan.Name, f.Artist, f.Title, releaseDateStr)
}
