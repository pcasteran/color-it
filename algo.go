package main

import "fmt"

// AlgorithmFn is the function type that will be used by all the implementations.
type AlgorithmFn func(board *Board, debug bool) ([]int, error)

// ColorPickerFn is the function type returning the color to play at the next step.
type ColorPickerFn func(board *Board) int

// Linear implementation using the provided colorPickerFn function to select the color to play at the next step.
func linearImpl(board *Board, colorPickerFn ColorPickerFn, debug bool) ([]int, error) {
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

		// Pick a color from the frontier.
		color := colorPickerFn(board)

		// Update the board.
		board.playStep(color)

		// Append the chosen color to the solution.
		solution = append(solution, color)
	}

	return solution, nil
}
