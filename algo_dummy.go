package main

import (
	"fmt"
	"math/rand"
)

// Dummy implementation randomly selecting a color in the frontier at each step.
func dummy(board *Board) ([]int, error) {
	var solution []int

	// Loop until the board is solved.
	for {
		// Print the board status as CSV.
		fmt.Printf("Step #%d (color %d)\n", len(solution), board.cells[0])
		boardCsv, err := serializeBoardToCsv(board)
		if err != nil {
			return nil, fmt.Errorf("unable to serialize the board as CSV: %w", err)
		}
		fmt.Println(boardCsv)

		// Check if the board is solved.
		if board.isSolved() {
			break
		}

		// Randomly pick a color from the frontier.
		color := randomPickColor(board)

		// Update the board.
		board.updateCompletedArea(color)

		// Append the chosen color to the solution.
		solution = append(solution, color)
	}

	return solution, nil
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
