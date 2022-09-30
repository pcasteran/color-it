package main

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

// Available algorithm implementations.
var implementations = map[string]AlgorithmFn{
	"dummy":    dummy,
	"max-area": maximizeStepArea,
}

func main() {
	// Parse the command line arguments.
	debug := flag.Bool("debug", false, "Enable the debug logs")
	impl := flag.String("impl", "max-area", "Name of the algorithm implementation to execute")
	checkSquare := flag.Bool("check-square", true, "Check whether the board is a square after loading it")
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
	solution, err := implFn(board, *debug)
	if err != nil {
		log.Fatal().Err(err).Msg("error during the algorithm execution")
	}
	log.Info().Int("nb-steps", len(solution)).Ints("solution", solution).Msgf("solution found !!")
}
