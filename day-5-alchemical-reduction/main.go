package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"unicode"
)

var fetchUnitType = unicode.ToLower

func readPolymer(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if scanner.Scan() {
		return scanner.Text(), nil
	}

	return "", nil
}

func fetchUnitTypes(polymer string) (unitTypes []rune) {
	uniqueTypes := map[rune]bool{}

	for _, unit := range polymer {
		unitType := fetchUnitType(unit)
		if _, ok := uniqueTypes[unitType]; !ok {
			unitTypes = append(unitTypes, unitType)
			uniqueTypes[unitType] = true
		}
	}

	return
}

func findReactingUnitIndexes(polymer string) (indexes []int) {
	var reactionFound bool
	var prevUnit rune

	for i, unit := range polymer {
		if i == 0 || reactionFound {
			reactionFound = false
			prevUnit = unit
			continue
		}

		if unit != prevUnit && fetchUnitType(unit) == fetchUnitType(prevUnit) {
			indexes = append(indexes, i-1)
			reactionFound = true
		}

		prevUnit = unit
	}

	return
}

func removeReactingUnits(polymer string, indexes []int) string {
	var newPolymer bytes.Buffer
	var reactionFound bool

	for i, unit := range polymer {
		if reactionFound {
			reactionFound = false
			continue
		}

		if len(indexes) != 0 && i == indexes[0] {
			reactionFound = true
			indexes = indexes[1:]
			continue
		}

		newPolymer.WriteRune(unit)
	}

	return newPolymer.String()
}

func removeUnitsOfType(polymer string, unitType rune) string {
	var newPolymer bytes.Buffer

	for _, unit := range polymer {
		if fetchUnitType(unit) != unitType {
			newPolymer.WriteRune(unit)
		}
	}

	return newPolymer.String()
}

func calculatePolymerAfterReactions(polymer string) string {
	for {
		reactingUnits := findReactingUnitIndexes(polymer)
		if len(reactingUnits) == 0 {
			break
		}

		polymer = removeReactingUnits(polymer, reactingUnits)
	}

	return polymer
}

func calculateShortesPolymerWithOneTypeRemoval(polymer string) string {
	shortestPolymer := ""

	for _, unitType := range fetchUnitTypes(polymer) {
		newPolymer := calculatePolymerAfterReactions(removeUnitsOfType(polymer, unitType))
		if len(shortestPolymer) == 0 || len(shortestPolymer) > len(newPolymer) {
			shortestPolymer = newPolymer
		}
	}

	return shortestPolymer
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments")
	}

	polymer, err := readPolymer(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("The Length of polymer after reactions is %d\n", len(calculatePolymerAfterReactions(polymer)))
	fmt.Printf("The shortest possible polymer length (with one unit type removal) is %d\n", len(calculateShortesPolymerWithOneTypeRemoval(polymer)))
}
