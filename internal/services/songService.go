package services

import (
	"context"
	"fmt"
	"music-lib/internal/db/models"
	"music-lib/internal/db/repository"
	"music-lib/internal/utils"

	"github.com/rs/zerolog/log"
)

type ISongService interface {
    CreateSong(ctx context.Context, song *models.Song) error
    GetSongs(ctx context.Context, f repository.SongFilter, page, limit int) ([]models.Song, error) 
    GetSong(ctx context.Context, id int) (*models.Song, error)
    UpdateSong(ctx context.Context, song, newSong *models.Song) (*models.Song, error)
    DeleteSong(ctx context.Context, id int) error
}

type SongService struct {
    Repo repository.ISongRepo
}

func NewSongService(songRepo repository.ISongRepo) ISongService {
    return SongService{songRepo}
}

func (s SongService) CreateSong(ctx context.Context, song *models.Song) error {
    if song.ID != nil {
        return fmt.Errorf("ID should not be set for a new song")
    }
    err := s.Repo.Save(ctx, song)
    if err != nil {
        log.Logger.Error().Err(err).Msgf("failed to save song")
        return err
    }
    return nil
}

func (s SongService) GetSong(ctx context.Context, id int) (*models.Song, error) {
    song, err := s.Repo.GetById(ctx, id)
    if err != nil {
        log.Logger.Error().Err(err).Msgf("failed to get song with id %d", id)
        return nil, err
    }
    return song, nil
}


// GetSongs fetches songs from database using filter and pagination
func (s SongService) GetSongs(ctx context.Context, f repository.SongFilter, page, limit int) ([]models.Song, error) {
    // Calculate offset
    offset := (page - 1) * limit
    songs, err := s.Repo.GetFiltered(ctx, f, offset, limit)
    if err != nil {
        log.Logger.Error().Err(err).Msgf("failed to get songs")
        return nil, err
    }

    return songs, nil
}

func (s SongService) UpdateSong(ctx context.Context, song, newSong *models.Song) (*models.Song, error) {
    // Update provided fields
    if newSong.Artist != "" {
        song.Artist = newSong.Artist
    }
    if newSong.Name != "" {
        song.Name = newSong.Name
    }
    // Empty date
    t := utils.CustomDate{}
    if newSong.ReleaseDate != t {
        song.ReleaseDate = newSong.ReleaseDate
    }
    if newSong.Lyrics != "" {
        song.Lyrics = newSong.Lyrics
    }
    if newSong.URL != "" {
        song.URL = newSong.URL
    }
    // Save updated song
    if err := s.Repo.Save(ctx, song); err != nil {
        log.Logger.Error().Err(err).Msgf("failed to update song with id %d", song.ID)
        return nil, err
    }
    return song, nil
}

func (s SongService) DeleteSong(ctx context.Context, id int) error {
    err := s.Repo.Delete(ctx, id)
    if err != nil {
        log.Logger.Error().Err(err).Msgf("failed to delete song with id %d", id)
        return err
    }
    return nil
}
