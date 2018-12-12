package main

import "fmt"

type cell struct {
	x, y int
}

func calculateCellPowerLevel(cell cell, serialNumber int) int {
	return (((cell.x+10)*cell.y+serialNumber)*(cell.x+10))/100%10 - 5
}

func calculateRegionPowerLevel(topLeft cell, serialNumber, size int) int {
	powerLevel := 0
	for x := topLeft.x; x < topLeft.x+size; x++ {
		for y := topLeft.y; y < topLeft.y+size; y++ {
			powerLevel += calculateCellPowerLevel(cell{x, y}, serialNumber)
		}
	}
	return powerLevel
}

func findMostPowerfulRegion(serialNumber, size int) (topLeft cell, powerLevel int) {
	maxFound := false

	for x := 1; x <= 300; x++ {
		for y := 1; y <= 300; y++ {
			if x+size < 300 && y+size < 300 {
				currentPowerLevel := calculateRegionPowerLevel(cell{x, y}, serialNumber, size)
				if !maxFound || powerLevel < currentPowerLevel {
					topLeft, powerLevel = cell{x, y}, currentPowerLevel
					maxFound = true
				}
			}
		}
	}

	return
}

func findMostPowerfulRegionAnySize(serialNumber int) (topLeft cell, powerLevel, size int) {
	maxFound := false

	for i := 1; i <= 300; i++ {
		currentTopLeft, currentPowerLevel := findMostPowerfulRegion(serialNumber, i)
		if !maxFound || powerLevel < currentPowerLevel {
			topLeft, powerLevel, size = currentTopLeft, currentPowerLevel, i
			maxFound = true
		}
	}

	return
}

func main() {
	performCalculations := func(serialNumber int) {
		topLeft, _ := findMostPowerfulRegion(serialNumber, 3)
		fmt.Println("Most powerful region with serial number", serialNumber, "of size 3 is", topLeft)

		topLeft, powerLevel, size := findMostPowerfulRegionAnySize(serialNumber)
		fmt.Println("Most powerful region with serial number", serialNumber, "is", topLeft, "with power level", powerLevel, "of size", size)
	}

	performCalculations(18)
	performCalculations(42)
	performCalculations(8199)
}
