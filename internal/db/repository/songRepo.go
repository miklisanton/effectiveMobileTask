package repository

import (
	"context"
	"database/sql"
	"fmt"
	"music-lib/internal/db/models"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type ISongRepo interface {
	GetAll(ctx context.Context) ([]models.Song, error)
	GetFiltered(ctx context.Context, filter SongFilter, offset int, limit int) ([]models.Song, error)
	GetById(ctx context.Context, id int) (*models.Song, error)
	Save(ctx context.Context, song *models.Song) error
	Delete(ctx context.Context, id int) error
}

type SongRepository struct {
	db *sqlx.DB
}

func NewSongRepository(db *sqlx.DB) ISongRepo {
	return &SongRepository{db}
}

// Save saves a song to db if id not set, otherwise updates the existing song
func (r *SongRepository) Save(ctx context.Context, song *models.Song) error {
	if song.ID != nil {
		// Update song
		query := `
            UPDATE song
            SET name=$1, artist=$2, lyrics=$3, release_date=$4, url=$5
            WHERE id=$6
            `
		log.Debug().Msgf("Running query: %s", query)
		_, err := r.db.ExecContext(ctx, query, song.Name, song.Artist, song.Lyrics, song.ReleaseDate, song.URL, *song.ID)
		if err != nil {
            if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
                // Unique violation
                return fmt.Errorf("Song with name %s and artist %s already exists", song.Name, song.Artist)
            }
			return err
		}
		return nil
	} else {
		// Create new song
		query := `
            INSERT INTO
            song(name, artist, lyrics, release_date, url)
            VALUES($1, $2, $3, $4, $5)
            RETURNING id
            `
		log.Debug().Msgf("Running query: %s", query)
		row := r.db.QueryRowContext(ctx, query, song.Name, song.Artist, song.Lyrics, song.ReleaseDate, song.URL)

		err := row.Err()
		if err != nil {
            if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
                // Unique violation
                return fmt.Errorf("Song with name %s and artist %s already exists", song.Name, song.Artist)
            }
			return err
		}
		err = row.Scan(&song.ID)
		if err != nil {
			return err
		}

		return nil
	}
}

func (r *SongRepository) GetAll(ctx context.Context) ([]models.Song, error) {
	songs := []models.Song{}
	query := `SELECT * FROM song ORDER BY id ASC`
	log.Debug().Msgf("Running query: %s", query)
	err := r.db.SelectContext(ctx, &songs, query)
	if err != nil {
		return nil, err
	}

	return songs, nil
}

func (r *SongRepository) GetById(ctx context.Context, id int) (*models.Song, error) {
	song := models.Song{}
	query := `SELECT * FROM song WHERE id=$1`
	log.Debug().Msgf("Running query: %s", query)
	err := r.db.GetContext(ctx, &song, query, id)
	if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("Song with id %d doesn't exist", id)
        }
		return nil, err
	}
	return &song, nil
}

func (r *SongRepository) GetFiltered(ctx context.Context, filter SongFilter, offset, limit int) ([]models.Song, error) {
	songs := []models.Song{}
	// Construct query from filter
	query := `SELECT * FROM song WHERE 1=1`
	log.Debug().Msgf("Running query: %s", query)
	count := 1
	if filter.Name != "" {
		query += ` AND name= :name`
		count++
	}
	if filter.Artist != "" {
		query += ` AND artist= :artist`
		count++
	}
	if filter.After != "" {
		query += ` AND release_date >= :after`
		count++
	}
	if filter.Before != "" {
		query += ` AND release_date <= :before`
		count++
	}

	boundQuery, filterArgs, err := r.db.BindNamed(query, filter)
	if err != nil {
		return nil, err
	}

	boundQuery += ` ORDER BY id ASC`
	boundQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", count, count+1)

	log.Debug().Msgf("Running query: %s", boundQuery)
	// Append limit and offset to the end of the query
	args := append(filterArgs, limit, offset)
	err = r.db.SelectContext(ctx, &songs, boundQuery, args...)
	if err != nil {
		return nil, err
	}

	return songs, nil
}

func (r *SongRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM song WHERE id = $1`
	log.Debug().Msgf("Running query: %s", query)
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
    if count != 1 {
        return fmt.Errorf("Song with id %d not found", id)
    }
	return nil
}
