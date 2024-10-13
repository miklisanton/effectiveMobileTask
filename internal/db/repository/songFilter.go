package repository

import (
	"fmt"
	"music-lib/internal/utils"
	"net/url"
	"strconv"
	"time"
)

type SongFilter struct {
	Name   string           `db:"name"`
	Artist string           `db:"artist"`
	After  utils.CustomDate `db:"after"`  // Song released after this date, inclusive
	Before utils.CustomDate `db:"before"` // Song released before this date, inclusive
}

// ParseQuery parses query parameters and returns SongFilter, page and limit
// query parameters: group, song, after, before, page, limit
// page and limit are used for pagination, default values are 1 and 10 respectively
func ParseQuery(query url.Values) (*SongFilter, int, int, error) {
	f := SongFilter{}
	// Default values
	page := 1
	limit := 10
	// Parse query parameters
	for key, value := range query {
		switch key {
		case "group":
			f.Artist = value[0]
		case "song":
			f.Name = value[0]
		case "after":
			t, err := time.Parse("02.01.2006", value[0])
			if err != nil {
				return nil, 0, 0, fmt.Errorf("invalid date format: %v, must be dd.mm.yyyy", value[0])
			}
			f.After = utils.CustomDate(t)
		case "before":
			t, err := time.Parse("02.01.2006", value[0])
			if err != nil {
				return nil, 0, 0, fmt.Errorf("invalid date format: %v, must be dd.mm.yyyy", value[0])
			}
			f.Before = utils.CustomDate(t)
		case "page":
			page, err := strconv.Atoi(value[0])
			if err != nil || page < 1 {
				return nil, 0, 0, fmt.Errorf("invalid page number: %v", value[0])
			}
		case "limit":
			limit, err := strconv.Atoi(value[0])
			if err != nil || limit < 1 {
				return nil, 0, 0, fmt.Errorf("invalid limit: %v", value[0])
			}
		default:
			return nil, 0, 0, fmt.Errorf("invalid query parameter: %s", key)
		}
	}
	return &f, page, limit, nil
}
