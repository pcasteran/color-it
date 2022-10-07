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

// Implementation selecting the color that maximizes the converted area for N steps in the tree of configurations.
func maximizeStepAreaDeep(board *Board, solutions chan []int, done chan void, debug bool) ([]int, error) {
	return linearImpl(board, solutions, done, pickColorWithLargestAreaDeep, debug)
}

// Returns the color from the frontier with the largest area for N steps in the tree of configurations.
func pickColorWithLargestAreaDeep(board *Board) int {
	depth := 3 // 3 is the best trade-off between performance and accuracy
	color, _ := doPickColorWithLargestAreaDeep(board, depth)
	return color
}

func doPickColorWithLargestAreaDeep(board *Board, depth int) (int, *Board) {
	// Check if the board is solved.
	if board.isSolved() {
		return -1, board
	}

	// Check if the evaluation is finished.
	if depth <= 0 {
		// Evaluation finished.
		return -1, board
	}

	// Get the list of colors in the frontier ordered by descending area size.
	colors := board.getColorsInFrontier()

	// Try all the colors in the frontier.
	resultColor := -1
	var resultBoard *Board = nil

	for _, color := range colors {
		// Clone and update the board.
		boardCopy := board.clone()
		boardCopy.playStep(color)

		// Continue the evaluation.
		_, bestBoard := doPickColorWithLargestAreaDeep(boardCopy, depth-1)

		// Check if we improved the local best solution.
		if resultBoard == nil || len(bestBoard.completedCells) > len(resultBoard.completedCells) {
			resultColor = color
			resultBoard = bestBoard
		}
	}

	return resultColor, resultBoard
}
