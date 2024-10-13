package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type CustomDate time.Time

type SongPostRequest struct {
    Group string `json:"group" validate:"required"`
    Song   string `json:"song" validate:"required"`
}

type SongPatchRequest struct {
    Group      string  `json:"group"`
    Song        string  `json:"song"`
    Lyrics      string  `json:"lyrics"`
    ReleaseDate CustomDate  `json:"release_date"`
    URL         string  `json:"url"`
}

type SongPutRequest struct {
    Group      string  `json:"group" validate:"required"`
    Song        string  `json:"song" validate:"required"`
    Lyrics      string  `json:"lyrics" validate:"required"`
    ReleaseDate CustomDate   `json:"release_date" validate:"required"`
    URL         string  `json:"url" validate:"required"`
}

func (j *CustomDate) UnmarshalJSON(b []byte) error {
    log.Logger.Debug().Msgf("UnmarshalJSON: %v", string(b))
    s := strings.Trim(string(b), `"`)
    t, err := time.Parse("02.01.2006", s)
    if err != nil {
        return err
    }
    *j = CustomDate(t)
    return nil
}

func (j CustomDate) MarshalJSON() ([]byte, error) {
    return json.Marshal(j.Format("02.01.2006"))
}

func (j CustomDate) Format(s string) string {
    t := time.Time(j)
    return t.Format(s)
}

func (j *CustomDate) Scan(value interface{}) error {
    if value == nil {
        *j = CustomDate(time.Time{})
        return nil
    }

    switch v := value.(type) {
    case time.Time:
        *j = CustomDate(v)
    default:
        return fmt.Errorf("cannot scan type %T into CustomDate", value)
    }
    return nil
}

func (j CustomDate) Value() (driver.Value, error) {
    t := time.Time(j)
    if t.IsZero() {
        return nil, nil
    }
    return t, nil
}


