package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
)

var pointRegexp = regexp.MustCompile(`position=<\s*(-?\d+),\s+(-?\d+)> velocity=<\s*(-?\d+),\s+(-?\d+)>`)

type point struct {
	x, y int
}

type velocity struct {
	x, y int
}

func updatePositions(points []point, velocities []velocity, timeDelta int) {
	for i := 0; i < len(points); i++ {
		points[i].x += velocities[i].x * timeDelta
		points[i].y += velocities[i].y * timeDelta
	}
}

type showPolicy struct {
	empty, point rune
}

func printGrid(writer io.Writer, points []point, policy showPolicy) error {
	w := bufio.NewWriter(writer)

	pointSet := map[point]bool{}
	for _, point := range points {
		pointSet[point] = true
	}

	topLeft, bottomRight := calculateVisibleGrid(points)

	for y := bottomRight.y; y <= topLeft.y; y++ {
		for x := topLeft.x; x <= bottomRight.x; x++ {

			point := point{x, y}
			if _, ok := pointSet[point]; ok {
				w.WriteRune(policy.point)
			} else {
				w.WriteRune(policy.empty)
			}
		}

		w.WriteRune('\n')
	}

	w.Flush()
	return nil
}

func calculateVisibleGrid(points []point) (topLeft, bottomRight point) {
	var maxXFound, maxYFound, minXFound, minYFound bool

	for _, point := range points {
		if !maxXFound || bottomRight.x < point.x {
			maxXFound = true
			bottomRight.x = point.x
		}

		if !minYFound || bottomRight.y > point.y {
			minYFound = true
			bottomRight.y = point.y
		}

		if !minXFound || topLeft.x > point.x {
			minXFound = true
			topLeft.x = point.x
		}

		if !maxYFound || topLeft.y < point.y {
			maxYFound = true
			topLeft.y = point.y
		}
	}

	return
}

func readPoints(reader io.Reader) (points []point, velocities []velocity, err error) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		match := pointRegexp.FindStringSubmatch(scanner.Text())
		if len(match) != 5 {
			continue
		}

		point, velocity := point{}, velocity{}
		if point.x, err = strconv.Atoi(match[1]); err != nil {
			return
		}
		if point.y, err = strconv.Atoi(match[2]); err != nil {
			return
		}
		if velocity.x, err = strconv.Atoi(match[3]); err != nil {
			return
		}
		if velocity.y, err = strconv.Atoi(match[4]); err != nil {
			return
		}

		points = append(points, point)
		velocities = append(velocities, velocity)
	}

	return
}

func readPointsFromFile(path string) ([]point, []velocity, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	defer f.Close()
	return readPoints(f)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments")
	}

	points, velocities, err := readPointsFromFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	/*
		I guessed the right grid size from the 3rd attempt using binary search approach :)
		More elaborated solution use some of this ideas:
			a) some interactive rewinding features (e.g., going forward/backward in time on
			   the left/right arrow key press and redrawing the grid)
			b) running a text detection algorithm on each point set (e.g. OpenCV's EAST or Google Text Recognition API)
	*/

	for i := 0; ; i++ {
		topLeft, bottomRight := calculateVisibleGrid(points)
		if topLeft.y-bottomRight.y <= 62 && bottomRight.x-topLeft.x <= 62 {
			printGrid(os.Stdout, points, showPolicy{empty: '.', point: '#'})
			fmt.Printf("%d second(s) passed\n", i)
			break
		}

		updatePositions(points, velocities, 1)
	}
}
