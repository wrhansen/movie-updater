package notion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type NotionResponse struct {
	Results    []*NotionMovie `json:"results"`
	NextCursor string         `json:"next_cursor"`
	HasMore    bool           `json:"has_more"`
	Type       string         `json:"type"`
	RequestId  string         `json:"request_id"`
	Object     string         `json:"object"`
}

type NotionMovie struct {
	Id         string     `json:"id"`
	Properties Properties `json:"properties"`
}

type Properties struct {
	MovieId     MovieId     `json:"MovieID"`
	Title       Title       `json:"Title"`
	MovieCover  Image       `json:"MovieCover"`
	Year        Year        `json:"Year"`
	Rating      Rating      `json:"Rating"`
	LatestWatch LatestWatch `json:"LatestWatch"`
	WatchCount  WatchCount  `json:"WatchCount"`
}

type Image struct {
	Files []File `json:"files"`
}
type File struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	External External `json:"external"`
}

type External struct {
	Url string `json:"url"`
}

type MovieId struct {
	RichText []RichText `json:"rich_text"`
}

type Title struct {
	Title []RichText `json:"title"`
}

type Year struct {
	Number int `json:"number"`
}

type Rating struct {
	Select Select `json:"select"`
}

type LatestWatch struct {
	Date Date `json:"date"`
}

type WatchCount struct {
	Number int `json:"number"`
}

type Date struct {
	Start string `json:"start"`
}

type Select struct {
	Name string `json:"name"`
}

type RichText struct {
	PlainText string `json:"plain_text"`
}

func (nm *NotionMovie) String() string {
	return fmt.Sprintf(
		"%s(%s): %s",
		nm.Properties.Title.Title[0].PlainText,
		strconv.Itoa(nm.Properties.Year.Number),
		nm.Properties.Rating.Select.Name,
	)
}

func GetMoviesFromNotionDatabase(apiKey string, dbId string, notionVersion string) ([]*NotionMovie, error) {
	url := fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", dbId)

	// Create a JSON body for the POST request
	body := map[string]any{
		"page_size": 100,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notion-Version", notionVersion)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read and print the JSON response as a string
	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(responseBytes))

	var responseData NotionResponse
	if err := json.Unmarshal(responseBytes, &responseData); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", responseData.Results)
	return responseData.Results, nil
}

func AddMoviesToNotionDatabase(movies []*NotionMovie, dbId string, apiKey string, notionVersion string) error {
	url := "https://api.notion.com/v1/pages"

	// Create a JSON body for the POST request
	for _, movie := range movies {
		body := map[string]any{
			"parent": map[string]any{"database_id": dbId},
			"properties": map[string]any{
				"Title": map[string]any{
					"title": []map[string]any{
						{
							"type": "text",
							"text": map[string]any{
								"content": movie.Properties.Title.Title[0].PlainText,
								"link":    nil,
							},
						},
					},
				},
				"MovieID": map[string]any{
					"rich_text": []map[string]any{
						{
							"type": "text",
							"text": map[string]any{"content": movie.Properties.MovieId.RichText[0].PlainText, "link": nil},
						},
					},
				},
				"Year": map[string]any{
					"number": movie.Properties.Year.Number,
				},
				"Rating": map[string]any{
					"select": map[string]any{
						"name": movie.Properties.Rating.Select.Name,
					},
				},
				"LatestWatch": map[string]any{
					"date": map[string]any{
						"start": movie.Properties.LatestWatch.Date.Start,
					},
				},
				"WatchCount": map[string]any{
					"number": movie.Properties.WatchCount.Number,
				},
				"MovieCover": map[string]any{
					"files": []map[string]any{
						{
							"type": "external",
							"name": "Movie Cover",
							"external": map[string]any{
								"url": movie.Properties.MovieCover.Files[0].External.Url,
							},
						},
					},
				},
			},
		}
		jsonBody, err := json.Marshal(body)
		fmt.Printf("jsonBody: %s\n", jsonBody)
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			return err
		}

		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Notion-Version", notionVersion)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode >= 400 {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			fmt.Printf("bodyBytes: %s\n", bodyBytes)
		}
		defer resp.Body.Close()
	}
	return nil
}

func UpdateWatchedDatesInNotionDatabase(movies []*NotionMovie) error {
	return nil
}
