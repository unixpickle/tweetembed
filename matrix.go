package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/unixpickle/essentials"
	"github.com/unixpickle/serializer"
	"github.com/unixpickle/wordembed"
	"github.com/unixpickle/wordembed/glove"
)

func CmdMatrix(args []string) {
	var tweetFile string
	var tokensFile string
	var outFile string
	fs := flag.NewFlagSet("matrix", flag.ExitOnError)
	fs.StringVar(&tweetFile, "tweets", "", "CSV file of tweets (empty means stdin)")
	fs.StringVar(&tokensFile, "tokens", "tokens_out", "vocabulary file of tokens")
	fs.StringVar(&outFile, "out", "matrix_out", "output file")
	fs.Parse(args)

	tweets, err := ReadTweets(tweetFile)
	if err != nil {
		essentials.Die(err)
	}

	var tokens wordembed.TokenSet
	if err := serializer.LoadAny(tokensFile, &tokens); err != nil {
		essentials.Die(err)
	}

	tokenizer := &wordembed.Tokenizer{}
	counter := &glove.CooccurCounter{
		Tokens:      tokens,
		Matrix:      glove.NewSparseMatrix(tokens.NumIDs(), tokens.NumIDs()),
		WeightWords: true,
	}

	var processed int
	for tweet := range tweets {
		counter.Add(tokenizer.Tokenize(tweet))
		processed++
		if processed%512 == 0 {
			fmt.Fprintf(os.Stderr, "\rprocessed %d tweets", processed)
		}
	}
	fmt.Fprintln(os.Stderr, "")

	if err := serializer.SaveAny(outFile, counter.Matrix); err != nil {
		essentials.Die(err)
	}
}
