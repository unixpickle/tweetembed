package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/unixpickle/anyvec/anyvec32"
	"github.com/unixpickle/essentials"
	"github.com/unixpickle/rip"
	"github.com/unixpickle/serializer"
	"github.com/unixpickle/wordembed/glove"
)

func CmdTrain(args []string) {
	rand.Seed(time.Now().UnixNano())

	var matrixFile string
	var vecSize int
	var outFile string
	var batchSize int
	var logInterval int
	fs := flag.NewFlagSet("train", flag.ExitOnError)
	fs.StringVar(&matrixFile, "matrix", "matrix_out", "co-occurrence matrix file")
	fs.IntVar(&vecSize, "vecsize", 128, "embedding vector size")
	fs.StringVar(&outFile, "out", "trainer_out", "output file")
	fs.IntVar(&batchSize, "batch", runtime.GOMAXPROCS(0)*32, "SGD mini-batch size")
	fs.IntVar(&logInterval, "logint", 16, "iteration log interval")
	fs.Parse(args)

	log.Println("Loading matrix...")
	var matrix *glove.SparseMatrix
	if err := serializer.LoadAny(matrixFile, &matrix); err != nil {
		essentials.Die(err)
	}

	log.Println("Loading trainer...")
	var trainer *glove.Trainer
	if err := serializer.LoadAny(outFile, &trainer); err != nil {
		fmt.Fprintln(os.Stderr, "Load error:", err)
		log.Println("Creating a new trainer...")
		trainer = glove.NewTrainer(anyvec32.CurrentCreator(), vecSize, matrix)
	} else {
		trainer.Cooccur = matrix
	}

	creator := trainer.Vectors.Data.Creator()
	ops := creator.NumOps()

	r := rip.NewRIP()

	sinceLastLog := 0
	costSum := creator.MakeNumeric(0)
	for !r.Done() {
		cost := trainer.Update(batchSize)
		costSum = ops.Add(costSum, cost)
		sinceLastLog++
		if sinceLastLog > logInterval {
			log.Printf("done %d updates: cost=%f", trainer.NumUpdates,
				ops.Div(costSum, creator.MakeNumeric(float64(sinceLastLog))))
			sinceLastLog = 0
			costSum = creator.MakeNumeric(0)
		}
	}

	log.Println("Saving result...")

	// Don't save extra copy of the huge matrix.
	n := trainer.Vectors.Rows
	trainer.Cooccur = glove.NewSparseMatrix(n, n)

	// Try to get GC to let go of the old matrix.
	matrix = nil
	runtime.GC()

	if err := serializer.SaveAny(outFile, trainer); err != nil {
		essentials.Die(err)
	}
}
