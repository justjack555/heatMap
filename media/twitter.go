package media

import (
	"fmt"
	"strings"
	"encoding/base64"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
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

func GetTweets(apiKey string, query string) (string, error) {
//	fmt.Println("GET_TWEETS: Returning tweets using API KEY ", apiKey, "...")
	var fullQuery, fullAuth strings.Builder
	var authErr error
	searchClient := &http.Client{}

	// First Authenticate Application
	if bearerToken == "" {
		bearerToken, authErr = postCred(apiKey)
	}

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

	// Obtain request object
	searchReq, reqErr := http.NewRequest("GET", fullQuery.String(), nil)
	if reqErr != nil {
		fmt.Println("GET_TWEETS: Error creating custom request. Exiting.")
		return bearerToken, reqErr
	}

	// Add token to header
	n, err = fullAuth.WriteString("Bearer ")
	if err != nil  || n != len("Bearer ") {
		fmt.Println("GET_TWEETS: Error writing authorization header to buffer. Exiting.")
		return "", err
	}
	n, err = fullAuth.WriteString(bearerToken)
	if err != nil  || n != len(bearerToken) {
		fmt.Println("GET_TWEETS: Error writing authorization header to buffer. Exiting.")
		return "", err
	}
/*	n, err = fullAuth.WriteString(". Signing is not required.")
	if err != nil  || n != len(". Signing is not required.") {
		fmt.Println("GET_TWEETS: Error writing authorization header to buffer. Exiting.")
		return "", err
	}
*/
	fmt.Println("GET_TWEETS: Full authorization header is: ", fullAuth.String())

	searchReq.Header.Add("Authorization", fullAuth.String())

	fmt.Println("Making search request to Twitter...")
	searchResp, getErr := searchClient.Do(searchReq)
	fmt.Println("Received search response from Twitter...")

	if getErr != nil {
		fmt.Println("GET_TWEETS: Error with GET for tweets. Exiting.")
		return "", getErr
	}

	// Handle bad response
	if searchResp.StatusCode != 200 {
		fmt.Println("GET_TWEETS: Unsuccessful GET. Status is: ", searchResp.Status, ". Exiting.")
		return "", HTTPError{searchResp.Status}
	}
	//	fmt.Println("POST_CRED: Status code of post is: ", resp.Status)

	// Close response body once function exits
	defer searchResp.Body.Close()

	// Read body
	respBody, respErr := ioutil.ReadAll(searchResp.Body)
	if respErr != nil {
		fmt.Println("GET_TWEETS: Error reading response body. Exiting.")
		return "", respErr
	}

	fmt.Println("GET_TWEETS: Tweets are: ", string(respBody))

	// JSON decode it
/*	if jsonErr := json.Unmarshal(respBody, &tweetsPayload); jsonErr != nil {
		fmt.Println("GET_TWEETS: Error JSON formatting response body. Exiting.")
		return "", jsonErr
	}
*/


	// Do request
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