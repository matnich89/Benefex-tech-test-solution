package db

import (
	"errors"
	"github.com/matnich89/benefex/communcation/model"
)

type FanbaseDB struct {
	fans map[string][]model.Fan
}

func NewFanbaseDB() *FanbaseDB {
	db := &FanbaseDB{make(map[string][]model.Fan)}
	db.populateStubData()
	return db
}

func (db *FanbaseDB) GetFansForArtist(artistName string) ([]model.Fan, error) {
	fans, ok := db.fans[artistName]
	if !ok {
		return []model.Fan{}, errors.New("artist not found")
	}
	return fans, nil
}

func (db *FanbaseDB) populateStubData() {

	db.fans["The Beatles"] = []model.Fan{} // I don't like the Beatles, so they have no fans :)

	db.fans["Epoch-alypse"] = []model.Fan{
		{Name: "Fan1", EmailAddress: "fan1@example.com"},
		{Name: "Fan2", EmailAddress: "fan2@example.com"},
	}

	db.fans["Elon Dusk"] = []model.Fan{
		{Name: "Fan2", EmailAddress: "fan2@example.com"},
	}

	db.fans["Epica"] = []model.Fan{
		{Name: "Fan3", EmailAddress: "fan3@example.com"},
		{Name: "Fan4", EmailAddress: "fan4@example.com"},
		{Name: "Fan5", EmailAddress: "fan5@example.com"},
		{Name: "Fan6", EmailAddress: "fan6@example.com"},
	}
}
