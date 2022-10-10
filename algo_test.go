package main

import (
	"github.com/rs/zerolog/log"
	"testing"
)

func benchmarkImplementation(b *testing.B, implFn AlgorithmFn, inputFile string) {
	// Prepare the implementation parameters.
	board, err := readInputFile(inputFile, false)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("unable to load the board input file")
	}
	solutions := make(chan []int, 100)

	// Launch a go routine that consumes all the solutions.
	go func() {
		for {
			_ = <-solutions
		}
	}()

	// Run the implementation to benchmark b.N times.
	for n := 0; n < b.N; n++ {
		_, err := implFn(board.clone(), solutions, nil, false)
		if err != nil {
			log.Fatal().Err(err).Msg("error during the algorithm execution")
		}
	}
}
