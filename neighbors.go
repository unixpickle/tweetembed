package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/unixpickle/essentials"
	"github.com/unixpickle/serializer"
	"github.com/unixpickle/wordembed/glove"
)

func CmdNeighbors(args []string) {
	var embeddingFile string
	var token string
	var numMatches int
	fs := flag.NewFlagSet("neighbors", flag.ExitOnError)
	fs.StringVar(&embeddingFile, "embedding", "embedding_out", "trained embedding")
	fs.StringVar(&token, "token", "", "token to lookup")
	fs.IntVar(&numMatches, "num", 5, "number of matches")
	fs.Parse(args)

	if token == "" {
		fmt.Fprintln(os.Stderr, "Missing -token flag.")
		fmt.Fprintln(os.Stderr)
		fs.PrintDefaults()
		os.Exit(1)
	}

	var embedding *glove.Embedding
	if err := serializer.LoadAny(embeddingFile, &embedding); err != nil {
		essentials.Die(err)
	}

	if !embedding.Tokens.Contains(token) {
		fmt.Fprintln(os.Stderr, "warning: token is not in vocabulary")
	}

	vec := embedding.Embed(token)
	matches, dists := embedding.Lookup(vec, numMatches)
	for i, matchID := range matches {
		fmt.Println(embedding.Tokens.Token(matchID), dists[i])
	}
}
