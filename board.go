package main

// Empty struct (no memory usage) to use as the value for the cell maps as Go doesn't have a set data structure.
type void struct{}

type Board struct {
	// Number of rows in the board.
	nbRows int

	// Number of columns in the board.
	nbCols int

	// Map of the board cells with the cell ID as key and the color as value.
	cells map[int]int

	// Set of the cells that are in the same contiguous area with the same color as the top-left cell.
	// They do not need to be processed anymore.
	completedCells map[int]void

	// Set of the cells adjacent to the completedCells area.
	// They may or may not be of the same color as the top-left cell.
	frontierCells map[int]void
}

func NewBoard(nbRows, nbCols int, cells map[int]int) *Board {
	// Create the board.
	board := &Board{
		nbRows:         nbRows,
		nbCols:         nbCols,
		cells:          cells,
		completedCells: make(map[int]void),
		frontierCells:  make(map[int]void),
	}

	// Initialize the board with:
	//   - a completed area consisting of only the top-left cell (ID = 0)
	//   - a frontier consisting of:
	//		- the cell immediately to the right (ID = 1)
	//		- the cell immediately to the bottom (ID = nbCols)
	board.completedCells[0] = void{}

	board.frontierCells[1] = void{}
	board.frontierCells[nbCols] = void{}

	// Update the frontier.
	board.updateFrontier()

	return board
}

func (board *Board) updateFrontier() {
	// Initialize the set of cells to process with the current frontier and clear the latter.
	cellsToProcess := board.frontierCells
	board.frontierCells = make(map[int]void, len(cellsToProcess))

	// Loop over the set of cells to process until it's empty.
	currentColor := board.cells[0]
	for {
		// Iterate over the cells to process.
		for cellId := range cellsToProcess {
			// Remove it.
			delete(cellsToProcess, cellId)

			// Check if the current cell has the same color as the top-left one.
			if board.cells[cellId] == currentColor {
				// Yes, add it to the completed area and mark the adjacent cells to be processed.
				board.completedCells[cellId] = void{}

				// Add the top, bottom, left and right adjacent cells if not already processed.
				row := cellId / board.nbCols
				col := cellId % board.nbCols

				processCell := func(cellId int) {
					_, alreadyCompleted := board.completedCells[cellId]
					if !alreadyCompleted {
						// The cell may already be present in the cellsToProcess map, but it's ok.
						// We save a map look-up by not checking it prior to adding it.
						cellsToProcess[cellId] = void{}
					}
				}

				// Top
				if row > 0 {
					processCell(cellId - board.nbCols)
				}

				// Bottom
				if row < (board.nbRows - 1) {
					processCell(cellId + board.nbCols)
				}

				// Left
				if col > 0 {
					processCell(cellId - 1)
				}

				// Right
				if col < (board.nbCols - 1) {
					processCell(cellId + 1)
				}
			} else {
				// No, add it to the frontier.
				board.frontierCells[cellId] = void{}
			}
		}

		if len(cellsToProcess) == 0 {
			// No more cell to process, we are done.
			break
		}
	}
}
