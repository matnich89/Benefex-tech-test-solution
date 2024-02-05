#!/bin/bash

curl -v http://localhost:8080/releases --data '
[
  {
    "artist": "The Beatles",
    "title": "Greatest Hits Vol. 4987843",
    "genre": "Pop",
    "releaseDate": "2021-02-09T00:00:00Z",
    "distribution": [
      {"type": "cd", "qty": 50000},
      {"type": "vinyl", "qty": 10000}
    ]
  },
  {
    "artist": "Epoch-alypse",
    "title": "Vinyl Countdown (Ltd. Edition)",
    "genre": "EDM",
    "releaseDate": "1970-01-01T00:00:00Z",
    "distribution": [
      {"type": "vinyl", "qty": 50}
    ]
  },
  {
    "artist": "Elon Dusk",
    "title": "Blockchain Bop (NFT Single)",
    "genre": "Spacecore",
    "releaseDate": "2024-03-28T00:00:00Z"
  },
  {
    "artist": "Epica",
    "title": "A Bright Dystopia",
    "genre": "Symphonic Metal",
    "releaseDate": "2024-05-25T00:00:00Z"
  },
  {
    "artist": "Mathew Nicholls",
    "title": "Not Real Album",
    "genre": "Mat Jam",
    "releaseDate": "2024-06-20T00:00:00Z"
   }
]
'
