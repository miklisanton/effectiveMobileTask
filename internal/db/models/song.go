package models

import (
	"music-lib/internal/utils"
)

type Song struct {
	ID          *int             `db:"id" json:"id"`
	Name        string           `db:"name" json:"song"`
	Artist      string           `db:"artist" json:"group"`
	Lyrics      string           `db:"lyrics" json:"lyrics"`
	ReleaseDate utils.CustomDate `db:"release_date" json:"release_date"`
	URL         string           `db:"url" json:"url"`
}
