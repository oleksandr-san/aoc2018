package main

import (
	"os"
	"testing"
)

func TestCalculateFinalFrequency(t *testing.T) {
	check := func(t *testing.T, path string, expectFrequency int) {
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		actualFrequency, err := calculateFinalFrequency(f)
		if err != nil {
			t.Fatal(err)
		}

		if expectFrequency != actualFrequency {
			t.Errorf("Expected frequency (%d) differs from actual (%d)\n", expectFrequency, actualFrequency)
		}
	}

	t.Run("Data from data.txt", func(t *testing.T) {
		check(t, "data.txt", 466)
	})
}

func TestFindFirstRepeatedFrequency(t *testing.T) {
	check := func(t *testing.T, path string, expectFrequency int) {
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		actualFrequency, err := findFirstRepeatedFrequency(f)
		if err != nil {
			t.Fatal(err)
		}

		if expectFrequency != actualFrequency {
			t.Errorf("Expected frequency (%d) differs from actual (%d)\n", expectFrequency, actualFrequency)
		}
	}

	t.Run("Data from data.txt", func(t *testing.T) {
		check(t, "data.txt", 750)
	})
}
