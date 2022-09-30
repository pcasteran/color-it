package main

// Implementation exploring the space of possibilities with a deep tree search to identify the optimal solution.
func deepSearch(board *Board, debug bool) ([]int, error) {
	// Evaluate the board and return the best steps solution.
	solution := evaluateBoard(board, []int{})
	return solution, nil
}

func evaluateBoard(board *Board, steps []int) []int {
	// Check if the board is solved.
	if board.isSolved() {
		return steps
	}

	// Get the set of available colors in the frontier.
	colors := make(map[int]void)
	for cellId := range board.frontierCells {
		color := board.cells[cellId]
		colors[color] = void{}
	}

	// Try all the colors in the frontier and continue the evaluation.
	var bestSolution []int = nil
	for color := range colors {
		// Clone and update the board.
		boardCopy := board.clone()
		boardCopy.playStep(color)

		// Continue the evaluation.
		stepsCopy := make([]int, len(steps))
		copy(stepsCopy, steps)
		stepsCopy = append(stepsCopy, color)

		solution := evaluateBoard(boardCopy, stepsCopy)
		if bestSolution == nil || len(solution) < len(bestSolution) {
			bestSolution = solution
		}
	}

	return bestSolution
}
