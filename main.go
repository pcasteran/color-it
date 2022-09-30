package main

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	// Parse the command line arguments.
	debug := flag.Bool("debug", false, "Set the log level to debug")
	algo := flag.String("algo", "dummy", "Name of the algorithm to execute")
	checkSquare := flag.Bool("check-square", true, "Check whether the board is a square")
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
	impl, exists := implementations[*algo]
	if !exists {
		log.Fatal().Err(err).Str("algorithm", *algo).Msg("invalid algorithm specified")
	}

	// Execute it.
	solution, err := impl(board)
	if err != nil {
		log.Fatal().Err(err).Msg("error during the algorithm execution")
	}
	log.Info().Int("nb-steps", len(solution)).Ints("solution", solution).Msgf("solution found !!")
}
