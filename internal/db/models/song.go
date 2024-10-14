package models

import (
	"music-lib/internal/utils"
)

type Song struct {
	ID          *int             `db:"id" json:"id" example:"1"`
	Name        string           `db:"name" json:"song" example:"Song name"`
	Artist      string           `db:"artist" json:"group" example:"Artist or group name"`
	Lyrics      string           `db:"lyrics" json:"lyrics" example:"Lyrics of the song"`
	ReleaseDate utils.CustomDate `db:"release_date" json:"release_date" format:"string" example:"02.01.2006"`
	URL         string           `db:"url" json:"url" example:"https://www.youtube.com/watch?v=12345"`
}
