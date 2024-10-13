package services

import (
	"encoding/json"
	"fmt"
	"io"
	"music-lib/internal/config"
	"net/http"
	"net/url"
	"time"
)

type IMusicInfoService interface {
    GetSongInfo(artist, name string) (*SongDetail, error)
}
// MusicInfoService is a an external service, that provides additional information about songs
type MusicInfoService struct {
	baseURL url.URL
	client  *http.Client
}

type SongDetail struct {
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func NewMusicInfoService(cfg *config.Config) (*MusicInfoService, error) {
	cl := &http.Client{
		Timeout: time.Duration(cfg.ExternalAPI.Timeout) * time.Second,
	}
	url, err := url.Parse(cfg.ExternalAPI.BaseURL)
	if err != nil {
		return nil, err
	}

	return &MusicInfoService{
		baseURL: *url,
		client:  cl,
	}, nil
}

func (ms *MusicInfoService) GetSongInfo(artist, name string) (*SongDetail, error) {
	// Construct URL for the request
	ms.baseURL.Path = "/info"
	ms.baseURL.Query().Add("group", artist)
	ms.baseURL.Query().Add("song", name)
	// Send request
	resp, err := ms.client.Get(ms.baseURL.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		return nil, fmt.Errorf("failed to get song info: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var songDetail SongDetail
	if err := json.Unmarshal(body, &songDetail); err != nil {
		return nil, err
	}

	return &songDetail, nil
}
