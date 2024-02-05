package model

/*
Obviously I don't have a concreate idea of what the payloads
Would like look like for each distributors api, but
these feel like a good attempt and always easy to change
later
*/
type CDOrder struct {
	Artist      string `json:"artist"`
	Album       string `json:"album"`
	Quantity    int64  `json:"quantity"`
	ReleaseDate string `json:"releaseDate"`
}

type VinylOrder struct {
	Artist        string `json:"artist"`
	Title         string `json:"title"`
	Quantity      int64  `json:"quantity"`
	DateOfRelease string `json:"dateOfRelease"`
}
