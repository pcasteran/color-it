package main

type void struct{}

type Board struct {
	dimension     int
	cells         map[int]int
	frontierCells map[int]void
}

func NewBoard(dimension int, cells map[int]int) *Board {
	// Create the board.
	board := &Board{
		dimension:     dimension,
		cells:         cells,
		frontierCells: make(map[int]void),
	}

	// Initialize the board with a frontier consisting of only the top-left cell (ID = 0).
	board.frontierCells[0] = void{}

	// Initialize the frontier cells.
	board.updateFrontierCells()

	return board
}

func (board *Board) updateFrontierCells() {
	// Iterate over all the frontier cells and check the adjacent cells.
	// The following rules will apply in order:
	//   1) If all the adjacent cells have the same color than the current cell, then the current cell is not part of
	//      the frontier anymore and is removed.
	//   2) The adjacent cell(s) that have the same color than the current cell are added to the frontier.
	// The preceding steps are repeated until no change is performed.
	currentColor := board.cells[0]
	for {
		changeOccured := false

		for cellId := range board.frontierCells {
			row := cellId / board.dimension
			col := cellId % board.dimension

			// Right
			rightAdded := false
			if col < (board.dimension - 1) {
				rightId := cellId + 1
				_, alreadyInFrontier := board.frontierCells[rightId]
				if !alreadyInFrontier && board.cells[rightId] == currentColor {
					board.frontierCells[rightId] = void{}
					rightAdded = true
				}
			}

			// Bottom
			bottomAdded := false
			if row < (board.dimension - 1) {
				bottomId := cellId + board.dimension
				_, alreadyInFrontier := board.frontierCells[bottomId]
				if !alreadyInFrontier && board.cells[bottomId] == currentColor {
					board.frontierCells[bottomId] = void{}
					bottomAdded = true
				}
			}

			// Check if some change occurred.
			if rightAdded || bottomAdded {
				changeOccured = true
			}

			// if both bottom and right are added to the frontier, the current cell can be removed.
			if rightAdded && bottomAdded {
				delete(board.frontierCells, cellId)
			}
		}

		// Check if some change occurred.
		if !changeOccured {
			// Done
			break
		}
	}
}
