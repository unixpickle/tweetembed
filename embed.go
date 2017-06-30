package main

import (
	"flag"
	"log"

	"github.com/unixpickle/essentials"
	"github.com/unixpickle/serializer"
	"github.com/unixpickle/wordembed"
	"github.com/unixpickle/wordembed/glove"
)

func CmdEmbed(args []string) {
	var trainerFile string
	var tokensFile string
	var outFile string
	fs := flag.NewFlagSet("embed", flag.ExitOnError)
	fs.StringVar(&tokensFile, "tokens", "tokens_out", "vocabulary file of tokens")
	fs.StringVar(&trainerFile, "trainer", "trainer_out", "trainer save file")
	fs.StringVar(&outFile, "out", "embedding_out", "output file")
	fs.Parse(args)

	var tokens wordembed.TokenSet
	if err := serializer.LoadAny(tokensFile, &tokens); err != nil {
		essentials.Die(err)
	}

	var trainer *glove.Trainer
	if err := serializer.LoadAny(trainerFile, &trainer); err != nil {
		essentials.Die(err)
	}

	log.Println("Creating embedding...")
	e := trainer.Embedding(tokens, true)

	log.Println("Saving result...")
	if err := serializer.SaveAny(outFile, e); err != nil {
		essentials.Die(err)
	}
}
