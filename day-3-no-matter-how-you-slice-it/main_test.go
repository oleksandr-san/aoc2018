package main

import (
	"os"
	"testing"
)

func TestCalculateOverlappingCellsCount(t *testing.T) {
	check := func(t *testing.T, path string, expected int) {
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		actual := calculateOverlappingCellsCount(f)
		if expected != actual {
			t.Errorf("Expected common part (%d) differs from actual (%d)\n", expected, actual)
		}
	}

	t.Run("Data from data.txt", func(t *testing.T) {
		check(t, "data.txt", 116491)
	})
}

func TestFindNonOverlappingClaimID(t *testing.T) {
	check := func(t *testing.T, path string, expectedID string) {
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		actualID := findNonOverlappingClaimID(f)
		if expectedID != actualID {
			t.Errorf("Expected id (%s) differs from actual (%s)\n", expectedID, actualID)
		}
	}

	t.Run("Data from data.txt", func(t *testing.T) {
		check(t, "data.txt", "707")
	})
}
