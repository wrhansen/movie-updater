# Movie Updater

This app parses letterboxd public RSS activity feed which shows all the most recent
movies you watched and the rating you gave it. After downloading the feed, it
adds/updates the movies to a personal Notion database.

This is a personal effort to keep a copy of all my data in a personal management
system (Notion).

The github action is setup to build and run the app once a day to stay on top
of keeping my notion database up to date with data from the RSS Feed.

# Environment Variables

* `LETTERBOXD_USERNAME`:  Your username in letterboxd to download the RSS feed from
* `NOTION_API_KEY`:  The API key to your Notion database integration
* `NOTION_DATABASE_ID`: The pre-created notion database ID to add the movie data to
* `NOTION_VERSION`: Version of the notion API to use (latest at this time is 2022-06-28)
