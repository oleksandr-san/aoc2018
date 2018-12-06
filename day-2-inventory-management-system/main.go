package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type lineHandler = func(string) bool

func readLines(r io.Reader, handler lineHandler) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if !handler(strings.TrimSpace(scanner.Text())) {
			break
		}
	}
}

func checkLetterRepeats(line string) (hasTwoLetter bool, hasThreeLetter bool) {
	repeatCnt := map[rune]int{}
	for _, letter := range line {
		repeatCnt[letter]++
	}

	for _, cnt := range repeatCnt {
		if cnt == 2 {
			hasTwoLetter = true
		} else if cnt == 3 {
			hasThreeLetter = true
		}
	}

	return
}

func calculateChecksumOfIDs(r io.Reader) int {
	twoLetterCnt, threeLetterCnt := 0, 0

	readLines(r, func(line string) bool {
		hasTwoLetter, hasThreeLetter := checkLetterRepeats(line)
		if hasTwoLetter {
			twoLetterCnt++
		}
		if hasThreeLetter {
			threeLetterCnt++
		}
		return true
	})

	return twoLetterCnt * threeLetterCnt
}

func findCommonPart(lhs, rhs string) string {
	lhsRunes, rhsRunes := []rune(lhs), []rune(rhs)
	if len(lhsRunes) != len(rhsRunes) {
		return ""
	}

	var commonPart bytes.Buffer
	for i, r := range lhsRunes {
		if r == rhsRunes[i] {
			commonPart.WriteRune(r)
		}
	}
	return commonPart.String()
}

func findCommonIDPart(r io.Reader) (commonPart string) {
	var lines []string

	readLines(r, func(line string) bool {
		for _, previousLine := range lines {
			localCommonPart := findCommonPart(line, previousLine)
			if len(localCommonPart) == len(line)-1 {
				commonPart = localCommonPart
				return false
			}
		}
		lines = append(lines, line)
		return true
	})

	return
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments")
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	//fmt.Printf("Resulting checksum is %d\n", calculateChecksumOfIDs(f))
	fmt.Printf("Resulting common part is %s\n", findCommonIDPart(f))
}
