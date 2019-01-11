package main

import (
	"log"
	"net/http"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

// either place key here or pass by -k flag
var developerKey string = "YOUR_KEY_HERE"

const maxResults = 1
const searchType = "video"

func SetYoutubeKey(key string) {
	if key == "" {
		return
	}
	developerKey = key
}

func YoutubeSearch(query *string) string {
	result := ""

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	call := service.Search.List("id,snippet").
		Q(*query).
		MaxResults(maxResults).
		Type(searchType)
	response, err := call.Do()
	handleError(err, "")

	if response == nil || len(response.Items) == 0 {
		return result
	}

	result = response.Items[0].Id.VideoId

	return result
}
