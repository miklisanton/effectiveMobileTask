package utils

type SongPostRequest struct {
    Group string `json:"group" validate:"required" example:"Artist or group name"`
    Song  string `json:"song" validate:"required" example:"Song name"`
}

type SongPatchRequest struct {
    Group       string     `json:"group" example:"Artist or group name"`
    Song        string     `json:"song" example:"Song name"`
    Lyrics      string     `json:"lyrics" example:"New lyrics of the song"`
    ReleaseDate CustomDate `json:"release_date" format:"string" example:"02.01.2006"`
    URL         string     `json:"url" example:"https://www.youtube.com/watch?v=12345"`
}

type SongPutRequest struct {
    Group       string     `json:"group" validate:"required" example:"Artist or group name"`
    Song        string     `json:"song" validate:"required" example:"Song name"`
    Lyrics      string     `json:"lyrics" validate:"required" example:"New lyrics of the song"`
	ReleaseDate CustomDate `json:"release_date" validate:"required" format:"string" example:"02.01.2006"`
    URL         string     `json:"url" validate:"required" example:"https://www.youtube.com/watch?v=12345"`
}
