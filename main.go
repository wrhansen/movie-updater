package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"movie-updater/movie"
	"movie-updater/notion"
	"movie-updater/parser"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v\n", err)
	}

	// Get Notion Database ID from environment variable
	dbId := os.Getenv("NOTION_DATABASE_ID")
	if dbId == "" {
		fmt.Println("NOTION_DATABASE_ID environment variable is not set")
		return
	}

	notionApiKey := os.Getenv("NOTION_API_KEY")
	if notionApiKey == "" {
		fmt.Println("NOTION_API_KEY environment variable is not set")
		return
	}

	// 2022-06-28
	notionVersion := os.Getenv("NOTION_VERSION")
	if notionVersion == "" {
		fmt.Println("NOTION_VERSION environment variable is not set")
		return
	}

	letterboxdUsername := os.Getenv("LETTERBOXD_USERNAME")
	if letterboxdUsername == "" {
		fmt.Println("LETTERBOXD_USERNAME environment variable is not set")
		return
	}

	// Get Movies from Letterboxd RSS Feed
	feedMovies, err := parser.ParseMoviesFromRSSFeed(
		fmt.Sprintf("https://letterboxd.com/%s/rss/", letterboxdUsername))
	if err != nil {
		fmt.Println("Error parsing RSS feed:", err)
		return
	}
	fmt.Printf("%v\n", feedMovies)

	// Get Movies from Notion Database
	notionMovies, err := notion.GetMoviesFromNotionDatabase(notionApiKey, dbId, notionVersion)
	if err != nil {
		fmt.Println("Error getting movies from Notion Database:", err)
		return
	}

	fmt.Printf("notionMovies: %v\n", notionMovies)

	// Figure out which movies to add
	newMovies := determineNewMovies(feedMovies, notionMovies)
	fmt.Printf("newMovies: %v\n", newMovies)

	// Figure out which movies have updates
	updateMovies := determineUpdatedMovies(feedMovies, notionMovies)
	fmt.Printf("updatedMovies: %v\n", updateMovies)

	// Add movies to Notion Database
	if len(newMovies) > 0 {
		err = notion.AddMoviesToNotionDatabase(newMovies, dbId, notionApiKey, notionVersion)
		if err != nil {
			fmt.Println("Error adding movies to Notion Database:", err)
			return
		}
	} else {
		fmt.Println("No New Movies to Add...")
	}

	// TODO: Finish implementing updates.
	// Update watched dates in Notion Database
	// if len(updateMovies) > 0 {
	// 	err = notion.UpdateWatchedDatesInNotionDatabase(updateMovies)
	// 	if err != nil {
	// 		fmt.Println("Error updating watched dates in Notion Database:", err)
	// 		return
	// 	}
	// }
}

func determineNewMovies(feed_movies []*movie.Movie, notionMovies []*notion.NotionMovie) []*notion.NotionMovie {
	// Create a map of movie IDs to NotionMovie structs
	notionMap := make(map[string]*notion.NotionMovie)
	for _, nm := range notionMovies {
		notionMap[nm.Properties.MovieId.RichText[0].PlainText] = nm
	}

	// Create a slice of movies that are in the feed but not in the Notion database
	var newMovies []*notion.NotionMovie
	for _, fm := range feed_movies {
		if _, ok := notionMap[fm.MovieID]; !ok {
			parsedYear, err := strconv.Atoi(fm.Year)
			if err != nil {
				fmt.Println("Error converting year to int:", err)
				parsedYear = 0
			}
			notionMovie := &notion.NotionMovie{
				Properties: notion.Properties{
					MovieId: notion.MovieId{
						RichText: []notion.RichText{
							{PlainText: fm.MovieID},
						},
					},
					Title: notion.Title{
						Title: []notion.RichText{
							{PlainText: fm.Title},
						},
					},
					Year: notion.Year{
						Number: parsedYear,
					},
					Rating: notion.Rating{
						Select: notion.Select{
							Name: fm.Rating,
						},
					},
					LatestWatch: notion.LatestWatch{
						Date: notion.Date{
							Start: fm.LatestWatch.Format("2006-01-02T00:00:00Z"),
						},
					},
					WatchCount: notion.WatchCount{
						Number: 1,
					},
					MovieCover: notion.Image{
						Files: []notion.File{
							{
								Name: "Movie Cover",
								Type: "external",
								External: notion.External{
									Url: fm.ImageUrl,
								},
							},
						},
					},
				},
			}
			newMovies = append(newMovies, notionMovie)
		}
	}

	return newMovies
}

func determineUpdatedMovies(feed_movies []*movie.Movie, notionMovies []*notion.NotionMovie) []*notion.NotionMovie {
	// Create a map of movie IDs to NotionMovie structs
	notionMap := make(map[string]*notion.NotionMovie)
	for _, nm := range notionMovies {
		notionMap[nm.Properties.MovieId.RichText[0].PlainText] = nm
	}

	// Create a slice of movies that are in both the feed and the Notion database
	var updatedMovies []*notion.NotionMovie
	for _, fm := range feed_movies {
		if nm, ok := notionMap[fm.MovieID]; ok {
			const layout = "2006-01-02T00:00:00.000+00:00"
			parsedTime, err := time.Parse(layout, nm.Properties.LatestWatch.Date.Start)
			if err != nil {
				fmt.Println("Error parsing time:", err)
				continue
			}

			// Update LatestWatch and WatchCount
			if fm.LatestWatch.After(parsedTime) {
				nm.Properties.LatestWatch.Date.Start = fm.LatestWatch.Format(layout)
				nm.Properties.WatchCount.Number++
				updatedMovies = append(updatedMovies, nm)
			}

			// Update image (if changed)
			if fm.ImageUrl != nm.Properties.MovieCover.Files[0].External.Url {
				nm.Properties.MovieCover.Files[0].Type = "external"
				nm.Properties.MovieCover.Files[0].Name = "Movie Cover"
				nm.Properties.MovieCover.Files[0].External.Url = fm.ImageUrl
			}
		}
	}
	return updatedMovies
}
