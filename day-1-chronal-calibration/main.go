package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type valueHandler = func(int)

func readValues(r io.Reader, handler valueHandler) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		value, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil {
			return err
		}

		handler(value)
	}

	return nil
}

func calculateFinalFrequency(r io.Reader) (int, error) {
	frequency := 0
	err := readValues(r, func(change int) {
		frequency += change
	})

	return frequency, err
}

func findFirstRepeatedFrequency(r io.Reader) (int, error) {
	changes := make([]int, 0)
	err := readValues(r, func(change int) {
		changes = append(changes, change)
	})

	if err != nil {
		return 0, err
	}

	currentFrequency, seenFrequencies := 0, map[int]bool{0: true}
	for {
		for _, change := range changes {
			currentFrequency += change
			if _, seen := seenFrequencies[currentFrequency]; seen {
				return currentFrequency, nil
			}

			seenFrequencies[currentFrequency] = true
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments")
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	frequency, err := findFirstRepeatedFrequency(f)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Resulting frequency is %d\n", frequency)
}
