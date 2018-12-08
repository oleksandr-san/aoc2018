package main

import "testing"

func TestCalculateStepsOrder(t *testing.T) {
	check := func(t *testing.T, path, expected string) {
		steps, err := readStepsFromFile(path)
		if err != nil {
			t.Fatal(err)
		}

		actual := calculateStepsOrder(steps)
		if expected != actual {
			t.Errorf("Expected result (%s) differs from actual (%s)\n", expected, actual)
		}
	}

	t.Run("Data from data_1.txt", func(t *testing.T) {
		check(t, "data_1.txt", "CABDFE")
	})

	t.Run("Data from data_2.txt", func(t *testing.T) {
		check(t, "data_2.txt", "EUGJKYFQSCLTWXNIZMAPVORDBH")
	})
}

func TestSimulateStepsExecution(t *testing.T) {
	check := func(t *testing.T, path string, expected int, workersCount, stepLatency int) {
		steps, err := readStepsFromFile(path)
		if err != nil {
			t.Fatal(err)
		}

		actual := simulateStepsExecution(steps, workersCount, stepLatency)
		if expected != actual {
			t.Errorf("Expected result (%d) differs from actual (%d)\n", expected, actual)
		}
	}

	t.Run("Data from data_1.txt", func(t *testing.T) {
		check(t, "data_1.txt", 15, 2, 0)
	})

	t.Run("Data from data_2.txt", func(t *testing.T) {
		check(t, "data_2.txt", 1014, 5, 60)
	})
}
