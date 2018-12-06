package main

import (
	"testing"
)

func TestFindLargestFiniteArea(t *testing.T) {
	check := func(t *testing.T, path string, expectedArea int) {
		coordinates, err := readCoordinatesFrom(path)
		if err != nil {
			t.Fatal(err)
		}

		actualArea := findLargestFiniteArea(coordinates)
		if expectedArea != actualArea {
			t.Errorf("Expected area (%d) differs from actual (%d)\n", expectedArea, actualArea)
		}
	}

	t.Run("Data from data_1.txt", func(t *testing.T) {
		check(t, "data_1.txt", 17)
	})

	t.Run("Data from data_2.txt", func(t *testing.T) {
		check(t, "data_2.txt", 4976)
	})
}

func TestFindAreaWithSumDistancesLessThan(t *testing.T) {
	check := func(t *testing.T, path string, n, expectedArea int) {
		coordinates, err := readCoordinatesFrom(path)
		if err != nil {
			t.Fatal(err)
		}

		actualArea := findAreaWithSumDistancesLessThan(coordinates, n)
		if expectedArea != actualArea {
			t.Errorf("Expected area (%d) differs from actual (%d)\n", expectedArea, actualArea)
		}
	}

	t.Run("Data from data_1.txt", func(t *testing.T) {
		check(t, "data_1.txt", 32, 16)
	})

	t.Run("Data from data_2.txt", func(t *testing.T) {
		check(t, "data_2.txt", 10000, 46462)
	})
}
