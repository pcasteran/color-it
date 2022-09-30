package main

import (
	"fmt"
)

// Implementation selecting the color that maximizes the converted area for each step.
func maximizeStepArea(board *Board, debug bool) ([]int, error) {
	var solution []int

	// Loop until the board is solved.
	for {
		// Print the board status as CSV.
		if debug {
			fmt.Printf("Step #%d (color %d)\n", len(solution), board.cells[0])
			boardCsv, err := serializeBoardToCsv(board)
			if err != nil {
				return nil, fmt.Errorf("unable to serialize the board as CSV: %w", err)
			}
			fmt.Println(boardCsv)
		}

		// Check if the board is solved.
		if board.isSolved() {
			break
		}

		// TODO
		color := 0

		// Update the board.
		board.playStep(color)

		// Append the chosen color to the solution.
		solution = append(solution, color)
	}

	return solution, nil
}
