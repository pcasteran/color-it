package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func readInputFile(filePath string) (*Board, error) {
	// Load the raw string file content.
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open the input file: %w", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			// TODO
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
	if len(cells) != (nbRows * nbRows) {
		return nil, fmt.Errorf("invalid row and column count, the board must be a square")
	}

	return NewBoard(nbRows, nbRows /* nbCols=nbRows as the board is a square */, cells), nil
}
