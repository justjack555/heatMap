package main

import (
	"fmt"
	"os"

	"github.com/justjack555/heatMap/media"
)

func main(){
	if len(os.Args[1:]) != 2 {
		fmt.Println("Usage: go run heatMap.go <TWITTER_API_KEY> <QUERY>")
		return
	}

	twitterKey := os.Args[1]
	query := os.Args[2]

	twitterRes, err := media.GetTweets(twitterKey, query)
	if err != nil {
		fmt.Println("Error message is: ", err)
		return
	}

	fmt.Println("Result of Get_Tweets is: ", twitterRes)
}