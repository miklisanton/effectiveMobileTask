package models

type Song struct {
	ID          *int   `db:"id"`
	Name        string `db:"name"`
	Artist      string `db:"artist"`
	Lyrics      string `db:"lyrics"`
	ReleaseDate string `db:"release_date"`
	URL         string `db:"url"`
}
