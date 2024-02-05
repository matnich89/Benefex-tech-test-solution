package model

import "time"

type Release struct {
	Artist       string                `json:"artist"`
	Title        string                `json:"title"`
	Genre        string                `json:"genre"`
	ReleaseDate  time.Time             `json:"releaseDate"`
	Distribution []ReleaseDistribution `json:"distribution"`
}

type ReleaseDistribution struct {
	Type string `json:"type"`
	Qty  int64  `json:"qty"`
}
