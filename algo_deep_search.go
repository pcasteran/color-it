package main

import (
	"github.com/rs/zerolog/log"
	"math"
)

// Implementation exploring the space of possibilities with a deep tree search to identify the optimal solution.
func deepSearch(board *Board, debug bool) ([]int, error) {
	// First compute a "good" solution to have an initial step count that will be used to prune the graph search.
	// It's very probably not the optimal solution, but it's fast to compute.
	initialSolutionStepCount := math.MaxInt
	initialSolution, err := maximizeStepArea(board.clone(), debug)
	if err != nil {
		// No initial solution, this is unfortunate but not blocking.
		log.Warn().Err(err).Msg("unable to compute the initial solution")
	} else {
		// The initial solution is valid.
		initialSolutionStepCount = len(initialSolution)
		log.Info().Int("step-count", initialSolutionStepCount).Msg("Initial solution found")
	}

	// Evaluate the board and return the best steps solution.
	solution := evaluateBoard(board, []int{}, initialSolutionStepCount)

	return solution, nil
}

// Recursive function to evaluate a board and the possible solution(s) from it.
func evaluateBoard(board *Board, steps []int, bestStepCount int) []int {
	// Check if the board is solved.
	if board.isSolved() {
		return steps
	}

	// Check if we can still hope to improve the current best solution.
	if !(len(steps) < (bestStepCount - 1)) {
		// We can't improve, just stop there for this branch.
		return nil
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

		stepsCopy := make([]int, len(steps)+1)
		copy(stepsCopy, steps)
		stepsCopy[len(stepsCopy)-1] = color

		// Continue the evaluation.
		solution := evaluateBoard(boardCopy, stepsCopy, bestStepCount)
		if solution != nil {
			// Check if the current solution is better than the best local one.
			solutionStepCount := len(solution)
			if bestSolution == nil || solutionStepCount < len(bestSolution) {
				// Yes, we improved the best local solution.
				bestSolution = solution

				// Check if we improved the overall best solution.
				if solutionStepCount < bestStepCount {
					bestStepCount = solutionStepCount
				}
			}
		}
	}

	return bestSolution
}
