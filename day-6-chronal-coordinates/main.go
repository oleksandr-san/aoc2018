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

type coordinate struct {
	x, y int
}

type region struct {
	center      coordinate
	visibleArea int
	isInfinite  bool
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func manhattanDistance(p1, p2 coordinate) int {
	return abs(p1.x-p2.x) + abs(p1.y-p2.y)
}

func readCoordinates(r io.Reader) ([]coordinate, error) {
	coordinates := []coordinate{}
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		rawCoordinates := strings.Split(scanner.Text(), ",")
		if len(rawCoordinates) != 2 {
			continue
		}

		x, err := strconv.Atoi(strings.TrimSpace(rawCoordinates[0]))
		if err != nil {
			return nil, err
		}

		y, err := strconv.Atoi(strings.TrimSpace(rawCoordinates[1]))
		if err != nil {
			return nil, err
		}

		coordinates = append(coordinates, coordinate{x, y})
	}

	return coordinates, nil
}

func readCoordinatesFrom(path string) ([]coordinate, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	return readCoordinates(f)
}

func calculateVisibleGrid(coordinates []coordinate) (topLeft, bottomRight coordinate) {
	var maxXFound, maxYFound, minXFound, minYFound bool

	for _, c := range coordinates {
		if !maxXFound || bottomRight.x < c.x {
			maxXFound = true
			bottomRight.x = c.x
		}

		if !minYFound || bottomRight.y > c.y {
			minYFound = true
			bottomRight.y = c.y
		}

		if !minXFound || topLeft.x > c.x {
			minXFound = true
			topLeft.x = c.x
		}

		if !maxYFound || topLeft.y < c.y {
			maxYFound = true
			topLeft.y = c.y
		}
	}

	return
}

func forEachCoordinateInVisibleGrid(coordinates []coordinate, handler func(coordinate, bool)) {
	topLeft, bottomRight := calculateVisibleGrid(coordinates)

	for x := topLeft.x; x <= bottomRight.x; x++ {
		for y := bottomRight.y; y <= topLeft.y; y++ {
			isEdge := x == topLeft.x || x == bottomRight.y || y == topLeft.y || y == bottomRight.y
			handler(coordinate{x, y}, isEdge)
		}
	}
}

func createRegions(coordinates []coordinate) (regions []*region) {
	for _, c := range coordinates {
		regions = append(regions, &region{center: c})
	}
	return
}

func findClosestRegion(c coordinate, regions []*region) *region {
	var closestRegions []*region
	minDistance := -1

	for _, r := range regions {
		distance := manhattanDistance(c, r.center)

		if minDistance == -1 || minDistance > distance {
			closestRegions = []*region{r}
			minDistance = distance
		} else if distance == minDistance {
			closestRegions = append(closestRegions, r)
		}
	}

	if len(closestRegions) == 1 {
		return closestRegions[0]
	}

	return nil
}

func calculateSumOfManhattanDistances(c coordinate, regions []*region) int {
	sum := 0
	for _, region := range regions {
		sum += manhattanDistance(c, region.center)
	}
	return sum
}

func findLargestFiniteRegion(regions []*region) (largestRegion *region) {
	for _, region := range regions {
		if region.isInfinite {
			continue
		}

		if largestRegion == nil || largestRegion.visibleArea < region.visibleArea {
			largestRegion = region
		}
	}

	return
}

func findLargestFiniteArea(coordinates []coordinate) int {
	regions := createRegions(coordinates)

	forEachCoordinateInVisibleGrid(coordinates, func(c coordinate, isEdge bool) {
		if region := findClosestRegion(c, regions); region != nil {
			region.visibleArea++
			if isEdge {
				region.isInfinite = true
			}
		}
	})

	largestRegion := findLargestFiniteRegion(regions)
	if largestRegion != nil {
		return largestRegion.visibleArea
	}

	return 0
}

func findAreaWithSumDistancesLessThan(coordinates []coordinate, n int) int {
	area := 0
	regions := createRegions(coordinates)

	forEachCoordinateInVisibleGrid(coordinates, func(c coordinate, isEdge bool) {
		if calculateSumOfManhattanDistances(c, regions) < n {
			area++
		}
	})

	return area
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments")
	}

	coordinates, err := readCoordinatesFrom(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(findLargestFiniteArea(coordinates))
	fmt.Println(findAreaWithSumDistancesLessThan(coordinates, 10000))
}
