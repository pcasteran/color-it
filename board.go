package main

import (
	"strconv"
	"strings"
)

// Board represents the current status of the game.
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
	//   - a temporary frontier consisting of:
	//		- the cell immediately to the right (ID = 1)
	//		- the cell immediately to the bottom (ID = nbCols)
	board.completedCells[0] = void{}

	board.frontierCells[1] = void{}
	board.frontierCells[nbCols] = void{}

	// Compute the initial frontier.
	board.updateFrontier()

	return board
}

func (board *Board) clone() *Board {
	// Create a new board.
	clone := &Board{
		nbRows:         board.nbRows,
		nbCols:         board.nbCols,
		cells:          make(map[int]int, len(board.cells)),
		completedCells: make(map[int]void, len(board.completedCells)),
		frontierCells:  make(map[int]void, len(board.frontierCells)),
	}

	// Deep copy the nested data structures.
	for cellId, color := range board.cells {
		clone.cells[cellId] = color
	}

	for cellId := range board.completedCells {
		clone.completedCells[cellId] = void{}
	}

	for cellId := range board.frontierCells {
		clone.frontierCells[cellId] = void{}
	}

	return clone
}

// Return a string identifier uniquely identifying a board configuration.
func (board *Board) getId() string {
	// Create a string builder and add all the cells' color.
	var builder strings.Builder
	builder.Grow(board.nbRows * board.nbCols)
	for iRow := 0; iRow < board.nbRows; iRow++ {
		for iCol := 0; iCol < board.nbCols; iCol++ {
			cellId := (iRow * board.nbCols) + iCol
			color := board.cells[cellId]
			builder.WriteString(strconv.Itoa(color))
		}
	}
	return builder.String()
}

// Update the current frontier by looking at all the cells inside it and checking if their color is the
// same as the top-level cell. If yes:
//  1. the cell is removed from the frontier
//  2. the cell is integrated in the completed area
//  3. the adjacent cells are integrated into the frontier (if necessary)
func (board *Board) updateFrontier() {
	// Initialize the set of cells to process with the current frontier and clear the latter.
	cellsToProcess := board.frontierCells
	board.frontierCells = make(map[int]void, len(cellsToProcess))

	// Closure function processing one cell.
	processCell := func(cellId int) {
		// Check if the cell has not been already processed.
		_, alreadyCompleted := board.completedCells[cellId]
		if !alreadyCompleted {
			// The cell may already be present in the cellsToProcess map, but it's ok.
			// We save a look-up in the map by not testing the presence of the key prior to adding it.
			cellsToProcess[cellId] = void{}
		}
	}

	// Loop over the set of cells to process until it's empty.
	currentColor := board.cells[0]
	for {
		// Iterate over the cells to process.
		for cellId := range cellsToProcess {
			// Remove it from the cells to process.
			delete(cellsToProcess, cellId)

			// Check if the current cell has the same color as the top-left one.
			if board.cells[cellId] == currentColor {
				// Yes, add it to the completed area and mark the adjacent cells to be processed.
				board.completedCells[cellId] = void{}

				// Add the top, bottom, left and right adjacent cells if not already processed.
				row := cellId / board.nbCols
				col := cellId % board.nbCols

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

// Execute a step by:
//  1. changing the color of all the cells inside the completed area to the specified color
//  2. extending the current frontier
func (board *Board) playStep(color int) {
	// Update the color of all the cells in the completed area.
	for cellId := range board.completedCells {
		board.cells[cellId] = color
	}

	// Update the frontier.
	board.updateFrontier()
}

// Returns whether the board is solved, i.e. no more cell needs to be processed.
func (board *Board) isSolved() bool {
	return len(board.frontierCells) == 0
}
