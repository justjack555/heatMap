package config

import (
	"os"
	"strings"
	"io"
	"bufio"
	"fmt"
)

type LineError struct {
	message string
}

func (lineError LineError) Error() string {
	return lineError.message
}

/**
 Function to load all environment variables from config file

 Returns a string array, error tuple

 Generally don't like passing a string array by value
 but since it is of max size 2, small elements, will
 accept it in this instance for the benefit of not
 hard coding any application environment variable names as constants
 */
func LoadEnv() ([2]string, error) {
	var twitterEnvs [2]string

	// Twitter environment variables
	f, err := os.Open("config/config.yml")
	if err != nil {
		fmt.Println("LOAD_ENV: Unable to open config file. Exiting.")
		return twitterEnvs, err
	}

	// Create buffered reader
	r := bufio.NewReader(f)

	// Read each line from file
	i := 0
	for line, isPrefix, lineErr := r.ReadLine(); lineErr != io.EOF; line, isPrefix, lineErr = r.ReadLine() {
		// Error and not io.EOF
		if lineErr != nil {
			fmt.Println("LOAD_ENV: Unable to read line. Exiting.")
			return twitterEnvs, err
		}

		// Config line was too big for buffer - return
		if isPrefix {
			fmt.Println("LOAD_ENV: Unable to read line into single buffer. Exiting.")

			// Should return a custom error, as this will likely be nil
			return twitterEnvs, LineError{"Unable to read line into single buffer."}
		}

		// Splice the line and store KV pair in environ variable
		pair := strings.Split(string(line), "=")
		if len(pair) != 2 {
			fmt.Println("LOAD_ENV: Unable to split line around equal sign. Exiting.")
			return twitterEnvs, LineError{"Unable to split line around equal sign."}
		}

//		fmt.Println("LOAD_ENV: Split kv pair is: ", pair[0], " : ", pair[1])
		os.Setenv(pair[0], pair[1])
		twitterEnvs[i] = pair[0]
		i++

//		fmt.Println("LOAD_ENV: Key ", pair[0], " has value: ", os.Getenv(pair[0]))
	}

	return twitterEnvs, nil
}
