package main

import (
	"math/rand"
)

// Dummy implementation randomly selecting a color in the frontier at each step.
func dummy(board *Board, solutions chan []int, done chan void, debug bool) ([]int, error) {
	return linearImpl(board, solutions, done, randomPickColor, debug)
}

// Returns a randomly picked color from the frontier.
func randomPickColor(board *Board) int {
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
