package main

import (
	"fmt"
	"os"

	"github.com/unixpickle/essentials"
)

var SubCommands = map[string]func([]string){
	"tokens": CmdTokens,
	"matrix": CmdMatrix,
	"train":  CmdTrain,
}

var Descriptions = map[string]string{
	"tokens": "build a vocabulary of tokens",
	"matrix": "build a co-occurrence matrix",
	"train":  "train a GloVe model",
}

func main() {
	if len(os.Args) < 2 {
		dieHelp()
	}

	if handler, ok := SubCommands[os.Args[1]]; ok {
		handler(os.Args[2:])
	} else {
		fmt.Fprintln(os.Stderr, "Unknown subcommand:", os.Args[1])
		fmt.Fprintln(os.Stderr)
		dieHelp()
	}
}

func dieHelp() {
	fmt.Fprintln(os.Stderr, "Usage: tweetembed <sub-command> [args]")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Available sub-commands:")

	var names []string
	var maxLen int
	for name := range SubCommands {
		names = append(names, name)
		maxLen = essentials.MaxInt(len(name), maxLen)
	}

	for _, name := range names {
		desc := Descriptions[name]
		for len(name) < maxLen {
			name += " "
		}
		fmt.Fprintln(os.Stderr, " "+name+"  "+desc)
	}

	fmt.Fprintln(os.Stderr)

	os.Exit(1)
}
