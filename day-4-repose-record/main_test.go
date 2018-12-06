package main

import (
	"os"
	"testing"
)

func TestFindMostSleepyMinuteData(t *testing.T) {
	check := func(t *testing.T, path string, expected int) {
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		actual := findMostSleepyMinuteData(f)
		if expected != actual {
			t.Errorf("Expected result (%d) differs from actual (%d)\n", expected, actual)
		}
	}

	t.Run("Data from data.txt", func(t *testing.T) {
		check(t, "data.txt", 103720)
	})
}

func TestFindMostSleepyGuardData(t *testing.T) {
	check := func(t *testing.T, path string, expected int) {
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		actual := findMostSleepyGuardData(f)
		if expected != actual {
			t.Errorf("Expected result (%d) differs from actual (%d)\n", expected, actual)
		}
	}

	t.Run("Data from data.txt", func(t *testing.T) {
		check(t, "data.txt", 110913)
	})
}
