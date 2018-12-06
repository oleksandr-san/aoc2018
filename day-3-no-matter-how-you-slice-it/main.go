package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var claimRegexp = regexp.MustCompile(`#(\d+)\s+@\s+(\d+),(\d+):\s+(\d+)x(\d+)`)

type claim struct {
	id            string
	left, top     int
	width, height int
}

type claimHandler = func(*claim) bool

type cell struct {
	x, y   int
	claims []*claim
}

type plane struct {
	cells map[string]*cell
}

func makeKey(x, y int) string {
	return fmt.Sprintf("%d:%d", x, y)
}

func newPlane() *plane {
	return &plane{
		cells: map[string]*cell{},
	}
}

func (p *plane) addCell(x, y int, relatedClaim *claim) *cell {
	key := makeKey(x, y)
	existingCell, ok := p.cells[key]
	if ok {
		existingCell.claims = append(existingCell.claims, relatedClaim)
	} else {
		existingCell = &cell{x, y, []*claim{relatedClaim}}
		p.cells[key] = existingCell
	}
	return existingCell
}

func (p *plane) forEachCell(handler func(*cell)) {
	for _, c := range p.cells {
		handler(c)
	}
}

func parseClaim(line string) *claim {
	match := claimRegexp.FindStringSubmatch(line)
	if len(match) != 6 {
		return nil
	}

	id := match[1]
	left, _ := strconv.Atoi(match[2])
	top, _ := strconv.Atoi(match[3])
	width, _ := strconv.Atoi(match[4])
	height, _ := strconv.Atoi(match[5])

	return &claim{id, left, top, width, height}
}

func readClaims(r io.Reader, handler claimHandler) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if !handler(parseClaim(strings.TrimSpace(scanner.Text()))) {
			break
		}
	}
}

func calculateOverlappingCellsCount(r io.Reader) int {
	plane := newPlane()

	readClaims(r, func(c *claim) bool {
		for i := 0; i < c.width; i++ {
			for j := 0; j < c.height; j++ {
				plane.addCell(i+c.left, j+c.top, c)
			}
		}
		return true
	})

	overlappingCells := 0
	plane.forEachCell(func(c *cell) {
		if len(c.claims) > 1 {
			overlappingCells++
		}
	})
	return overlappingCells
}

func findNonOverlappingClaimID(r io.Reader) string {
	plane := newPlane()
	nonOverlappingClaims := map[*claim]bool{}

	readClaims(r, func(c *claim) bool {
		for i := 0; i < c.width; i++ {
			for j := 0; j < c.height; j++ {
				plane.addCell(i+c.left, j+c.top, c)
			}
		}
		nonOverlappingClaims[c] = true
		return true
	})

	plane.forEachCell(func(c *cell) {
		if len(c.claims) > 1 {
			for _, claim := range c.claims {
				delete(nonOverlappingClaims, claim)
			}
		}
	})

	for claim := range nonOverlappingClaims {
		return claim.id
	}

	return ""
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

	fmt.Printf("Non-overlapping claim is %s\n", findNonOverlappingClaimID(f))
	//fmt.Printf("Overlapping cells count is %d\n", calculateOverlappingCellsCount(f))
}
