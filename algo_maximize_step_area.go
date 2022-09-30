package main

import (
	"math/rand"
)

// Implementation selecting the color that maximizes the converted area for each step.
func maximizeStepArea(board *Board, debug bool) ([]int, error) {
	return linearImpl(board, pickColorWithLargestArea, debug)
}

// Returns the color from the frontier with the largest area.
func pickColorWithLargestArea(board *Board) int {
	// Get the set of available colors in the frontier.
	colors := make(map[int]void)
	for cellId := range board.frontierCells {
		color := board.cells[cellId]
		colors[color] = void{}
	}

	// Get the frontier colors as a slice.
	choices := make([]int, len(colors))
	i := 0
	for color := range colors {
		choices[i] = color
		i++
	}

	// Pick one color randomly.
	choiceIdx := rand.Intn(len(choices))
	color := choices[choiceIdx]

	return color
}
