package main

import "testing"

func TestCalculateMetadataSum(t *testing.T) {
	check := func(t *testing.T, path string, expected int) {
		node, err := readNodeTree(path)
		if err != nil {
			t.Fatal(err)
		}

		actual := calculateMetadataSum(node)
		if expected != actual {
			t.Errorf("Expected result (%d) differs from actual (%d)\n", expected, actual)
		}
	}

	t.Run("Data from data_1.txt", func(t *testing.T) {
		check(t, "data_1.txt", 138)
	})

	t.Run("Data from data_2.txt", func(t *testing.T) {
		check(t, "data_2.txt", 40848)
	})
}

func TestCalculateSpecificSum(t *testing.T) {
	check := func(t *testing.T, path string, expected int) {
		node, err := readNodeTree(path)
		if err != nil {
			t.Fatal(err)
		}

		actual := calculateSpecificSum(node)
		if expected != actual {
			t.Errorf("Expected result (%d) differs from actual (%d)\n", expected, actual)
		}
	}

	t.Run("Data from data_1.txt", func(t *testing.T) {
		check(t, "data_1.txt", 66)
	})

	t.Run("Data from data_2.txt", func(t *testing.T) {
		check(t, "data_2.txt", 34466)
	})
}
