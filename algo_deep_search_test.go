package main

import "testing"

import (
	"github.com/rs/zerolog/log"
)

func BenchmarkDeepSearch(b *testing.B) {
	// Prepare the implementation parameters.
	board, err := readInputFile("samples/30_30_3-1.csv", false)
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
		_, err := deepSearch(board.clone(), solutions, nil, false)
		if err != nil {
			log.Fatal().Err(err).Msg("error during the algorithm execution")
		}
	}
}
