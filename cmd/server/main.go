package main

import (
	"fmt"
	"music-lib/internal/config"
	"music-lib/internal/db/drivers"
	"music-lib/internal/db/repository"
	"music-lib/internal/handlers"
	"music-lib/internal/services"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/swaggo/echo-swagger"
	_ "music-lib/docs"
)

var (
	db  *sqlx.DB
	cfg *config.Config
)

func init() {
	// Setup logger
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()
	log.Info().Msg("Logger initialized")
	cfgPath, err := config.ParseCLI()
	// Parse CLI arguments
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse CLI")
	}
	// Read config
	cfg, err = config.NewConfig(cfgPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read config")
	}
	log.Info().Msg("Config loaded")
	// Initialize database
	connURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Db.User,
		cfg.Db.Password,
		cfg.Db.Host,
		cfg.Db.Port,
		cfg.Db.Name,
	)
	db, err = drivers.Connect(connURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	log.Info().Msg("Database connected")
}

// @title Songs API
// @version 1.0
// @description This is an API for managing songs library.

// @host localhost:8080
// @BasePath /api1/public
func main() {
	// Setup services
	songRepo := repository.NewSongRepository(db)
	songService := services.NewSongService(songRepo)
	musicInfoService, err := services.NewMusicInfoService(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize music info service")
	}
	// Setup controllers
	songController := handlers.NewSongController(songService, musicInfoService, cfg)
	// Setup echo
	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Logger.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Msg("request")

			return nil
		},
	}))
	pg := e.Group("/api1/public")

	// Endpoints
	pg.POST("/songs", songController.CreateSong)
	pg.GET("/songs", songController.GetSongs)
	pg.GET("/songs/:id", songController.GetSong)
	pg.PUT("/songs/:id", songController.PutSong)
	pg.PATCH("/songs/:id", songController.PatchSong)
	pg.DELETE("/songs/:id", songController.DeleteSong)
	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	// Start server
	e.Logger.Fatal(e.Start(":" + cfg.Server.Port))
}
