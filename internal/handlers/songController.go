package handlers

import (
	"context"
	"fmt"
	"music-lib/internal/config"
	"music-lib/internal/db/models"
	"music-lib/internal/services"
	"music-lib/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type SongController struct {
    SongService services.ISongService
    MusicInfoService services.IMusicInfoService
    Timeout time.Duration
}

func NewSongController(
        songS services.ISongService,
        musicInfoS services.IMusicInfoService,
        cfg config.Config) *SongController {

    timeout := time.Duration(cfg.ExternalAPI.Timeout) * time.Second

    return &SongController{songS, musicInfoS, timeout}
}

func (sc *SongController) CreateSong(c echo.Context) error {
    // New context with timeout
    ctx, cancel := context.WithTimeout(c.Request().Context(), sc.Timeout)
    defer cancel()
    // Extract song name and artist from request
    songRequest := new(utils.SongRequest)
    if err := c.Bind(songRequest); err != nil {
        return c.JSON(http.StatusBadRequest, utils.Response{Message: err.Error()})
    }
    if err := validator.New().Struct(songRequest); err != nil {
        return c.JSON(http.StatusBadRequest, utils.Response{Message: err.Error()})
    }
    // Fetch song details from external service
    songDetail, err := sc.MusicInfoService.GetSongInfo(songRequest.Artist, songRequest.Name)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error()})
    }
    // Save song to db
    song := &models.Song{
        Artist:      songRequest.Artist,
        Name:        songRequest.Name,
        Lyrics:      songDetail.Text,
        URL:         songDetail.Link,
        ReleaseDate: songDetail.ReleaseDate,
    }
    if err := sc.SongService.CreateSong(ctx, song); err != nil {
        return c.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error()})
    }
    return c.JSON(http.StatusCreated, utils.Response{Message: "Song created", Data: song})
}

func (sc *SongController) GetSong(c echo.Context) error {
    // New context with timeout
    ctx, cancel := context.WithTimeout(c.Request().Context(), sc.Timeout)
    defer cancel()
    // Extract song id from request
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, utils.Response{Message: fmt.Sprintf("Invalid song id %s", c.Param("id"))})
    }
    // Fetch song from db
    song, err := sc.SongService.GetSong(ctx, id)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error()})
    }
    return c.JSON(http.StatusOK, utils.Response{Data: song})
}

func (sc *SongController) GetSongs(c echo.Context) error {
    // New context with timeout
    ctx, cancel := context.WithTimeout(c.Request().Context(), sc.Timeout)
    defer cancel()
    // Parse query params
    query := c.Request().URL.Query()
    f, p, l, err := utils.ParseQuery(query)
    if err != nil {
        return c.JSON(http.StatusBadRequest, utils.Response{Message: err.Error()})
    }
    // Retrieve filtered songs from db
    songs ,err := sc.SongService.GetSongs(ctx, *f, p, l)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error()})
    }
    return c.JSON(http.StatusOK, utils.Response{Data: songs})
}

func (sc *SongController) PatchSong(c echo.Context) error {
    // New context with timeout
    ctx, cancel := context.WithTimeout(c.Request().Context(), sc.Timeout)
    defer cancel()
    // Extract song id from request
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, utils.Response{Message: fmt.Sprintf("Invalid song id %s", c.Param("id"))})
    }
    // Extract song details from request
    sReq := new(utils.SongPatchRequest)
    if err := c.Bind(sReq); err != nil {
        return c.JSON(http.StatusBadRequest, utils.Response{Message: err.Error()})
    }
    if err := validator.New().Struct(sReq); err != nil {
        return c.JSON(http.StatusBadRequest, utils.Response{Message: err.Error()})
    }
    // Update song in db
    song := &models.Song{
        Artist:      sReq.Artist,
        Name:        sReq.Name,
        Lyrics:      sReq.Lyrics,
        URL:         sReq.URL,
        ReleaseDate: sReq.ReleaseDate,
    }
    updatedSong, err := sc.SongService.UpdateSong(ctx, id, song)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error()})
    }
    return c.JSON(http.StatusOK, utils.Response{Message: "Song updated", Data: updatedSong})
}

func (sc *SongController) PutSong(c echo.Context) error {
    // New context with timeout
    ctx, cancel := context.WithTimeout(c.Request().Context(), sc.Timeout)
    defer cancel()
    // Extract song id from request
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, utils.Response{Message: fmt.Sprintf("Invalid song id %s", c.Param("id"))})
    }
    // Extract song details from request
    sReq := new(utils.SongPutRequest)
    if err := c.Bind(sReq); err != nil {
        return c.JSON(http.StatusBadRequest, utils.Response{Message: err.Error()})
    }
    if err := validator.New().Struct(sReq); err != nil {
        return c.JSON(http.StatusBadRequest, utils.Response{Message: err.Error()})
    }
    // Update song in db
    song := &models.Song{
        Artist:      sReq.Artist,
        Name:        sReq.Name,
        Lyrics:      sReq.Lyrics,
        URL:         sReq.URL,
        ReleaseDate: sReq.ReleaseDate,
    }
    updatedSong, err := sc.SongService.UpdateSong(ctx, id, song)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error()})
    }
    return c.JSON(http.StatusOK, utils.Response{Message: "Song updated", Data: updatedSong})
}

func (sc *SongController) DeleteSong(c echo.Context) error {
    // New context with timeout
    ctx, cancel := context.WithTimeout(c.Request().Context(), sc.Timeout)
    defer cancel()
    // Extract song id from request
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, utils.Response{Message: fmt.Sprintf("Invalid song id %s", c.Param("id"))})
    }
    // Delete song from db
    if err := sc.SongService.DeleteSong(ctx, id); err != nil {
        return c.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error()})
    }
    return c.JSON(http.StatusOK, utils.Response{Message: "Song deleted"})
}

