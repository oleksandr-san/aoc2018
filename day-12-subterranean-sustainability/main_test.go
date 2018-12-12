package main

import (
	"reflect"
	"testing"
)

func TestExtendPots(t *testing.T) {
	tests := []struct {
		initialPots, expectedPots pots
	}{
		{
			initialPots:  pots{values: []byte("#..#")},
			expectedPots: pots{values: []byte("..#..#.."), firstPotNumber: -2},
		},
		{
			initialPots:  pots{values: []byte(".#..#.")},
			expectedPots: pots{values: []byte("..#..#.."), firstPotNumber: -1},
		},
		{
			initialPots:  pots{values: []byte("..#..#..")},
			expectedPots: pots{values: []byte("..#..#..")},
		},
	}

	for _, test := range tests {
		actualPots := extendPots(test.initialPots)
		if !reflect.DeepEqual(actualPots, test.expectedPots) {
			t.Errorf("Expected pots (%v) differ from the actual ones (%v)\n", test.expectedPots, actualPots)
		}
	}
}

func TestCalculateSumOfPlantedPotNumbers(t *testing.T) {
	tests := []struct {
		pots        pots
		expectedSum int
	}{
		{
			pots:        pots{values: []byte("#..#")},
			expectedSum: 3,
		},
		{
			pots:        pots{values: []byte(".#....##....#####...#######....#.#..##."), firstPotNumber: -3},
			expectedSum: 325,
		},
	}

	for _, test := range tests {
		actualSum := calculateSumOfPlantedPotNumbers(test.pots)
		if actualSum != test.expectedSum {
			t.Errorf("Expected pots (%d) differ from the actual ones (%d)\n", test.expectedSum, actualSum)
		}
	}
}

func TestEvaluatePotRule(t *testing.T) {
	tests := []struct {
		pots         pots
		idx          int
		expectedRule string
	}{
		{pots{values: []byte("#")}, 0, "..#.."},
		{pots{values: []byte(".")}, 0, "....."},
		{pots{values: []byte("#.#")}, 1, ".#.#."},
	}

	for _, test := range tests {
		actualRule := evaluatePotRule(test.pots, test.idx)
		if actualRule != test.expectedRule {
			t.Errorf("Expected rule (%s) differ from the actual one (%s)\n", test.expectedRule, actualRule)
		}
	}
}
