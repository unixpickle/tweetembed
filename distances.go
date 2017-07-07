package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/unixpickle/anyvec"
	"github.com/unixpickle/essentials"
	"github.com/unixpickle/serializer"
	"github.com/unixpickle/wordembed/glove"
)

func CmdDistances(args []string) {
	var embeddingFile string
	fs := flag.NewFlagSet("distances", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: correlation [flags] [words...]")
		fs.PrintDefaults()
	}
	fs.StringVar(&embeddingFile, "embedding", "embedding_out", "trained embedding")
	fs.Parse(args)
	if len(fs.Args()) == 0 {
		fs.Usage()
		os.Exit(1)
	}

	var embedding *glove.Embedding
	if err := serializer.LoadAny(embeddingFile, &embedding); err != nil {
		essentials.Die(err)
	}

	embeddings := make([]anyvec.Vector, len(fs.Args()))
	for i, token := range fs.Args() {
		if !embedding.Tokens.Contains(token) {
			fmt.Fprintf(os.Stderr, "warning: token '%s' not in vocabulary\n", token)
		}
		embeddings[i] = embedding.Embed(token)
	}

	for _, vec := range embeddings {
		var distStrs []string
		for _, otherVec := range embeddings {
			diff := vec.Copy()
			diff.Sub(otherVec)
			distStrs = append(distStrs, fmt.Sprintf("%f", anyvec.Norm(diff)))
		}
		fmt.Println(strings.Join(distStrs, ","))
	}
}
