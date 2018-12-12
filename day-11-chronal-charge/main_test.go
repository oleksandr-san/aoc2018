package main

import "testing"

func TestCalculateCellPowerLevel(t *testing.T) {
	tests := []struct {
		fuelCell     cell
		serialNumber int
		powerLevel   int
	}{
		{cell{122, 79}, 57, -5},
		{cell{217, 196}, 39, 0},
		{cell{101, 153}, 71, 4},
	}

	for _, test := range tests {
		actualPowerLevel := calculateCellPowerLevel(test.fuelCell, test.serialNumber)
		if actualPowerLevel != test.powerLevel {
			t.Errorf(
				"For test case (%v) expected power level (%d) differs from the actual one (%d)\n",
				test,
				test.powerLevel,
				actualPowerLevel,
			)
		}
	}
}

func TestCalculateRegionPowerLevel(t *testing.T) {
	tests := []struct {
		topLeft      cell
		serialNumber int
		powerLevel   int
	}{
		{cell{33, 45}, 18, 29},
		{cell{21, 61}, 42, 30},
	}

	for _, test := range tests {
		actualPowerLevel := calculateRegionPowerLevel(test.topLeft, test.serialNumber, 3)
		if actualPowerLevel != test.powerLevel {
			t.Errorf(
				"For test case (%v) expected power level (%d) differs from the actual one (%d)\n",
				test,
				test.powerLevel,
				actualPowerLevel,
			)
		}
	}
}

func TestFindMostPowerfulRegion(t *testing.T) {
	tests := []struct {
		serialNumber int
		topLeft      cell
	}{
		{18, cell{33, 45}},
		{42, cell{21, 61}},
		{8199, cell{235, 87}},
	}

	for _, test := range tests {
		actualTopLeft, _ := findMostPowerfulRegion(test.serialNumber, 3)
		if actualTopLeft != test.topLeft {
			t.Errorf(
				"For test case (%v) expected cell (%v) differs from the actual one (%v)\n",
				test,
				test.topLeft,
				actualTopLeft,
			)
		}
	}
}
