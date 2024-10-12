package main

import (
	"fmt"
	"music-lib/internal/config"
	"music-lib/internal/db/drivers"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

func main() {
	e := echo.New()
	// Setup request logging middleware
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

	// Endpoints
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":" + cfg.Server.Port))
}
