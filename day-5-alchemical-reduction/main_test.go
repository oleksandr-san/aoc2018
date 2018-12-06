package main

import (
	"testing"
)

func TestCalculatePolymerAfterReactions(t *testing.T) {
	check := func(t *testing.T, path string, expected int) {
		polymer, err := readPolymer(path)
		if err != nil {
			t.Fatal(err)
		}

		actual := len(calculatePolymerAfterReactions(polymer))
		if expected != actual {
			t.Errorf("Expected result (%d) differs from actual (%d)\n", expected, actual)
		}
	}

	t.Run("Data from data.txt", func(t *testing.T) {
		check(t, "data.txt", 10384)
	})
}

func TestCalculateShortesPolymerWithOneTypeRemoval(t *testing.T) {
	check := func(t *testing.T, path string, expected int) {
		polymer, err := readPolymer(path)
		if err != nil {
			t.Fatal(err)
		}

		actual := len(calculateShortesPolymerWithOneTypeRemoval(polymer))
		if expected != actual {
			t.Errorf("Expected result (%d) differs from actual (%d)\n", expected, actual)
		}
	}

	t.Run("Data from data.txt", func(t *testing.T) {
		check(t, "data.txt", 5412)
	})
}
