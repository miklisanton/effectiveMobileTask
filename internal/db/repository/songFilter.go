package repository

type SongFilter struct {
    Name        string `db:"name"`
    Artist      string `db:"artist"`
    After       string `db:"after"`     // Song released after this date, inclusive 
    Before      string `db:"before"`    // Song released before this date, inclusive
}
