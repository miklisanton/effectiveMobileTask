package services

import (
	"encoding/json"
	"fmt"
	"io"
	"music-lib/internal/config"
	"music-lib/internal/utils"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type IMusicInfoService interface {
    GetSongInfo(artist, name string) (*SongDetail, error)
}
// MusicInfoService is a an external service, that provides additional information about songs
type MusicInfoService struct {
	baseURL string
	client  *http.Client
}

type SongDetail struct {
    ReleaseDate utils.CustomDate `json:"releaseDate" validate:"required"`
    Text        string `json:"text" validate:"required"`
    Link        string `json:"link" validate:"required"`
}

func NewMusicInfoService(cfg *config.Config) (*MusicInfoService, error) {
	cl := &http.Client{
		Timeout: time.Duration(cfg.ExternalAPI.Timeout) * time.Second,
	}

	return &MusicInfoService{
		baseURL: cfg.ExternalAPI.BaseURL,
		client:  cl,
	}, nil
}

func (ms *MusicInfoService) GetSongInfo(artist, name string) (*SongDetail, error) {
	// Construct URL for the request
    u, err := url.Parse(ms.baseURL)
    if err != nil {
        return nil, err
    }
    u.Path = strings.TrimRight(u.Path, "/") + "/info" 
    queryParams := url.Values{}
	queryParams.Add("group", artist)
	queryParams.Add("song", name)
    u.RawQuery = queryParams.Encode()
	// Send request
    log.Info().Msgf("Sending request to %s", u.String())
	resp, err := ms.client.Get(u.String())
	if err != nil {
        log.Error().Err(err).Msg("failed to send request")
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body")
	}

	var songDetail SongDetail
	if err := json.Unmarshal(body, &songDetail); err != nil {
        log.Error().Err(err).Msg("failed to unmarshal response")
		return nil, fmt.Errorf("failed to unmarshal response")
	}
    if err := validator.New().Struct(songDetail); err != nil {
        return nil, fmt.Errorf("invalid response: %v", err)
    }

    log.Debug().Msgf("Release date: %s", songDetail.ReleaseDate)

	return &songDetail, nil
}
