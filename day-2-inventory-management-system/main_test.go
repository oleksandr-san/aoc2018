package main

import (
	"os"
	"testing"
)

func TestCalculateChecksumOfIDs(t *testing.T) {
	check := func(t *testing.T, path string, expectedChecksum int) {
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		actualChecksum := calculateChecksumOfIDs(f)
		if expectedChecksum != actualChecksum {
			t.Errorf("Expected checksum (%d) differs from actual (%d)\n", expectedChecksum, actualChecksum)
		}
	}

	t.Run("Data from data.txt", func(t *testing.T) {
		check(t, "data.txt", 4693)
	})
}

func TestFindCommonIDPart(t *testing.T) {
	check := func(t *testing.T, path string, expectedPart string) {
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		actualPart := findCommonIDPart(f)
		if expectedPart != actualPart {
			t.Errorf("Expected common part (%s) differs from actual (%s)\n", expectedPart, actualPart)
		}
	}

	t.Run("Data from data.txt", func(t *testing.T) {
		check(t, "data.txt", "pebjqsalrdnckzfihvtxysomg")
	})
}
