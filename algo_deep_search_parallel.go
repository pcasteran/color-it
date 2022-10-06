package main

import (
	"github.com/rs/zerolog/log"
	"runtime"
	"sync/atomic"
)

// Implementation exploring the space of possibilities with a deep tree search to identify the optimal solution.
// Same logic than deepSearch but using multiple goroutines to take advantage of the available CPU cores.
func deepSearchParallel(board *Board, solutions chan []int, done chan void, debug bool) ([]int, error) {
	// First compute a "good" solution to have an initial step count that will be used to prune the graph search.
	// It's very probably not the optimal solution, but it's fast to compute.
	initialSolutionStepCount := computeInitialSolutionStepCount(board, solutions, debug)

	// Evaluate the board.
	ctx := NewDeepSearchParallelContext(int64(initialSolutionStepCount), solutions)
	evaluateBoardParallel(board, []int{}, ctx)

	// Notify that the execution is finished.
	done <- void{}

	return nil, nil
}

// Recursive function to evaluate a board and the possible solution(s) from it.
func evaluateBoardParallel(board *Board, steps []int, ctx *DeepSearchParallelContext) {
	// Check if the board is solved.
	if board.isSolved() {
		ctx.processNewSolution(steps)
		return
	}

	// Check if we can still hope to improve the current best solution.
	bestSolutionStepCount := ctx.bestSolutionStepCount.Load()
	currentStepCount := int64(len(steps))
	if !(currentStepCount < (bestSolutionStepCount - 1)) {
		// We can't improve, just stop there for this branch.
		return
	}

	// Get the set of available colors in the frontier.
	colors := make(map[int]void)
	for cellId := range board.frontierCells {
		color := board.cells[cellId]
		colors[color] = void{}
	}

	// Try all the colors in the frontier and continue the evaluation.
	for color := range colors {
		playAndEvaluateParallel(board, steps, color, ctx)
	}
}

func playAndEvaluateParallel(board *Board, steps []int, color int, ctx *DeepSearchParallelContext) {
	// Clone and update the board.
	boardCopy := board.clone()
	boardCopy.playStep(color)

	// Copy the steps and append the current color.
	stepsCopy := make([]int, len(steps)+1)
	copy(stepsCopy, steps)
	stepsCopy[len(stepsCopy)-1] = color

	// Continue the evaluation.
	evaluateBoardParallel(boardCopy, stepsCopy, ctx)
}

type DeepSearchParallelContext struct {
	availableGoroutinesCounter atomic.Int64

	bestSolutionStepCount atomic.Int64

	solutions chan []int
}

func NewDeepSearchParallelContext(bestSolutionStepCount int64, solutions chan []int) *DeepSearchParallelContext {
	// Create and initialize a new context.
	ctx := &DeepSearchParallelContext{
		solutions: solutions,
	}
	ctx.availableGoroutinesCounter.Store(int64(runtime.NumCPU()))
	ctx.bestSolutionStepCount.Store(bestSolutionStepCount)

	return ctx
}

func (ctx *DeepSearchParallelContext) processNewSolution(solution []int) {
	solutionStepCount := int64(len(solution))
	for {
		// Get the current best solution step count.
		current := ctx.bestSolutionStepCount.Load()

		// Check if we improved the best solution.
		if solutionStepCount < current {
			// Try to atomically update the best solution step count.
			swapped := ctx.bestSolutionStepCount.CompareAndSwap(current, solutionStepCount)
			if swapped {
				// Push the new solution to the channel.
				ctx.solutions <- solution

				// Atomic update done, exit the for loop.
				break
			} else {
				// Atomic update failed, the for loop will retry it.
				log.Warn().Msg("best solution step count atomic update failed")
			}
		} else {
			// No need to update the best solution step count.
			break
		}
	}
}
