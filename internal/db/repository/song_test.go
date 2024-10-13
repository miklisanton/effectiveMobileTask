package repository

import (
	"context"
	"fmt"
	"music-lib/internal/db/drivers"
	"music-lib/internal/db/models"
	"music-lib/internal/utils"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var songRepo ISongRepo

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, ".env")); os.IsNotExist(err) {
			parent := filepath.Dir(dir)
			if parent == dir {
				// Reached the root of the filesystem, and .env wasn't found
				return "", os.ErrNotExist
			}
			dir = parent
		} else {
			return dir, nil
		}
	}
}

func PrintSong(song *models.Song) {
	fmt.Printf("\nID: %d\n Name: %s\n Artist: %s\n Lyrics: %s\n ReleaseDate: %s\n URL: %s\n",
		*song.ID, song.Name, song.Artist, song.Lyrics, song.ReleaseDate.Format("02.01.2006"), song.URL)
}

func TestMain(m *testing.M) {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()
	log.Info().Msg("Logger initialized")

	root, err := findProjectRoot()
	if err != nil {
		panic(fmt.Sprintf("Failed to find project root: %v", err))
	}
	if err = os.Chdir(root); err != nil {
		panic(fmt.Sprintf("Failed to change directory: %v", err))
	}

	db, err := drivers.Connect("postgresql://anton:1111@localhost:5432/test_db?sslmode=disable")
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	songRepo = NewSongRepository(db)

	m.Run()

	err = drivers.Down(db)
	if err != nil {
		panic(fmt.Sprintf("Failed to run cleanup: %v", err))
	}
}

func TestSave(t *testing.T) {
	// Create a new song
	song := models.Song{
		Name:        "Song Name",
		Artist:      "Song Artist",
		Lyrics:      "Song Lyrics",
		ReleaseDate: utils.CustomDate(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
		URL:         "https://song.url",
	}
	// Save the song
	err := songRepo.Save(context.Background(), &song)
	if err != nil {
		t.Fatalf("Error saving song: %v", err)
	}

	t.Logf("Saved song with id: %d", *song.ID)
}

func TestSave2(t *testing.T) {
	// Create a new song
	song := models.Song{
		Name:   "Gangsta's Paradise",
		Artist: "Coolio",
		Lyrics: "As I walk through the valley of the shadow of death\n" +
			"I take a look at my life and realize there's not much left\n" +
			"Cause I've been blastin' and laughin' so long that\n" +
			"Even my mama thinks that my mind is gone\n" +
			"But I ain't never crossed",
		ReleaseDate: utils.CustomDate(time.Date(1995, 11, 7, 0, 0, 0, 0, time.UTC)),
		URL:         "https://www.youtube.com/watch?v=fPO76Jlnz6c",
	}
	// Save the song
	err := songRepo.Save(context.Background(), &song)
	if err != nil {
		t.Fatalf("Error saving song: %v", err)
	}

	t.Logf("Saved song with id: %d", *song.ID)
}

func TestSave3(t *testing.T) {
	// Create a new song
	song := models.Song{
		Name:        "Song Name 2",
		Artist:      "Song Artist",
		Lyrics:      "Song Lyrics",
		ReleaseDate: utils.CustomDate(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
		URL:         "https://song.url",
	}
	// Save the song
	err := songRepo.Save(context.Background(), &song)
	if err != nil {
		t.Fatalf("Error saving song: %v", err)
	}

	t.Logf("Saved song with id: %d", *song.ID)
}

func TestSaveSame(t *testing.T) {
	// Create a new song
	song := models.Song{
		Name:        "Song Name",
		Artist:      "Song Artist",
		Lyrics:      "Song Lyrics",
		ReleaseDate: utils.CustomDate(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
		URL:         "https://song.com",
	}
	// Save the song
	err := songRepo.Save(context.Background(), &song)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	t.Logf("Error saving song: %v", err)
}

func TestUpdate(t *testing.T) {
	// Update existing song
	id := 1
	song := models.Song{
		ID:          &id,
		Name:        "Song Name",
		Artist:      "Song Artist",
		Lyrics:      "Song Updated Lyrics",
		ReleaseDate: utils.CustomDate(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
		URL:         "https://updated.song.com",
	}
	// Save the song
	err := songRepo.Save(context.Background(), &song)
	if err != nil {
		t.Fatalf("Error updating song: %v", err)
	}

	t.Logf("Updated song with id: %d", *song.ID)
}

func TestGetById(t *testing.T) {
	song, err := songRepo.GetById(context.Background(), 1)
	if err != nil {
		t.Fatalf("Error getting song by name: %v", err)
	}

	PrintSong(song)
}

func TestGetByIdNotFound(t *testing.T) {
	song, err := songRepo.GetById(context.Background(), 1000)
	if err == nil {
		t.Fatalf("Expected error, got nil and song: %v", song)
	}

	t.Logf("Error while trying to get song that doesnt exist: %v", err)
}

func TestGetFiltered(t *testing.T) {
	filter := SongFilter{
		After:  "1995-11-08",
		Artist: "Song Artist",
	}

	songs, err := songRepo.GetFiltered(context.Background(), filter, 0, 2)
	if err != nil {
		t.Fatalf("Error getting filtered songs: %v", err)
	}

	for _, song := range songs {
		PrintSong(&song)
	}
}

func TestDelete(t *testing.T) {
    err := songRepo.Delete(context.Background(), 1)
	if err != nil {
		t.Fatalf("Error deleting song: %v", err)
	}
}

func TestDeleteNotExist(t *testing.T) {
    err := songRepo.Delete(context.Background(), 100)
	if err == nil {
		t.Fatalf("Expected error, got nil and deleted")
	}
}

func TestGetAll(t *testing.T) {
	songs, err := songRepo.GetAll(context.Background())
	if err != nil {
		t.Fatalf("Error getting all songs: %v", err)
	}

	for _, song := range songs {
		PrintSong(&song)
	}
}
