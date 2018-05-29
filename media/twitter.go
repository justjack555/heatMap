package media

import (
	"fmt"
	"strings"
	"encoding/json"
)

const resourceUrl string = "https://api.twitter.com/1.1/search/tweets.json"
const startQuery string = "?q="
const locationPin string = "39.8283,-98.5795,1500mi"

func GetTweets(apiKey string, query string) (string, error) {
	fmt.Println("GET_TWEETS: Returning tweets...")
	var fullQuery strings.Builder
	n, err := fullQuery.WriteString(resourceUrl)
	if err != nil  || n != len(resourceUrl) {
		fmt.Println("GET_TWEETS: Error writing query to buffer. Exiting.")
		return "", err
	}

	n, err = fullQuery.WriteString(startQuery)
	if err != nil  || n != len(startQuery) {
		fmt.Println("GET_TWEETS: Error writing query to buffer. Exiting.")
		return "", err
	}

	// Replace spaces
	query = strings.Replace(query, " ", "%20", -1)
	n, err = fullQuery.WriteString(query)
	if err != nil  || n != len(query) {
		fmt.Println("GET_TWEETS: Error writing query to buffer. Exiting.")
		return "", err
	}

	reqMap := map[string]string{
		"q" : fullQuery.String(),
		"geocode" : locationPin,
	}

	reqJson, err := json.Marshal(reqMap)
	if err != nil {
		fmt.Println("GET_TWEETS: Error encoding map in JSON. Exiting.")
		return "", err
	}

	return string(reqJson), nil
}