package main

import (
	"fmt"
	"log"
)

// TODO: function documentation

func main() {
	// Load the board input file.
	board, err := readInputFile("samples/1.csv")
	if err != nil {
		// TODO: log
		log.Fatal("Error while loading the board input file: ", err)
	}

	// TODO
	impl := "dummy"
	algo, exists := implementations[impl]
	if !exists {
		// TODO: log
		log.Fatal(fmt.Sprintf("Invalid algorithm specified [%s]", impl))
	}
	solution, err := algo(board)
	if err != nil {
		// TODO: log
		log.Fatal("Error during the algorithm execution: ", err)
	}
	log.Printf("Solution found in %d steps: %v\n", len(solution), solution)
}
