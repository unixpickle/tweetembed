package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/unixpickle/essentials"
	"github.com/unixpickle/serializer"
	"github.com/unixpickle/wordembed/glove"
)

func CmdEntries(args []string) {
	rand.Seed(time.Now().UnixNano())

	var matrixFile string
	fs := flag.NewFlagSet("entries", flag.ExitOnError)
	fs.StringVar(&matrixFile, "matrix", "matrix_out", "co-occurrence matrix file")
	fs.Parse(args)

	log.Println("Loading matrix...")
	var matrix *glove.SparseMatrix
	if err := serializer.LoadAny(matrixFile, &matrix); err != nil {
		essentials.Die(err)
	}

	log.Println("Counting entries...")
	fmt.Println("Count:", matrix.NumEntries())
}
