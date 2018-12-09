package main

import (
	"fmt"
	"regexp"
	"strconv"
)

var gameRuleRegexp = regexp.MustCompile(`(\d+) players; last marble is worth (\d+) points`)

type gameRules struct {
	playersCount    int
	lastMarbleValue int
}

func readGameRules(rawRules string) gameRules {
	rules := gameRules{}
	match := gameRuleRegexp.FindStringSubmatch(rawRules)

	if len(match) == 3 {
		rules.playersCount, _ = strconv.Atoi(match[1])
		rules.lastMarbleValue, _ = strconv.Atoi(match[2])
	}

	return rules
}

type marble struct {
	value      int
	next, prev *marble
}

func (m *marble) getNext(steps int) *marble {
	next := m
	for i := 0; i < steps; i++ {
		next = next.next
	}
	return next
}

func (m *marble) getPrev(steps int) *marble {
	prev := m
	for i := 0; i < steps; i++ {
		prev = prev.prev
	}
	return prev
}

type gameEngine struct {
	playerScores  []int
	currentPlayer int
	currentMarble *marble
}

func (e *gameEngine) play(rules gameRules) int {
	zeroMarble := &marble{}
	zeroMarble.prev, zeroMarble.next = zeroMarble, zeroMarble

	e.playerScores = make([]int, rules.playersCount)
	e.currentPlayer = 0
	e.currentMarble = zeroMarble

	nextMarbleValue := e.currentMarble.value + 1
	for ; nextMarbleValue <= rules.lastMarbleValue; nextMarbleValue++ {
		nextMarble := &marble{value: nextMarbleValue}

		if nextMarble.value%23 != 0 {
			firstMarble, secondMarble := e.currentMarble.getNext(1), e.currentMarble.getNext(2)

			nextMarble.prev, nextMarble.next = firstMarble, secondMarble
			firstMarble.next, secondMarble.prev = nextMarble, nextMarble
			e.currentMarble = nextMarble
		} else {
			removedMarble := e.currentMarble.getPrev(7)
			removedMarble.prev.next, removedMarble.next.prev = removedMarble.next, removedMarble.prev
			e.currentMarble = removedMarble.next

			e.playerScores[e.currentPlayer] += nextMarble.value + removedMarble.value
		}

		e.currentPlayer++
		if e.currentPlayer == rules.playersCount {
			e.currentPlayer = 0
		}
	}

	return e.maxScore()
}

func (e *gameEngine) maxScore() int {
	maxScore := 0
	for _, score := range e.playerScores {
		if maxScore < score {
			maxScore = score
		}
	}
	return maxScore
}

func main() {
	engine := gameEngine{}
	rules := readGameRules("441 players; last marble is worth 7103200 points")
	fmt.Printf("Winning player score is %d\n", engine.play(rules))
}
