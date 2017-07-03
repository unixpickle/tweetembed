package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/unixpickle/anyvec"
	"github.com/unixpickle/essentials"
	"github.com/unixpickle/serializer"
	"github.com/unixpickle/wordembed/glove"
)

func CmdAnalogy(args []string) {
	var embeddingFile string
	var numMatches int
	fs := flag.NewFlagSet("analogy", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: analogy [flags] <a> <b> <c>")
		fs.PrintDefaults()
	}
	fs.StringVar(&embeddingFile, "embedding", "embedding_out", "trained embedding")
	fs.IntVar(&numMatches, "num", 5, "number of matches to show")
	fs.Parse(args)
	if len(fs.Args()) != 3 {
		fs.Usage()
		os.Exit(1)
	}

	var embedding *glove.Embedding
	if err := serializer.LoadAny(embeddingFile, &embedding); err != nil {
		essentials.Die(err)
	}

	embeddings := make([]anyvec.Vector, 3)
	for i, token := range fs.Args() {
		if !embedding.Tokens.Contains(token) {
			fmt.Fprintf(os.Stderr, "warning: token '%s' not in vocabulary\n", token)
		}
		embeddings[i] = embedding.Embed(token)
	}

	embeddings[1].Sub(embeddings[0])
	embeddings[2].Add(embeddings[1])
	matches, dists := embedding.Lookup(embeddings[2], numMatches)
	for i, matchID := range matches {
		fmt.Println(embedding.Tokens.Token(matchID), dists[i])
	}
}
