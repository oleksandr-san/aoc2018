package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	spreadRuleLength = 5
	plant            = '#'
	empty            = '.'
)

type pots struct {
	values         []byte
	firstPotNumber int
}

func (p pots) String() string {
	return fmt.Sprintf("{firstPotNumber=%d, pots=%s}", p.firstPotNumber, string(p.values))
}

type spreadRules map[string]byte

func readConfig(r io.Reader) (pots pots, rules spreadRules, err error) {
	scanner := bufio.NewScanner(r)
	if !scanner.Scan() {
		err = errors.New("unexpected EOF")
		return
	}

	text := scanner.Text()
	statePrefix := "initial state: "
	if !strings.HasPrefix(text, statePrefix) {
		err = errors.New("invalid file format")
		return
	}

	pots.values = []byte(text[len(statePrefix):])
	rules = spreadRules{}

	for scanner.Scan() {
		text := scanner.Text()
		if len(text) < spreadRuleLength+5 {
			continue
		}

		rules[text[:spreadRuleLength]] = text[spreadRuleLength+4]
	}

	return
}

func extendPots(originalPots pots) (extendedPots pots) {
	extendedPots.firstPotNumber = originalPots.firstPotNumber

	for i := 0; i < spreadRuleLength/2; i++ {
		if originalPots.values[i] == plant {
			for j := 0; j < spreadRuleLength/2-i; j++ {
				extendedPots.values = append(extendedPots.values, empty)
				extendedPots.firstPotNumber--
			}
			break
		}
	}

	extendedPots.values = append(extendedPots.values, originalPots.values...)

	for i := 0; i < spreadRuleLength/2; i++ {
		if originalPots.values[len(originalPots.values)-1-i] == plant {
			for j := 0; j < spreadRuleLength/2-i; j++ {
				extendedPots.values = append(extendedPots.values, empty)
			}
			break
		}
	}

	return
}

func evaluatePotRule(pots pots, idx int) string {
	var bytes bytes.Buffer
	for i := 0; i < spreadRuleLength/2; i++ {
		spreadIdx := idx - spreadRuleLength/2 + i
		if spreadIdx >= 0 && spreadIdx < len(pots.values) {
			bytes.WriteByte(pots.values[spreadIdx])
		} else {
			bytes.WriteByte(empty)
		}
	}

	bytes.WriteByte(pots.values[idx])

	for i := 0; i < spreadRuleLength/2; i++ {
		spreadIdx := idx + 1 + i
		if spreadIdx >= 0 && spreadIdx < len(pots.values) {
			bytes.WriteByte(pots.values[spreadIdx])
		} else {
			bytes.WriteByte(empty)
		}
	}
	return bytes.String()
}

func evaluateNextGeneration(originalPots pots, rules spreadRules) pots {
	extendedPots := extendPots(originalPots)
	newPots := pots{firstPotNumber: extendedPots.firstPotNumber}

	for i := range extendedPots.values {
		if newPot, ok := rules[evaluatePotRule(extendedPots, i)]; ok {
			newPots.values = append(newPots.values, newPot)
		} else {
			newPots.values = append(newPots.values, empty)
		}
	}

	return newPots
}

func calculateSum(pots pots) int {
	sum := 0
	for i, value := range pots.values {
		if value == plant {
			sum += pots.firstPotNumber + i
		}
	}
	return sum
}

func calculateSumAfterNGenerations(
	pots pots, rules spreadRules, generationsCnt int,
) int {
	for i := 0; i < generationsCnt; i++ {
		pots = evaluateNextGeneration(pots, rules)
	}

	return calculateSum(pots)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments")
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	initialPots, spreadRules, err := readConfig(f)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(
		"After 20 generations the sum is",
		calculateSumAfterNGenerations(initialPots, spreadRules, 20),
	)

	// This solution is taken from https://www.reddit.com/r/adventofcode/comments/a5eztl/2018_day_12_solutions/ebm4c9d
	pots := initialPots
	lastSum, diff := 0, 0
	for i := 0; i < 2000; i++ {
		pots = evaluateNextGeneration(pots, spreadRules)
		sum := calculateSum(pots)

		diff = sum - lastSum
		lastSum = sum
	}

	fmt.Println(
		"After 50 billion generations the sum is",
		calculateSum(pots)+(50000000000-2000)*diff,
	)
}
