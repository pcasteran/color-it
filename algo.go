package main

// Algorithm is the function type that will be used by all the implementations.
type Algorithm func(board *Board) ([]int, error)

// Available algorithm implementations.
var implementations = map[string]Algorithm{
	"dummy": dummy,
}
