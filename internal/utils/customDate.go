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

func (j *CustomDate) UnmarshalJSON(b []byte) error {
    log.Logger.Debug().Msgf("UnmarshalJSON: %v", string(b))
    s := strings.Trim(string(b), `"`)
    t, err := time.Parse("02.01.2006", s)
    if err != nil {
        return fmt.Errorf("wrong date format, need dd.mm.yyyy")
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
