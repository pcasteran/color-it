package main

// Implementation selecting the color that maximizes the converted area for each step.
func maximizeStepArea(board *Board, solutions chan []int, done chan void, debug bool) ([]int, error) {
	return linearImpl(board, solutions, done, pickColorWithLargestArea, debug)
}

// Returns the color from the frontier with the largest area.
func pickColorWithLargestArea(board *Board) int {
	// Get the list of colors in the frontier ordered by descending area size.
	// Return the first one (guaranteed to exist as the board is not solved).
	return board.getColorsInFrontier()[0]
}
