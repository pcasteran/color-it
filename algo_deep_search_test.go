package main

import "testing"

func BenchmarkDeepSearch(b *testing.B) {
	benchmarkImplementation(b, deepSearch, "samples/30_30_3-1.csv")
}
