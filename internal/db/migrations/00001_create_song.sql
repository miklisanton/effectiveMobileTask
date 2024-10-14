-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE song (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    artist VARCHAR(255) NOT NULL,
    lyrics TEXT NOT NULL,
    release_date DATE NOT NULL,
    url VARCHAR(255) NOT NULL,
    UNIQUE (name, artist)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE song;
-- +goose StatementEnd
