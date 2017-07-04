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

func CmdCorrelation(args []string) {
	var embeddingFile string
	fs := flag.NewFlagSet("correlation", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: correlation [flags] <a> <b>")
		fs.PrintDefaults()
	}
	fs.StringVar(&embeddingFile, "embedding", "embedding_out", "trained embedding")
	fs.Parse(args)
	if len(fs.Args()) != 2 {
		fs.Usage()
		os.Exit(1)
	}

	var embedding *glove.Embedding
	if err := serializer.LoadAny(embeddingFile, &embedding); err != nil {
		essentials.Die(err)
	}

	embedding.Normalize()

	embeddings := make([]anyvec.Vector, 2)
	for i, token := range fs.Args() {
		if !embedding.Tokens.Contains(token) {
			fmt.Fprintf(os.Stderr, "warning: token '%s' not in vocabulary\n", token)
		}
		embeddings[i] = embedding.Embed(token)
	}

	fmt.Println(embeddings[0].Dot(embeddings[1]))
}
