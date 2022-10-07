package main

import (
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

// Available algorithm implementations.
var implementations = map[string]AlgorithmFn{
	"dummy":         dummy,
	"max-area":      maximizeStepArea,
	"max-area-deep": maximizeStepAreaDeep,
	"deep-search":   deepSearch,
}

func main() {
	// Parse the command line arguments.
	debug := flag.Bool("debug", false, "Enable the debug logs")
	impl := flag.String("impl", "deep-search", "Name of the algorithm implementation to execute")
	checkSquare := flag.Bool("check-square", true, "Check whether the board is a square after loading it")
	timeoutSec := flag.Int("timeout", 115, "Timeout in seconds of the execution")
	flag.Parse()

	inputFile := flag.Arg(0)

	// Configure logging. Default level is info, unless the debug flag is present.
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Load the board input file.
	board, err := readInputFile(inputFile, *checkSquare)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("input-file", inputFile).
			Bool("check-square", *checkSquare).
			Msg("unable to load the board input file")
	}

	// Get the algorithm implementation.
	implFn, exists := implementations[*impl]
	if !exists {
		impls := make([]string, 0, len(implementations))
		for i := range implementations {
			impls = append(impls, i)
		}
		log.Fatal().
			Strs("available", impls).
			Str("selected", *impl).
			Msg("invalid algorithm implementation specified")
	}

	// Execute it.
	var bestSolution []int = nil
	solutions := make(chan []int, 100)
	done := make(chan void)
	timeout := time.After(time.Duration(*timeoutSec) * time.Second)
	go func() {
		_, err := implFn(board, solutions, done, *debug)
		if err != nil {
			log.Fatal().Err(err).Msg("error during the algorithm execution")
		}
	}()

mainLoop:
	for {
		select {
		case solution := <-solutions:
			// A new solution has been pushed to the channel.
			if bestSolution == nil || len(solution) < len(bestSolution) {
				log.Info().Int("nb-steps", len(solution)).Ints("solution", solution).Msg("new best solution found")
				bestSolution = solution
			}
		case <-done:
			// The algorithm execution is finished.
			log.Info().Msg("algorithm execution finished")
			break mainLoop
		case <-timeout:
			// Timeout, the algorithm execution must be stopped.
			log.Warn().Msg("timeout reached during the algorithm execution")
			break mainLoop
		}
	}

	// Print the best solution found.
	log.Info().Int("nb-steps", len(bestSolution)).Ints("solution", bestSolution).Msg("best solution")
	for _, color := range bestSolution {
		fmt.Println(color)
	}

	// Generate the output file.
	err = writeOutputFile("out.csv", bestSolution)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("unable to write the solution to the output file")
	}
}
