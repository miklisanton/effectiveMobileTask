package handlers

import (
	"context"
	"fmt"
	"music-lib/internal/config"
	"music-lib/internal/db/models"
	"music-lib/internal/db/repository"
	"music-lib/internal/services"
	"music-lib/internal/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type SongController struct {
	SongService      services.ISongService
	MusicInfoService services.IMusicInfoService
	Timeout          time.Duration
}

func NewSongController(
	songS services.ISongService,
	musicInfoS services.IMusicInfoService,
	cfg *config.Config) *SongController {

	timeout := time.Duration(cfg.Server.Timeout) * time.Second

	return &SongController{songS, musicInfoS, timeout}
}

// @Summary      Create a new song
// @Description  Create a new song by providing the group and song name. The song details are fetched from an external API.
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        song body utils.SongPostRequest true "Song request"
// @Success      201  {object}  utils.Response{message=string, data=models.Song} "Song created"
// @Failure      400  {object}  utils.Response{message=string, data=[]string} "Invalid request"
// @Failure      404  {object}  utils.Response{message=string} "Song not found"
// @Failure      409  {object}  utils.Response{message=string} "Song already exists"
// @Failure      500  {object}  utils.Response{message=string} "External API error or internal server error"
// @Router       /songs [post]
func (sc *SongController) CreateSong(c echo.Context) error {
	// New context with timeout
	ctx, cancel := context.WithTimeout(c.Request().Context(), sc.Timeout)
	defer cancel()
	// Extract song name and artist from request
	songRequest := new(utils.SongPostRequest)
	if err := c.Bind(songRequest); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{Message: err.Error()})
	}
	if err := validator.New().Struct(songRequest); err != nil {
		var errors []string
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("%s: %s", e.Field(), e.Tag()))
		}

		return c.JSON(
			http.StatusBadRequest,
			utils.Response{Message: "invalid request", Data: errors})
	}
	// Fetch song details from external service
	songDetail, err := sc.MusicInfoService.GetSongInfo(songRequest.Group, songRequest.Song)
	if err != nil {
		if err.Error() == "404" {
			return c.JSON(
				http.StatusNotFound,
				utils.Response{Message: "Song not found"})
		}
		return c.JSON(
			http.StatusInternalServerError,
			utils.Response{Message: "external api error: " + err.Error()})
	}
	// Save song to db
	song := &models.Song{
		Artist:      songRequest.Group,
		Name:        songRequest.Song,
		Lyrics:      songDetail.Text,
		URL:         songDetail.Link,
		ReleaseDate: songDetail.ReleaseDate,
	}
	if err := sc.SongService.CreateSong(ctx, song); err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return c.JSON(
				http.StatusConflict,
				utils.Response{Message: err.Error()})
		}
		return c.JSON(
			http.StatusInternalServerError,
			utils.Response{Message: err.Error()})
	}
	return c.JSON(
		http.StatusCreated,
		utils.Response{Message: "Song created", Data: song})
}

// @Summary      Get a song by ID
// @Description  Retrieve song details by ID
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Song ID"
// @Success      200  {object}  utils.Response{message=string, data=models.Song} "Song received"
// @Failure      400  {object}  utils.Response{message=string} "Invalid song ID"
// @Failure      404  {object}  utils.Response{message=string} "Song not found"
// @Failure      500  {object}  utils.Response{message=string} "Internal server error"
// @Router       /songs/{id} [get]
func (sc *SongController) GetSong(c echo.Context) error {
	// New context with timeout
	ctx, cancel := context.WithTimeout(c.Request().Context(), sc.Timeout)
	defer cancel()
	// Extract song id from request
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			utils.Response{Message: fmt.Sprintf("Invalid song id %s", c.Param("id"))})
	}
	// Fetch song from db
	song, err := sc.SongService.GetSong(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(
				http.StatusNotFound,
				utils.Response{Message: err.Error()})
		}
		return c.JSON(
			http.StatusInternalServerError,
			utils.Response{Message: err.Error()})
	}
	return c.JSON(
		http.StatusOK,
		utils.Response{Message: "Song received", Data: song})
}

