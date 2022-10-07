package main

import (
	"github.com/rs/zerolog/log"
	"math"
	"runtime"
	"sync"
)

// Implementation exploring the space of possibilities with a deep tree search to identify the optimal solution.
func deepSearch(board *Board, solutions chan []int, done chan void, debug bool) ([]int, error) {
	// First compute a "good" solution to have an initial step count that will be used to prune the graph search.
	// It's very probably not the optimal solution, but it's fast to compute.
	initialSolutionStepCount := computeInitialSolutionStepCount(board, solutions, debug)

	// Evaluate the board and return the best steps solution.
	ctx := &DeepSearchContext{
		debug:                 debug,
		bestSolutionStepCount: initialSolutionStepCount,
		processedCache:        make(map[string]*DeepSearchCacheEntry),
		solutions:             solutions,
	}
	solution := evaluateBoard(board, []int{}, ctx)

	// Print debug stats.
	ctx.printStats(true)

	// Notify that the execution is finished.
	if done != nil {
		done <- void{}
	}

	return solution, nil
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
func evaluateBoard(board *Board, steps []int, ctx *DeepSearchContext) []int {
	// Print debug stats.
	ctx.evaluationCounter++
	if ctx.evaluationCounter%10_000 == 0 {
		ctx.printStats(false)
	}

	// Get the current step count.
	currentStepCount := len(steps)

	// Check if the board is solved.
	if board.isSolved() {
		ctx.solvedCounter++

		// Check if we improved the overall best solution.
		if currentStepCount < ctx.bestSolutionStepCount {
			ctx.bestSolutionStepCount = currentStepCount

			// Push the new solution to the channel.
			ctx.solutions <- steps

			// TODO: clear cache with count >= currentStepCount ?
		}

		return steps
	}

	// Check if we can still hope to improve the current best solution.
	if currentStepCount >= (ctx.bestSolutionStepCount - 1) {
		// We can't improve, just stop there.
		ctx.prunedCounter++
		return nil
	}

	// Check if we have already processed this board configuration.
	boardId := board.getId()
	cacheEntry, alreadyProcessed := ctx.processedCache[boardId]
	if alreadyProcessed {
		ctx.cacheHitCounter++

		// Check if we are improving the best solution for this board.
		previousStepCount := cacheEntry.stepCount
		if currentStepCount < previousStepCount {
			ctx.cacheImprovedCounter++

			// We are improving, check if there was a valid previous solution.
			if cacheEntry.bestSolution != nil {
				ctx.cacheMergedCounter++

				// Merge the previous best solution with the current one.
				diff := previousStepCount - currentStepCount
				newStepCount := len(cacheEntry.bestSolution) - diff
				newBestSolution := make([]int, newStepCount)
				for i := 0; i < newStepCount; i++ {
					if i < currentStepCount {
						newBestSolution[i] = steps[i]
					} else {
						newBestSolution[i] = cacheEntry.bestSolution[i+diff]
					}
				}

				// Return the merged solution.
				return cacheEntry.bestSolution
			}

			// We are improving, but they were no valid previous solution for this board configuration.
			// Maybe there can be one now with a smaller step count, so we continue the evaluation (i.e. no return).
		} else {
			// We are not improving, stop now.
			return nil
		}
	}

	// Get the set of available colors in the frontier.
	colors := make(map[int]void)
	for cellId := range board.frontierCells {
		color := board.cells[cellId]
		colors[color] = void{}
	}

	// Try all the colors in the frontier and continue the evaluation.
	var localBestSolution []int = nil
	for color := range colors {
		// Clone and update the board.
		boardCopy := board.clone()
		boardCopy.playStep(color)

		// Copy the steps and append the current color.
		stepsCopy := make([]int, currentStepCount+1)
		copy(stepsCopy, steps)
		stepsCopy[len(stepsCopy)-1] = color

		// Continue the evaluation.
		solution := evaluateBoard(boardCopy, stepsCopy, ctx)

		// Update the cache.
		boardCopyId := board.getId()
		ctx.processedCache[boardCopyId] = &DeepSearchCacheEntry{
			stepCount:    currentStepCount + 1,
			bestSolution: solution,
		}

		// Check if we improved the local best solution.
		if solution != nil {
			// Check if the current solution is better than the best local one.
			solutionStepCount := len(solution)
			if localBestSolution == nil || solutionStepCount < len(localBestSolution) {
				// Yes, we improved the best local solution.
				localBestSolution = solution
			}
		}
	}

	// Update the cache.
	ctx.processedCache[boardId] = &DeepSearchCacheEntry{
		stepCount:    currentStepCount,
		bestSolution: localBestSolution,
	}

	return localBestSolution
}

type DeepSearchCacheEntry struct {
	stepCount    int
	bestSolution []int
}

type DeepSearchContext struct {
	debug                 bool
	bestSolutionStepCount int
	processedCache        map[string]*DeepSearchCacheEntry
	solutions             chan []int

	// Debug stats.
	evaluationCounter    int
	solvedCounter        int
	prunedCounter        int
	cacheHitCounter      int
	cacheImprovedCounter int
	cacheMergedCounter   int
}

func (ctx *DeepSearchContext) printStats(finished bool) {
	if ctx.debug {
		msg := "progress"
		if finished {
			msg = "finished"
		}

		log.Debug().
			Int("best", ctx.bestSolutionStepCount).
			Int("evaluation", ctx.evaluationCounter).
			Int("solved", ctx.solvedCounter).
			Int("pruned", ctx.prunedCounter).
			Int("cache-size", len(ctx.processedCache)).
			Int("cache-hit", ctx.cacheHitCounter).
			Int("cache-improvement", ctx.cacheImprovedCounter).
			Int("cache-merge", ctx.cacheMergedCounter).
			Msg(msg)
	}
}
