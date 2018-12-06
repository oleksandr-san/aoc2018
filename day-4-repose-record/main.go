package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var guardIDRegexp = regexp.MustCompile(`#(\d+)`)

type rawRecord struct {
	date, body string
}

type rawRecordsByDate []*rawRecord

func (s rawRecordsByDate) Len() int {
	return len(s)
}

func (s rawRecordsByDate) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s rawRecordsByDate) Less(i, j int) bool {
	return s[i].date < s[j].date
}

func parseRawRecord(line string) *rawRecord {
	dateEndIdx := strings.IndexRune(line, ']')
	if dateEndIdx == -1 {
		return nil
	}

	return &rawRecord{line[:dateEndIdx+1], line[dateEndIdx+2:]}
}

func readRawRecords(r io.Reader, handler func(*rawRecord) bool) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		rawRecord := parseRawRecord(strings.TrimSpace(scanner.Text()))
		if rawRecord != nil && !handler(rawRecord) {
			break
		}
	}
}

func readShiftRecords(rawRecordList []*rawRecord, handler func(int, []int)) {
	guardID, fallAsleepTime := -1, -1
	var minutesAsleep []int

	for _, r := range rawRecordList {
		if strings.HasPrefix(r.body, "Guard") {
			match := guardIDRegexp.FindStringSubmatch(r.body)
			if len(match) == 2 {
				guardID, _ = strconv.Atoi(match[1])
			}
		} else if strings.HasPrefix(r.body, "falls asleep") {
			fallAsleepTime, _ = strconv.Atoi(r.date[15:17])
		} else if strings.HasPrefix(r.body, "wakes up") {
			wakeUpTime, _ := strconv.Atoi(r.date[15:17])
			for i := fallAsleepTime; i < wakeUpTime; i++ {
				minutesAsleep = append(minutesAsleep, i)
			}
			handler(guardID, minutesAsleep)
			minutesAsleep, fallAsleepTime = nil, -1
		}
	}
}

func findGuardWithMaxTimeAsleep(shiftRecords map[int][][]int) (int, int) {
	guardID, maxTimeAsleep := -1, -1

	for id, timeAsleepList := range shiftRecords {
		timeAsleep := 0
		for _, minutesAsleep := range timeAsleepList {
			timeAsleep += len(minutesAsleep)
		}

		if maxTimeAsleep < timeAsleep {
			maxTimeAsleep = timeAsleep
			guardID = id
		}
	}

	return guardID, maxTimeAsleep
}

func findMostSleepyMinute(timeAsleep [][]int) (int, int) {
	sleepCnt := map[int]int{}

	for _, minutesAsleep := range timeAsleep {
		for _, minute := range minutesAsleep {
			sleepCnt[minute]++
		}
	}

	maxKey, maxVal := -1, -1
	for key, val := range sleepCnt {
		if maxVal < val {
			maxKey, maxVal = key, val
		}
	}
	return maxKey, maxVal
}

func collectShiftRecords(r io.Reader) map[int][][]int {
	var rawRecordList []*rawRecord
	readRawRecords(r, func(r *rawRecord) bool {
		rawRecordList = append(rawRecordList, r)
		return true
	})

	sort.Sort(rawRecordsByDate(rawRecordList))

	shiftRecords := map[int][][]int{}
	readShiftRecords(rawRecordList, func(guardID int, minutesAsleep []int) {
		shiftRecords[guardID] = append(shiftRecords[guardID], minutesAsleep)
	})

	return shiftRecords
}

func findMostSleepyMinuteData(r io.Reader) int {
	shiftRecords := collectShiftRecords(r)
	guardID, _ := findGuardWithMaxTimeAsleep(shiftRecords)
	mostSleepyMinute, _ := findMostSleepyMinute(shiftRecords[guardID])
	return guardID * mostSleepyMinute
}

func findMostSleepyGuardData(r io.Reader) int {
	shiftRecords := collectShiftRecords(r)
	guardID, maxMinute, maxCnt := -1, -1, -1
	for id, timeAsleep := range shiftRecords {
		minute, cnt := findMostSleepyMinute(timeAsleep)
		if maxCnt < cnt {
			guardID, maxMinute, maxCnt = id, minute, cnt
		}
	}

	return guardID * maxMinute
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

	//fmt.Printf("Most sleepy guard data is %d\n", findMostSleepyGuardData(f))
	fmt.Printf("Most sleepy minute data is %d\n", findMostSleepyGuardData(f))
}
