package media

import (
	"fmt"
	"strings"
	"encoding/base64"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"

	// Imports the Google Cloud Natural Language API client package.
	language "cloud.google.com/go/language/apiv1"
	"golang.org/x/net/context"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
//	"time"
)

// Custom unsuccessful HTTP request error type
type HTTPError struct {
	status string
}

type BearerError struct {
	message string
}
// Custom structure to hold token response
type BearerToken struct {
	Token_type string
	Access_token string
}

type UserType struct {
	Name string
	Screen_name string
	Location string
}

// Tweet structure
type Tweet struct {
//	Created_at time.Time
	User UserType
	Text string
//	Geo string
//	Coordinates string
	Retweet_count int
	Favorite_count int
}

type Tweets struct {
	Statuses []Tweet
}

// Type to hold data for each location
type location struct {
	score float32
}

func (httpError HTTPError) Error() string {
	return httpError.status
}

func (bearerError BearerError) Error() string {
	return bearerError.message
}

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
	var tokenPayload BearerToken


	// Make custom client
	client := &http.Client{}

	// Add body
	body := url.Values{}
	body.Set("grant_type", "client_credentials")
//	fmt.Println("POST_CRED: Encoded body is ", body.Encode())

	// Build custom request object
	req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(body.Encode()))
	if err != nil {
		fmt.Println("Post_CRED: Error creating custom request. Exiting.")
		return bearerToken, err
	}

	// Bearer + apiKey
	keyEnc := base64.StdEncoding.EncodeToString([]byte(apiKey))
	n, strErr := authString.WriteString("Basic ")
	if strErr != nil  || n != len("Basic ") {
		fmt.Println("Post_CRED: Error writing Authorization to buffer. Exiting.")
		return bearerToken, strErr
	}

	n, strErr = authString.WriteString(keyEnc)
	if strErr != nil  || n != len(keyEnc) {
		fmt.Println("Post_CRED: Error writing Authorization to buffer. Exiting.")
		return bearerToken, strErr
	}

	// Add two custom headers
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", authString.String())

//	fmt.Println("POST_CRED: Body value for key grant_type: ", req.FormValue("grant_type"))
//	fmt.Println("POST_CRED: request object is: ", req)

	// Make POST request
	fmt.Println("Making Authorization request to Twitter...")
	resp, postErr := client.Do(req)
	fmt.Println("Received Authorization response from Twitter...")

	if postErr != nil {
		fmt.Println("Post_CRED: Error with POST to Auth. Exiting.")
		return bearerToken, postErr
	}

	// Handle bad response
	if resp.StatusCode != 200 {
		fmt.Println("POST_CRED: Unsuccessful POST. Status is: ", resp.Status, ". Exiting.")
		return bearerToken, HTTPError{resp.Status}
	}
//	fmt.Println("POST_CRED: Status code of post is: ", resp.Status)

	// Close response body once function exits
	defer resp.Body.Close()

	// Read body
	respBody, respErr := ioutil.ReadAll(resp.Body)
	if respErr != nil {
		fmt.Println("POST_CRED: Error reading response body. Exiting.")
		return bearerToken, respErr
	}

	// JSON decode it
	if jsonErr := json.Unmarshal(respBody, &tokenPayload); jsonErr != nil {
		fmt.Println("POST_CRED: Error JSON formatting response body. Exiting.")
		return bearerToken, jsonErr
	}

	// Ensure that token_type is bearer
	if tokenPayload.Token_type != "bearer" {
		fmt.Println("POST_CRED: Token payload type is not bearer: ", tokenPayload.Token_type, ". Exiting.")
		return bearerToken, BearerError{tokenPayload.Token_type}
	}

//	fmt.Println("POST_CRED: Response token is: ", tokenPayload.Access_token)

	return tokenPayload.Access_token, nil
}

