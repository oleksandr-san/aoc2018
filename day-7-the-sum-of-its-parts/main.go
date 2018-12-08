package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
)

var stepDependencyRegexp = regexp.MustCompile(`Step (\S+?) must be finished before step (\S+?) can begin\.`)

type step struct {
	name                    string
	dependencies, dependend []*step
}

type stepsByName []*step

func (s stepsByName) Len() int {
	return len(s)
}

func (s stepsByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s stepsByName) Less(i, j int) bool {
	return s[i].name < s[j].name
}

type steps map[string]*step

func (s steps) ensureStep(name string) *step {
	st, ok := s[name]
	if !ok {
		st = &step{name: name}
		s[name] = st
	}
	return st
}

func (s steps) addStepDependency(from, to string) {
	fromStep := s.ensureStep(from)
	toStep := s.ensureStep(to)

	fromStep.dependencies = append(fromStep.dependencies, toStep)
	toStep.dependend = append(toStep.dependend, fromStep)
}

func (s steps) sortInterals() {
	for _, st := range s {
		sort.Sort(stepsByName(st.dependencies))
		sort.Sort(stepsByName(st.dependend))
	}
}

func collectStepsData(steps steps) (incompletedSteps map[*step]bool, availableSteps []*step) {
	incompletedSteps = map[*step]bool{}
	for _, step := range steps {
		incompletedSteps[step] = true
		if len(step.dependencies) == 0 {
			availableSteps = append(availableSteps, step)
		}
	}

	return
}

func updateAvailableSteps(availableSteps []*step, incompletedSteps map[*step]bool, completedStep *step) []*step {
	isCompleted := func(step *step) bool {
		_, ok := incompletedSteps[step]
		return !ok
	}

	areCompleted := func(steps []*step) bool {
		for _, step := range steps {
			if !isCompleted(step) {
				return false
			}
		}
		return true
	}

	for _, dependentStep := range completedStep.dependend {
		if !isCompleted(dependentStep) && areCompleted(dependentStep.dependencies) {
			availableSteps = append(availableSteps, dependentStep)
		}
	}

	return availableSteps
}

func calculateStepsOrder(steps steps) string {
	var stepNames bytes.Buffer

	incompletedSteps, availableSteps := collectStepsData(steps)
	for len(availableSteps) != 0 {
		sort.Sort(stepsByName(availableSteps))

		step := availableSteps[0]
		availableSteps = availableSteps[1:]

		stepNames.WriteString(step.name)

		delete(incompletedSteps, step)
		availableSteps = updateAvailableSteps(availableSteps, incompletedSteps, step)
	}

	return stepNames.String()
}

type task struct {
	step    *step
	endTime int
}

type workQueue struct {
	tasks         []task
	currentTime   int
	workersCount  int
	workerLatency int
}

func (q *workQueue) hasWorkers() bool {
	return len(q.tasks) < q.workersCount
}

func (q *workQueue) assignSteps(steps []*step) []*step {
	for q.hasWorkers() && len(steps) > 0 {
		step := steps[0]
		steps = steps[1:]

		task := task{step, q.currentTime + q.workerLatency + int(rune(step.name[0])-rune('A')+1)}
		q.tasks = append(q.tasks, task)
	}

	return steps
}

func (q *workQueue) waitForSome() (completedSteps []*step) {
	earliestEndTime := 0
	for _, task := range q.tasks {
		if earliestEndTime == 0 || earliestEndTime > task.endTime {
			earliestEndTime = task.endTime
		}
	}
	q.currentTime = earliestEndTime

	for i := len(q.tasks) - 1; i >= 0; i-- {
		if q.tasks[i].endTime == q.currentTime {
			completedSteps = append(completedSteps, q.tasks[i].step)
			q.tasks = append(q.tasks[:i], q.tasks[i+1:]...)
		}
	}

	return
}

func simulateStepsExecution(steps steps, workersCount, workerLatency int) int {
	queue := workQueue{workersCount: workersCount, workerLatency: workerLatency}

	incompletedSteps, availableSteps := collectStepsData(steps)
	for len(incompletedSteps) != 0 {
		availableSteps = queue.assignSteps(availableSteps)

		for _, completedStep := range queue.waitForSome() {
			delete(incompletedSteps, completedStep)
			availableSteps = updateAvailableSteps(availableSteps, incompletedSteps, completedStep)
		}
	}

	return queue.currentTime
}

func readSteps(r io.Reader) steps {
	steps := steps{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		match := stepDependencyRegexp.FindStringSubmatch(scanner.Text())
		if len(match) != 3 {
			continue
		}

		steps.addStepDependency(match[2], match[1])
	}

	steps.sortInterals()
	return steps
}

func readStepsFromFile(path string) (steps, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return readSteps(f), nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments")
	}

	steps, err := readStepsFromFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Steps order: %s\n", calculateStepsOrder(steps))
	fmt.Printf("Step execution time: %d\n", simulateStepsExecution(steps, 5, 60))
}
