package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const startPos string = "S"
const endPos string = "E"
const lowestElevation string = "a"
const maxClimbDiff int = 1

type Point struct {
	X, Y, Priority int
}

func (p Point) toString() string {
	return "(" + strconv.Itoa(p.X) + "," + strconv.Itoa(p.Y) + ")"
}

func (p Point) getNeighbors() (neighbors []Point) {
	neighbors = append(neighbors, Point{p.X + 1, p.Y, 0})
	neighbors = append(neighbors, Point{p.X - 1, p.Y, 0})
	neighbors = append(neighbors, Point{p.X, p.Y + 1, 0})
	neighbors = append(neighbors, Point{p.X, p.Y - 1, 0})
	return
}

type PriorityQueue []Point

func (pq *PriorityQueue) isEmpty() bool {
	return len(*pq) == 0
}

func (pq *PriorityQueue) push(elem Point) {
	*pq = append(*pq, elem)
}

func (pq *PriorityQueue) pop() (Point, bool) {
	if pq.isEmpty() {
		return Point{-1, -1, -1}, false
	}
	elem := (*pq)[0]
	if len(*pq) == 1 {
		*pq = []Point{}
	} else {
		*pq = (*pq)[1:]
	}
	return elem, true
}

type Grid [][]string

func (g Grid) isEmpty() bool {
	return len(g) == 0
}

func (g Grid) canMove(dest_i, dest_j int, curr_value string) bool {
	if g.isEmpty() {
		return false
	}
	if dest_i < 0 || dest_j < 0 || dest_i >= len(g) || dest_j >= len(g[0]) {
		return false
	}

	if curr_value == startPos {
		curr_value = "a"
	}

	dest_value := g[dest_i][dest_j]

	if dest_value == endPos {
		dest_value = "z"
	}

	return int(dest_value[0])-int(curr_value[0]) <= maxClimbDiff
}

func (g Grid) findShortestPath(start []Point) int {
	pq := PriorityQueue{}
	for _, p := range start {
		pq.push(p)
	}
	visited := map[string]bool{}

	for !pq.isEmpty() {
		curr_point, ok := pq.pop()
		if !ok {
			break
		}
		curr_val := g[curr_point.X][curr_point.Y]

		if curr_val == endPos {
			return curr_point.Priority
		}

		for _, p := range curr_point.getNeighbors() {
			_, ok := visited[p.toString()]
			if !ok && g.canMove(p.X, p.Y, curr_val) {
				visited[p.toString()] = true
				pq.push(Point{p.X, p.Y, curr_point.Priority + 1})
			}
		}

	}
	return -1

}

func solveProblem(filePath string, solvePart1 bool) int {
	f, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	topo := parseGrid(scanner, solvePart1)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var start []Point
	if solvePart1 {
		start = []Point{findStartingPosition(topo)}
	} else {
		start = findAllStartingPositions(topo)
	}

	return topo.findShortestPath(start)
}

func parseGrid(scanner *bufio.Scanner, solvePart1 bool) (topo Grid) {
	for scanner.Scan() {
		row := scanner.Text()
		topo = append(topo, strings.Split(row, ""))
	}
	return
}

func findStartingPosition(topo Grid) Point {
	for i, row := range topo {
		for j, char := range row {
			if string(char) == startPos {
				return Point{i, j, 0}
			}
		}
	}
	return Point{-1, -1, -1}
}

func findAllStartingPositions(topo Grid) []Point {
	starts := []Point{}
	for i, row := range topo {
		for j, char := range row {
			if string(char) == lowestElevation || string(char) == startPos {
				starts = append(starts, Point{i, j, 0})
			}
		}
	}
	return starts
}

func main() {
	inputArgPtr := flag.String("i", "./12_inputs.txt", "The input data (.txt file).")
	boolArgPtr := flag.Bool("solvePart1", false, "Whether to get the solution for Part 1 or for Part 2.")
	// Parse command line into the defined flags
	flag.Parse()

	nbSteps := solveProblem(*inputArgPtr, *boolArgPtr)

	if *boolArgPtr {
		fmt.Printf("Moving from the starting position to the location that should get the best signal requires at best %d steps.\n", nbSteps)
	} else {
		fmt.Printf("Moving starting from any square with elevation a to the location that should get the best signal requires as best %d steps.\n", nbSteps)
	}
}
