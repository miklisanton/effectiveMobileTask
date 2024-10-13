package utils

import (
    "fmt"
    "net/url"
    "strconv"
    "time"
    "music-lib/internal/db/repository"
)

// ParseQuery parses query parameters and returns SongFilter, page and limit
// query parameters: artist, name, after, before, page, limit
// page and limit are used for pagination, default values are 1 and 10 respectively
func ParseQuery(query url.Values) (*repository.SongFilter, int, int, error) {
    f := repository.SongFilter{}
    // Default values
    page := 1
    limit := 10
    // Parse query parameters
    for key, value := range query {
        switch key {
            case "artist":
                f.Artist = value[0]
            case "name":
                f.Name = value[0]
            case "after":
                _, err := time.Parse("2006-01-02", value[0])
                if err != nil {
                    return nil, 0, 0, fmt.Errorf("invalid date format: %v, must be yyyy-mm-dd", value[0])
                }
                f.After = value[0]
            case "before":
                _, err := time.Parse("2006-01-02", value[0])
                if err != nil {
                    return nil, 0, 0, fmt.Errorf("invalid date format: %v, must be yyyy-mm-dd", value[0])
                }
                f.Before = value[0]
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