// @Summary      Get songs with optional filtering and pagination
// @Description  Retrieve a list of songs with optional filters such as group, song name, and date range, and supports pagination with page and limit parameters.
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        group   query     string  false  "Filter by group/artist name"
// @Param        song    query     string  false  "Filter by song name"
// @Param        after   query     string  false  "Filter by songs released after date (dd.mm.yyyy)"
// @Param        before  query     string  false  "Filter by songs released before date (dd.mm.yyyy)"
// @Param        page    query     int     false  "Page number for pagination, default 1"
// @Param        limit   query     int     false  "Limit per page, default 10"
// @Success      200  {object}  utils.Response{message=string, data=[]models.Song} "Songs received"
// @Failure      400  {object}  utils.Response{message=string} "Error while parsing query params"
// @Failure      500  {object}  utils.Response{message=string} "Internal server error"
// @Router       /songs [get]
func (sc *SongController) GetSongs(c echo.Context) error {
	// New context with timeout
	ctx, cancel := context.WithTimeout(c.Request().Context(), sc.Timeout)
	defer cancel()
	// Parse query params
	query := c.Request().URL.Query()
	f, p, l, err := repository.ParseQuery(query)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			utils.Response{Message: "Error while parsing query params: " + err.Error()})
	}
	// Retrieve filtered songs from db
	songs, err := sc.SongService.GetSongs(ctx, *f, p, l)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			utils.Response{Message: err.Error()})
	}
	return c.JSON(
		http.StatusOK,
		utils.Response{Message: "Songs received", Data: songs})
}

// @Summary      Partially update a song
// @Description  Update one or more fields of an existing song by providing the song ID and the fields to update.
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        id    path     int  true  "Song ID"
// @Param        song  body     utils.SongPatchRequest  true  "Fields to update"
// @Success      200  {object}  utils.Response{message=string, data=models.Song} "Song updated"
// @Failure      400  {object}  utils.Response{message=string} "Invalid song ID or request"
// @Failure      404  {object}  utils.Response{message=string} "Song not found"
// @Failure      500  {object}  utils.Response{message=string} "Internal server error"
// @Router       /songs/{id} [patch]
func (sc *SongController) PatchSong(c echo.Context) error {
	// New context with timeout
	ctx, cancel := context.WithTimeout(c.Request().Context(), sc.Timeout)
	defer cancel()
	// Extract song id from request
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			utils.Response{Message: fmt.Sprintf("Invalid song id %s", c.Param("id"))})
	}
	// Extract song details from request
	sReq := new(utils.SongPatchRequest)
	if err := c.Bind(sReq); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			utils.Response{Message: err.Error()})
	}
	if err := validator.New().Struct(sReq); err != nil {
		var errors []string
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("%s: %s", e.Field(), e.Tag()))
		}

		return c.JSON(
			http.StatusBadRequest,
			utils.Response{Message: "invalid request", Data: errors})
	}
	// Get original song from db
	song, err := sc.SongService.GetSong(ctx, id)
	if err != nil {
		log.Logger.Error().Msgf(err.Error())
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(
				http.StatusNotFound,
				utils.Response{Message: err.Error()})
		}
		return c.JSON(
			http.StatusInternalServerError,
			utils.Response{Message: err.Error()})
	}
	// Update song in db
	newSong := &models.Song{
		ID:          &id,
		Artist:      sReq.Group,
		Name:        sReq.Song,
		Lyrics:      sReq.Lyrics,
		URL:         sReq.URL,
		ReleaseDate: sReq.ReleaseDate,
	}
	updatedSong, err := sc.SongService.UpdateSong(ctx, song, newSong)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error()})
	}
	return c.JSON(
		http.StatusOK,
		utils.Response{Message: "Song updated", Data: updatedSong})
}

