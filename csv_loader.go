package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

// Read the input CSV file specified by its path and load a board from its content.
func readInputFile(filePath string, checkSquare bool) (*Board, error) {
	// Load the raw string file content.
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open the input file: %w", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("unable to close the input file")
		}
	}(f)

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to parse the input CSV file: %w", err)
	}

	// Parse it.
	cells := make(map[int]int)
	for iRow, columns := range records {
		for iCol, colorStr := range columns {
			cellId := (iRow * len(columns)) + iCol
			color, err := strconv.Atoi(colorStr)
			if err != nil {
				return nil, fmt.Errorf("invalid color for row=%d, col=%d, color=%s : %w", iRow+1, iCol+1, colorStr, err)
			}
			cells[cellId] = color
		}
	}

	// Check that the board is a square.
	nbRows := len(records)
	nbCols := len(cells) / nbRows
	if checkSquare && len(cells) != (nbRows*nbRows) {
		return nil, fmt.Errorf("invalid row and column count, the board must be a square")
	}

	return NewBoard(nbRows, nbCols, cells), nil
}

// Serialize a board to it's CSV string representation.
func serializeBoardToCsv(board *Board) (string, error) {
	// Create a CSV writer on top of a byte buffer.
	buffer := new(bytes.Buffer)
	csvWriter := csv.NewWriter(buffer)

	// Add a record to the writer for each row in the board.
	for iRow := 0; iRow < board.nbRows; iRow++ {
		// Create the row record.
		record := make([]string, board.nbCols)
		for iCol := 0; iCol < board.nbCols; iCol++ {
			cellId := (iRow * board.nbCols) + iCol
			color := board.cells[cellId]
			record[iCol] = strconv.Itoa(color)
		}

		// Add it to the writer.
		err := csvWriter.Write(record)
		if err != nil {
			return "", fmt.Errorf("unable to serialize the record as CSV: %w", err)
		}
	}

	// Flush the writer and return the resulting CSV string.
	csvWriter.Flush()
	csvStr := buffer.String()

	return csvStr, nil
}
