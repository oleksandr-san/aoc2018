package main

import "testing"

func TestWinnerScore(t *testing.T) {
	tests := []struct {
		rules         string
		expectedScore int
	}{
		{"10 players; last marble is worth 1618 points", 8317},
		{"13 players; last marble is worth 7999 points", 146373},
		{"17 players; last marble is worth 1104 points", 2764},
		{"21 players; last marble is worth 6111 points", 54718},
		{"30 players; last marble is worth 5807 points", 37305},
		{"441 players; last marble is worth 71032 points", 393229},
		{"441 players; last marble is worth 7103200 points", 3273405195},
	}

	engine := gameEngine{}
	for _, test := range tests {
		actualScore := engine.play(readGameRules(test.rules))
		if actualScore != test.expectedScore {
			t.Fatalf("Expected winner score (%d) differs from actual one (%d)\n", test.expectedScore, actualScore)
		}
	}
}