// @Summary      Fully update a song or create a new one
// @Description  Replace an existing song by providing the song ID and the full song data. If the song doesn't exist, create a new one.
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        id    path     int  true  "Song ID"
// @Param        song  body     utils.SongPutRequest  true  "Full song details"
// @Success      200  {object}  utils.Response{message=string, data=models.Song} "Song updated"
// @Success      201  {object}  utils.Response{message=string, data=models.Song} "Song created"
// @Failure      400  {object}  utils.Response{message=string} "Invalid song ID or request"
// @Failure      409  {object}  utils.Response{message=string} "Song already exists"
// @Failure      500  {object}  utils.Response{message=string} "Internal server error"
// @Router       /songs/{id} [put]
func (sc *SongController) PutSong(c echo.Context) error {
	// New context with timeout
	ctx, cancel := context.WithTimeout(c.Request().Context(), sc.Timeout)
	defer cancel()
	// Extract song id from request
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			utils.Response{Message: fmt.Sprintf("Invalid song id %s", c.Param("id"))})
	}
	// Extract song details from request
	sReq := new(utils.SongPutRequest)
	if err := c.Bind(sReq); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			utils.Response{Message: err.Error()})
	}
	log.Logger.Debug().Msgf("SongPutRequest date: %s", sReq.ReleaseDate.Format("02.01.2006"))
	if err := validator.New().Struct(sReq); err != nil {
		var errors []string
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("%s: %s", e.Field(), e.Tag()))
		}

		return c.JSON(
			http.StatusBadRequest,
			utils.Response{Message: "invalid request", Data: errors})
	}
	newSong := &models.Song{
		Artist:      sReq.Group,
		Name:        sReq.Song,
		Lyrics:      sReq.Lyrics,
		URL:         sReq.URL,
		ReleaseDate: sReq.ReleaseDate,
	}
	// Get original song from db
	song, err := sc.SongService.GetSong(ctx, id)
	if err != nil {
		// Create new song in db
		if err := sc.SongService.CreateSong(ctx, newSong); err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				return c.JSON(
					http.StatusConflict,
					utils.Response{Message: err.Error()})
			}
			return c.JSON(
				http.StatusInternalServerError,
				utils.Response{Message: err.Error()})
		}
		return c.JSON(
			http.StatusCreated,
			utils.Response{Message: "Song created", Data: newSong})
	} else {
		// Update newSong in db
		newSong.ID = &id
		updatedSong, err := sc.SongService.UpdateSong(ctx, song, newSong)
		if err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				utils.Response{Message: err.Error()})
		}
		return c.JSON(
			http.StatusOK,
			utils.Response{Message: "Song updated", Data: updatedSong})
	}
}

// @Summary      Delete a song by ID
// @Description  Remove song from database
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        id   path     int  true  "Song ID"
// @Success      200  {object}  utils.Response{message=string} "Song deleted"
// @Failure      400  {object}  utils.Response{message=string} "Invalid song ID"
// @Failure      404  {object}  utils.Response{message=string} "Song not found"
// @Failure      500  {object}  utils.Response{message=string} "Internal server error"
// @Router       /songs/{id} [delete]
func (sc *SongController) DeleteSong(c echo.Context) error {
	// New context with timeout
	ctx, cancel := context.WithTimeout(c.Request().Context(), sc.Timeout)
	defer cancel()
	// Extract song id from request
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			utils.Response{Message: fmt.Sprintf("Invalid song id %s", c.Param("id"))})
	}
	// Delete song from db
	if err := sc.SongService.DeleteSong(ctx, id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(
				http.StatusNotFound,
				utils.Response{Message: err.Error()})
		}
		return c.JSON(
			http.StatusInternalServerError,
			utils.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, utils.Response{Message: "Song deleted"})
}
