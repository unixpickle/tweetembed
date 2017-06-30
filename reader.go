package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

// ReadTweets reads tweets from a file.
// If the file name is empty, stdin is used.
func ReadTweets(file string) (<-chan string, error) {
	var ioStream *os.File
	if file == "" {
		fmt.Fprintln(os.Stderr, "reading from standard input.")
		ioStream = os.Stdin
	} else {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		ioStream = f
	}
	reader := csv.NewReader(ioStream)

	resChan := make(chan string, 16)
	go func() {
		if ioStream != os.Stdin {
			defer ioStream.Close()
		}
		defer close(resChan)
		for {
			entry, err := reader.Read()
			if err != nil {
				return
			}
			resChan <- entry[len(entry)-1]
		}
	}()

	return resChan, nil
}