func GetTweets(apiKey string, query string) error {
//	fmt.Println("GET_TWEETS: Returning tweets using API KEY ", apiKey, "...")
	var fullQuery, fullAuth strings.Builder
	var authErr error
	var tweets Tweets

	searchClient := &http.Client{}

	// First Authenticate Application
	if bearerToken == "" {
		bearerToken, authErr = postCred(apiKey)
	}

	if authErr != nil {
		fmt.Println("GET_TWEETS: Error Receiving bearer token. Exiting")
		return authErr
	}

	// Add base url
	n, err := fullQuery.WriteString(resourceUrl)
	if err != nil  || n != len(resourceUrl) {
		fmt.Println("GET_TWEETS: Error writing query to buffer. Exiting.")
		return err
	}

	// Query indicator
	n, err = fullQuery.WriteString(startQuery)
	if err != nil  || n != len(startQuery) {
		fmt.Println("GET_TWEETS: Error writing query to buffer. Exiting.")
		return err
	}

	// Replace spaces
	query = strings.Replace(query, " ", "%20", -1)
	n, err = fullQuery.WriteString(query)
	if err != nil  || n != len(query) {
		fmt.Println("GET_TWEETS: Error writing query to buffer. Exiting.")
		return err
	}

	// Add location to request
	n, err = fullQuery.WriteString(locationPin)
	if err != nil  || n != len(locationPin) {
		fmt.Println("GET_TWEETS: Error writing query to buffer. Exiting.")
		return err
	}

	// Obtain request object
	searchReq, reqErr := http.NewRequest("GET", fullQuery.String(), nil)
	if reqErr != nil {
		fmt.Println("GET_TWEETS: Error creating custom request. Exiting.")
		return reqErr
	}

	// Add token to header
	n, err = fullAuth.WriteString("Bearer ")
	if err != nil  || n != len("Bearer ") {
		fmt.Println("GET_TWEETS: Error writing authorization header to buffer. Exiting.")
		return err
	}
	n, err = fullAuth.WriteString(bearerToken)
	if err != nil  || n != len(bearerToken) {
		fmt.Println("GET_TWEETS: Error writing authorization header to buffer. Exiting.")
		return err
	}

	fmt.Println("GET_TWEETS: Full authorization header is: ", fullAuth.String())

	searchReq.Header.Add("Authorization", fullAuth.String())

	fmt.Println("Making search request to Twitter...")
	searchResp, getErr := searchClient.Do(searchReq)
	fmt.Println("Received search response from Twitter...")

	if getErr != nil {
		fmt.Println("GET_TWEETS: Error with GET for tweets. Exiting.")
		return getErr
	}

	// Handle bad response
	if searchResp.StatusCode != 200 {
		fmt.Println("GET_TWEETS: Unsuccessful GET. Status is: ", searchResp.Status, ". Exiting.")
		return HTTPError{searchResp.Status}
	}
	//	fmt.Println("POST_CRED: Status code of post is: ", resp.Status)

	// Close response body once function exits
	defer searchResp.Body.Close()

	// Read body
	respBody, respErr := ioutil.ReadAll(searchResp.Body)
	if respErr != nil {
		fmt.Println("GET_TWEETS: Error reading response body. Exiting.")
		return respErr
	}

//	fmt.Println("GET_TWEETS: Tweets are: ", string(respBody))

	// JSON decode it
	if jsonErr := json.Unmarshal(respBody, &tweets); jsonErr != nil {
		fmt.Println("GET_TWEETS: Error JSON formatting response body. Exiting.")
		return jsonErr
	}

	if sentErr := analyzeTweets(&tweets); sentErr != nil {
		fmt.Println("GET_TWEETS: Error analyzing tweets. Exiting.")
		return sentErr
	}
/*	for i, tweet := range tweets.Statuses {
		fmt.Println("GET_TWEETS: ", i, "th tweet is: ", tweet)
	}
*/
	return nil
}

/**
 * Function to analyze sentiment of each tweet
 * Accepts a pointer to a Tweets structure, loads
 * sentiment scores into database
 * Returns error if one occurs
 */
func analyzeTweets(tweets *Tweets) error {
	locations := make(map[string]location)

	// Create Request
	sentReq := &languagepb.AnalyzeSentimentRequest{
		Document: &languagepb.Document{
			Type: languagepb.Document_PLAIN_TEXT,
			Source: &languagepb.Document_Content{
//				Content: "",
			},
		},
		EncodingType: languagepb.EncodingType_UTF8,
	}

	ctx := context.Background()

	// Creates a client.
	client, err := language.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Analyze each tweet
	for i, tweet := range tweets.Statuses {
		src, ok := sentReq.Document.Source.(*languagepb.Document_Content)
		if !ok {
			log.Fatalf("Failed to perform type assertion on Source. Exiting")
		}

		src.Content = tweet.Text

		// Detects the sentiment of the text.
		sentiment, err := client.AnalyzeSentiment(ctx, sentReq)
		if err != nil {
			log.Fatalf("Failed to analyze text: %v", err)
		}

		fmt.Printf("%dth tweet text is: %v\n", i, tweet.Text)
		prev, ok := locations[tweet.User.Location]
		prev.score = prev.score + sentiment.DocumentSentiment.Score
		locations[tweet.User.Location] = prev
/*		if sentiment.DocumentSentiment.Score >= 0 {
			fmt.Println("Sentiment: positive with value", sentiment.DocumentSentiment.Score)
		} else {
			fmt.Println("Sentiment: negative with value", sentiment.DocumentSentiment.Score)
		}
*/	}

	return nil
}