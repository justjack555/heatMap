package media

import (
	"fmt"
	"strings"
	"encoding/base64"
	"net/url"
	"net/http"
//	"encoding/json"
)

const resourceUrl string = "https://api.twitter.com/1.1/search/tweets.json"
const startQuery string = "?q="
const locationPin string = "&geocode=39.8283,-98.5795,1500mi"

const contentType string = `application/x-www-form-urlencoded;charset=UTF-8`
const tokenEndpoint string = "https://api.twitter.com/oauth2/token"
var bearerToken string = ""

/**
 Function to obtain Bearer Token for GET request

 Makes POST request to twitter App Auth API
 */
func postCred (apiKey string)(string, error){
	var authString strings.Builder

	// Make custom client
	client := &http.Client{}

	// Add body
	body := url.Values{}
	body.Set("grant_type", "client_credentials")

	// Build custom request object
	req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(body.Encode()))
	if err != nil {
		fmt.Println("Post_CRED: Error creating custom request. Exiting.")
		return "", err
	}

	// Bearer + apiKey
	keyEnc := base64.StdEncoding.EncodeToString([]byte(apiKey))
	n, strErr := authString.WriteString("Bearer ")
	if strErr != nil  || n != len("Bearer ") {
		fmt.Println("Post_CRED: Error writing Authorization to buffer. Exiting.")
		return "", strErr
	}

	n, strErr = authString.WriteString(keyEnc)
	if strErr != nil  || n != len(keyEnc) {
		fmt.Println("Post_CRED: Error writing Authorization to buffer. Exiting.")
		return "", strErr
	}

	// Add two custom headers
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", authString.String())

	// Make POST request
	resp, postErr := client.Do(req)
	if postErr != nil {
		fmt.Println("Post_CRED: Error with POST to Auth. Exiting.")
		return "", postErr
	}

	fmt.Println("POST_CRED: Status code of post is: ", resp.Status)

	return "", nil
}

func GetTweets(apiKey string, query string) (string, error) {
	fmt.Println("GET_TWEETS: Returning tweets using API KEY ", apiKey, "...")
	var fullQuery strings.Builder
	var authErr error

	// First Authenticate Application
	bearerToken, authErr = postCred(apiKey)
	if authErr != nil {
		fmt.Println("GET_TWEETS: Error Receiving bearer token. Exiting")
		return "", authErr
	}
	// Add base url
	n, err := fullQuery.WriteString(resourceUrl)
	if err != nil  || n != len(resourceUrl) {
		fmt.Println("GET_TWEETS: Error writing query to buffer. Exiting.")
		return "", err
	}

	// Query indicator
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

	// Add location to request
	n, err = fullQuery.WriteString(locationPin)
	if err != nil  || n != len(locationPin) {
		fmt.Println("GET_TWEETS: Error writing query to buffer. Exiting.")
		return "", err
	}

/*	reqMap := map[string]string{
		"q" : fullQuery.String(),
		"geocode" : locationPin,
	}

	reqJson, err := json.Marshal(reqMap)
	if err != nil {
		fmt.Println("GET_TWEETS: Error encoding map in JSON. Exiting.")
		return "", err
	}
*/
	return fullQuery.String(), nil
}