package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/justjack555/heatMap/config"
	"github.com/justjack555/heatMap/media"
)
var twitterEnv [2]string

func invokeTwitter() error {
	var twitterCred strings.Builder
	query := os.Args[1]

	// Concatenate key, colon, and secret
	for i, str := range twitterEnv {
//		fmt.Println("INVOKE_TWITTER: Index value ", i, ", arg value: ", str)
		twitterCred.WriteString(os.Getenv(str))

		// Append colon after consumer key
		if i == 0 {
			twitterCred.WriteString(":")
		}
	}

//	fmt.Println("INVOKE_TWITTER: Final twitter credential: ", twitterCred.String())

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
	var err error
	if len(os.Args[1:]) != 1 {
		fmt.Println("Usage: go run heatMap.go <QUERY>")
		return
	}

	// Load application specific environment variables
	twitterEnv, err = config.LoadEnv()
	if err != nil {
		fmt.Println("Error message is: ", err)
		return
	}

	// Invoke twitter handler to retrieve tweets
	err = invokeTwitter()
	if err != nil {
		fmt.Println("Error message is: ", err)
		return
	}

}