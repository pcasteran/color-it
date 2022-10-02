package main

// Implementation selecting the color that maximizes the converted area for each step.
func maximizeStepArea(board *Board, solutions chan []int, done chan void, debug bool) ([]int, error) {
	return linearImpl(board, solutions, done, pickColorWithLargestArea, debug)
}

// Returns the color from the frontier with the largest area.
func pickColorWithLargestArea(board *Board) int {
	// Compute the set of cells (areas) accessible from the frontier and grouped by color.
	areasByColor := make(map[int]map[int]void)

	// Initialize the set of cells to process with the current frontier.
	cellsToProcess := make(map[int]void, len(board.frontierCells))
	for cellId := range board.frontierCells {
		cellsToProcess[cellId] = void{}
	}

	// Initialized the set of processed cells to an empty set.
	processedCells := make(map[int]void)

	// Closure function processing one cell.
	processCell := func(cellId, expectedColor int) {
		// Check if the cell color is the same as the expected one.
		color := board.cells[cellId]
		if color == expectedColor {
			// Check if the cell has not been already processed.
			_, alreadyProcessed := processedCells[cellId]
			if !alreadyProcessed {
				// The cell may already be present in the cellsToProcess map, but it's ok.
				// We save a look-up in the map by not testing the presence of the key prior to adding it.
				cellsToProcess[cellId] = void{}
			}
		}
	}

	// Loop over the set of cells to process until it's empty.
	for {
		// Iterate over the cells to process.
		for cellId := range cellsToProcess {
			// Remove it from the cells to process.
			delete(cellsToProcess, cellId)
			processedCells[cellId] = void{}

			// Add it to the area corresponding to its color.
			color := board.cells[cellId]
			area := areasByColor[color]
			if area == nil {
				// Lazy initialization of the area set for this color.
				area = make(map[int]void)
				areasByColor[color] = area
			}
			area[cellId] = void{}

			// Check if the top, bottom, left and right adjacent cells are of the same color.
			row := cellId / board.nbCols
			col := cellId % board.nbCols

			// Top
			if row > 0 {
				processCell(cellId-board.nbCols, color)
			}

			// Bottom
			if row < (board.nbRows - 1) {
				processCell(cellId+board.nbCols, color)
			}

			// Left
			if col > 0 {
				processCell(cellId-1, color)
			}

			// Right
			if col < (board.nbCols - 1) {
				processCell(cellId+1, color)
			}
		}

		if len(cellsToProcess) == 0 {
			// No more cell to process, we are done.
			break
		}
	}

	// Iterate over the frontier colors and keep the one with the largest area (i.e. the greater cell count).
	largestAreaSize, largestAreaColor := -1, -1
	for color, area := range areasByColor {
		areaSize := len(area)
		if largestAreaSize == -1 || areaSize > largestAreaSize {
			largestAreaSize = areaSize
			largestAreaColor = color
		}
	}

	return largestAreaColor
}
