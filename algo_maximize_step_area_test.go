package main

import "testing"

func BenchmarkMaximizeStepArea(b *testing.B) {
	benchmarkImplementation(b, maximizeStepArea, "samples/30_30_3-1.csv")
}

func BenchmarkMaximizeStepAreaDeep(b *testing.B) {
	benchmarkImplementation(b, maximizeStepAreaDeep, "samples/30_30_3-1.csv")
}
