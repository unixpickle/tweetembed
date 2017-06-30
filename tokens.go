package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/unixpickle/essentials"
	"github.com/unixpickle/serializer"
	"github.com/unixpickle/wordembed"
)

func CmdTokens(args []string) {
	var tweetFile string
	var numTokens int
	var outFile string
	fs := flag.NewFlagSet("tokens", flag.ExitOnError)
	fs.StringVar(&tweetFile, "tweets", "", "CSV file of tweets (empty means stdin)")
	fs.IntVar(&numTokens, "num", 100000, "number of tokens to save")
	fs.StringVar(&outFile, "out", "tokens_out", "output file")
	fs.Parse(args)

	tweets, err := ReadTweets(tweetFile)
	if err != nil {
		essentials.Die(err)
	}

	tokenizer := &wordembed.Tokenizer{}
	counts := wordembed.TokenCounts{}

	var processed int
	for tweet := range tweets {
		for _, token := range tokenizer.Tokenize(tweet) {
			counts.Add(token)
		}
		processed++
		if processed%512 == 0 {
			fmt.Fprintf(os.Stderr, "\rprocessed %d tweets (%d tokens)", processed,
				len(counts))
		}
	}
	fmt.Fprintln(os.Stderr)

	toks := counts.MostCommon(numTokens)
	if err := serializer.SaveAny(outFile, toks); err != nil {
		essentials.Die(err)
	}
}
