package main

import (
	"context"

	"os"
	"os/signal"
	"syscall"
	"vocab/harvester"
	"vocab/lib"
	"vocab/vocabulary/forest"
)

func main() {
	print("Starting vocab-ls...\n")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	inputReader := lib.NewInputReader(os.Stdin)
	outputWriter := lib.NewOutputWriter(os.Stdout)
	logger := lib.NewLogger(os.Stderr)
	forest := forest.NewForest(ctx, func(any) {})
	h := harvester.NewHarvester(
		ctx,
		forest,
		inputReader.Read,
		outputWriter.Write,
		logger,
	)
	h.Start()
}
