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
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type SongController struct {
    SongService services.ISongService
    MusicInfoService services.IMusicInfoService
    Timeout time.Duration
}

func NewSongController(
        songS services.ISongService,
        musicInfoS services.IMusicInfoService,
        cfg *config.Config) *SongController {

    timeout := time.Duration(cfg.ExternalAPI.Timeout) * time.Second

    return &SongController{songS, musicInfoS, timeout}
}

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
        return c.JSON(
            http.StatusInternalServerError,
            utils.Response{Message: err.Error()})
    }
    return c.JSON(
        http.StatusCreated,
        utils.Response{Message: "Song created", Data: song})
}

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
        return c.JSON(
            http.StatusInternalServerError,
            utils.Response{Message: err.Error()})
    }
    return c.JSON(
        http.StatusOK,
        utils.Response{Message: "Song received", Data: song})
}

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
    songs ,err := sc.SongService.GetSongs(ctx, *f, p, l)
    if err != nil {
        return c.JSON(
            http.StatusInternalServerError,
            utils.Response{Message: err.Error()})
    }
    return c.JSON(
        http.StatusOK,
        utils.Response{Message:"Songs received", Data: songs})
}

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
        return c.JSON(
            http.StatusNotFound, 
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
        return c.JSON(
            http.StatusInternalServerError, 
            utils.Response{Message: err.Error()})
    }
    return c.JSON(http.StatusOK, utils.Response{Message: "Song deleted"})
}

