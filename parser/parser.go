package parser

import (
	"fmt"
	"movie-updater/movie"
	"time"

	"github.com/mmcdole/gofeed"
)

func ParseMoviesFromRSSFeed(url string) ([]*movie.Movie, error) {
	// Create a new parser
	fp := gofeed.NewParser()

	// Parse the RSS feed from a URL
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("error parsing RSS feed: %v", err)
	}

	var movies = []*movie.Movie{}

	// Print each item title
	for _, item := range feed.Items {
		const format = "2006-01-02"
		parsedTime, err := time.Parse(format, item.Extensions["letterboxd"]["watchedDate"][0].Value)
		if err != nil {
			return nil, fmt.Errorf("error parsing time: %v", err)
		}

		var id string
		if movieID, ok := item.Extensions["tmdb"]["movieId"]; ok && len(movieID) > 0 {
			id = movieID[0].Value
		} else if tvID, ok := item.Extensions["tmdb"]["tvId"]; ok && len(tvID) > 0 {
			id = tvID[0].Value
		} else {
			id = ""
		}

		parsedMovie := &movie.Movie{
			Title:       item.Extensions["letterboxd"]["filmTitle"][0].Value,
			ImageUrl:    item.Image.URL,
			Rating:      item.Extensions["letterboxd"]["memberRating"][0].Value,
			MovieID:     id,
			Year:        item.Extensions["letterboxd"]["filmYear"][0].Value,
			LatestWatch: parsedTime,
		}

		movies = append(movies, parsedMovie)
	}

	return movies, nil
}
