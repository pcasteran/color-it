package main

import (
	"github.com/rs/zerolog/log"
	"math"
	"runtime"
	"sync"
)

type DeepSearchContext struct {
	bestSolutionStepCount int
	processedCache        map[string]int
	solutions             chan []int
}

// Implementation exploring the space of possibilities with a deep tree search to identify the optimal solution.
func deepSearch(board *Board, solutions chan []int, done chan void, debug bool) ([]int, error) {
	// First compute a "good" solution to have an initial step count that will be used to prune the graph search.
	// It's very probably not the optimal solution, but it's fast to compute.
	initialSolutionStepCount := computeInitialSolutionStepCount(board, solutions, debug)

	// Evaluate the board.
	ctx := &DeepSearchContext{
		bestSolutionStepCount: initialSolutionStepCount,
		processedCache:        make(map[string]int),
		solutions:             solutions,
	}
	evaluateBoard(board, []int{}, ctx)

	// Notify that the execution is finished.
	done <- void{}

	return nil, nil
}

// Compute an initial solution using a fast but not optimal implementation and return the best step count found.
func computeInitialSolutionStepCount(board *Board, solutions chan []int, debug bool) int {
	// The fast implementation has some random parts (map iteration order). Thus, we launch multiple instances
	// in parallel and return the best one.
	var waitGroup sync.WaitGroup
	numCores := runtime.NumCPU()
	initialSolutionsStepCounts := make([]int, numCores)
	for i := 0; i < numCores; i++ {
		waitGroup.Add(1)
		id := i
		go func() {
			defer waitGroup.Done()

			// Call the fast implementation.
			solution, err := maximizeStepArea(board.clone(), solutions, nil, debug)
			if err != nil {
				// No solution found, this is unfortunate but not blocking.
				log.Warn().Err(err).Int("id", id).Msg("unable to compute the initial solution")
				initialSolutionsStepCounts[id] = math.MaxInt
			} else {
				// The initial solution is valid.
				initialSolutionsStepCounts[id] = len(solution)
			}
		}()
	}
	waitGroup.Wait()

	// Compute the best initial solution step count.
	initialSolutionStepCount := math.MaxInt
	for _, count := range initialSolutionsStepCounts {
		if count < initialSolutionStepCount {
			initialSolutionStepCount = count
		}
	}
	log.Info().Int("step-count", initialSolutionStepCount).Msg("initial solution found")

	return initialSolutionStepCount
}

// Recursive function to evaluate a board and the possible solution(s) from it.
func evaluateBoard(board *Board, steps []int, ctx *DeepSearchContext) {
	// Get the current step count.
	currentStepCount := len(steps)

	// Check if we have already processed this board.
	boardId := board.getId()
	previousBestStepCount, alreadyProcessed := ctx.processedCache[boardId]
	if alreadyProcessed {
		// Check if we are improving the best solution for this board.
		if currentStepCount >= previousBestStepCount {
			// We are not improving, stop now.
			return
		}
	} else {
		// Update the cache.
		ctx.processedCache[boardId] = currentStepCount
	}

	// Check if the board is solved.
	if board.isSolved() {
		// Check if we improved the overall best solution.
		if currentStepCount < ctx.bestSolutionStepCount {
			ctx.bestSolutionStepCount = currentStepCount

			// Push the new solution to the channel.
			ctx.solutions <- steps
		}

		return
	}

	// Check if we can still hope to improve the current best solution.
	if !(currentStepCount < (ctx.bestSolutionStepCount - 1)) {
		// We can't improve, just stop there.
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
		// Clone and update the board.
		boardCopy := board.clone()
		boardCopy.playStep(color)

		// Copy the steps and append the current color.
		stepsCopy := make([]int, currentStepCount+1)
		copy(stepsCopy, steps)
		stepsCopy[len(stepsCopy)-1] = color

		// Continue the evaluation.
		evaluateBoard(boardCopy, stepsCopy, ctx)
	}
}
