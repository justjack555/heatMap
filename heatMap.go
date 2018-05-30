package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/justjack555/heatMap/media"
)

func invokeTwitter() error {
	var twitterCred strings.Builder
//	twitterKey := os.Args[1]
//	twitterSec := os.Args[2]
	query := os.Args[3]

	// Concatenate key, colon, and secret
	for i, str := range os.Args[1:3] {
		fmt.Println("INVOKE_TWITTER: Index value ", i, ", arg value: ", str)
		twitterCred.WriteString(str)

		// Append colon after consumer key
		if i == 0 {
			twitterCred.WriteString(":")
		}
	}

	fmt.Println("INVOKE_TWITTER: Final twitter credential: ", twitterCred.String())

	// With full credential, send to GetTweets
	twitterRes, err := media.GetTweets(twitterCred.String(), query)
	if err != nil {
		fmt.Println("Error message is: ", err)
		return err
	}

	fmt.Println("Result of Get_Tweets is: ", twitterRes)

	return nil
}

func main(){
	if len(os.Args[1:]) != 3 {
		fmt.Println("Usage: go run heatMap.go <TWITTER_API_KEY> <TWITTER_SECRET_KEY> <QUERY>")
		return
	}

	// Invoke twitter handler to retrieve tweets
	err := invokeTwitter()

	if err != nil {
		fmt.Println("Error message is: ", err)
		return
	}

}