package utils

type SongPostRequest struct {
	Group string `json:"group" validate:"required"`
	Song  string `json:"song" validate:"required"`
}

type SongPatchRequest struct {
	Group       string     `json:"group"`
	Song        string     `json:"song"`
	Lyrics      string     `json:"lyrics"`
	ReleaseDate CustomDate `json:"release_date"`
	URL         string     `json:"url"`
}

type SongPutRequest struct {
	Group       string     `json:"group" validate:"required"`
	Song        string     `json:"song" validate:"required"`
	Lyrics      string     `json:"lyrics" validate:"required"`
	ReleaseDate CustomDate `json:"release_date" validate:"required"`
	URL         string     `json:"url" validate:"required"`
}
