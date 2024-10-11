package drivers

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func Connect(URL string) (*sqlx.DB, error) {
    db, err := sql.Open("postgres", URL)
    if err != nil {
        return nil, err
    }

    if err := db.Ping(); err != nil {
        return nil, err
    }

    if err := runMigrations(db); err != nil {
        return nil, err
    }
    return sqlx.NewDb(db, "postgres"), nil
}

func runMigrations(db *sql.DB) error {
    goose.SetDialect("postgres")
    return goose.Up(db, "./internal/db/migrations")
}

func Down(db *sqlx.DB) error {
    goose.SetDialect("postgres")
    var sqlDB *sql.DB = db.DB
    return goose.Down(sqlDB, "./internal/db/migrations")
}

