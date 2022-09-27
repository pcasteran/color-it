package main

import (
	"log"
)

// TODO: function documentation

func main() {
	// Load the board input file.
	_, err := readInputFile("samples/2.csv")
	if err != nil {
		// TODO: log
		log.Fatal("Error while loading the board input file: ", err)
	}
}
