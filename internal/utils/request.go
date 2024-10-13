package utils

type SongRequest struct {
    Artist string `json:"artist" validate:"required"`
    Name   string `json:"name" validate:"required"`
}

type SongPatchRequest struct {
    Artist      string  `json:"artist"`
    Name        string  `json:"name"`
    Lyrics      string  `json:"lyrics"`
    ReleaseDate string  `json:"release_date" validate:"datetime=2006-01-02"`
    URL         string  `json:"url" validate:"url"`
}

type SongPutRequest struct {
    Artist      string  `json:"artist" validate:"required"`
    Name        string  `json:"name" validate:"required"`
    Lyrics      string  `json:"lyrics" validate:"required"`
    ReleaseDate string  `json:"release_date" validate:"required, datetime=2006-01-02"`
    URL         string  `json:"url" validate:"required, url"`
}

