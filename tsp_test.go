package main

import (
	"testing"
)

func TestTSP(t *testing.T) {

	assertCorrectLegCount := func(t *testing.T, got, want int) {
		t.Helper()
		if got != want {
			t.Errorf("got '%d' want '%d'", got, want)
		}
	}

	t.Run("Small", func(t *testing.T) {
		got := len(SolveShortestByShortest(10))
		want := 10
		assertCorrectLegCount(t, got, want)
	})
}