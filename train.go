package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/unixpickle/anyvec/anyvec32"
	"github.com/unixpickle/essentials"
	"github.com/unixpickle/rip"
	"github.com/unixpickle/serializer"
	"github.com/unixpickle/wordembed/glove"
)

func CmdTrain(args []string) {
	var matrixFile string
	var vecSize int
	var outFile string
	var batchSize int
	var logInterval int
	fs := flag.NewFlagSet("train", flag.ExitOnError)
	fs.StringVar(&matrixFile, "matrix", "matrix_out", "co-occurrence matrix file")
	fs.IntVar(&vecSize, "vecsize", 128, "embedding vector size")
	fs.StringVar(&outFile, "out", "trainer_out", "output file")
	fs.IntVar(&batchSize, "batch", runtime.GOMAXPROCS(0), "SGD mini-batch size")
	fs.IntVar(&logInterval, "logint", 512, "iteration log interval")
	fs.Parse(args)

	var matrix *glove.SparseMatrix
	if err := serializer.LoadAny(matrixFile, &matrix); err != nil {
		essentials.Die(err)
	}

	var trainer *glove.Trainer
	if err := serializer.LoadAny(outFile, &trainer); err != nil {
		fmt.Fprintln(os.Stderr, "Creating a new trainer...")
		trainer = glove.NewTrainer(anyvec32.CurrentCreator(), vecSize, matrix)
	}

	creator := trainer.Vectors.Data.Creator()
	ops := creator.NumOps()

	r := rip.NewRIP()

	sinceLastLog := 0
	costSum := creator.MakeNumeric(0)
	for !r.Done() {
		cost := trainer.Update(batchSize)
		costSum = ops.Add(costSum, cost)
		sinceLastLog += batchSize
		if sinceLastLog > logInterval {
			log.Printf("done %d updates: cost=%f", trainer.NumUpdates,
				ops.Div(costSum, creator.MakeNumeric(float64(sinceLastLog))))
			sinceLastLog = 0
			costSum = creator.MakeNumeric(0)
		}
	}

	log.Println("Saving result...")
	if err := serializer.SaveAny(outFile, trainer); err != nil {
		essentials.Die(err)
	}
}