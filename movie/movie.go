package movie

import (
	"fmt"
	"time"
)

type Movie struct {
	Title       string
	ImageUrl    string
	Rating      string
	MovieID     string
	Year        string
	LatestWatch time.Time
	WatchCount  int
}

func (m *Movie) String() string {
	return fmt.Sprintf("%s: %s (%s) -- %s", m.MovieID, m.Title, m.Year, m.Rating)
}
