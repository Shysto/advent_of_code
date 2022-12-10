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

const left string = "L"
const up string = "U"
const right string = "R"
const down string = "D"
const sep string = " "

const part1NbKnots int = 2
const part2NbKnots int = 10

type Point struct {
	X, Y int
}

type Rope []Point

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func xor(X, Y bool) bool {
	return (X || Y) && !(X && Y)
}

func (p Point) toString() string {
	return "(" + strconv.Itoa(p.X) + "," + strconv.Itoa(p.Y) + ")"
}

func (p Point) isEqual(point Point) bool {
	return (p.X == point.X) && (p.Y == point.Y)
}

func (p Point) isAdjacent(neighbor Point) bool {
	isAdjacentHorizontal := abs(neighbor.X-p.X) == 1 && (neighbor.Y == p.Y)
	isAdjacentVertical := abs(neighbor.Y-p.Y) == 1 && (neighbor.X == p.X)
	return xor(isAdjacentHorizontal, isAdjacentVertical)
}

func (p Point) isDiagonal(neighbor Point) bool {
	return (abs(neighbor.X-p.X) == 1) && (abs(neighbor.Y-p.Y) == 1)
}

func (p Point) goLeft() (new_p Point) {
	new_p = Point{p.X, p.Y}
	new_p.X--
	return
}

func (p Point) goUp() (new_p Point) {
	new_p = Point{p.X, p.Y}
	new_p.Y++
	return
}

func (p Point) goRight() (new_p Point) {
	new_p = Point{p.X, p.Y}
	new_p.X++
	return
}

func (p Point) goDown() (new_p Point) {
	new_p = Point{p.X, p.Y}
	new_p.Y--
	return
}

func (p Point) updateHead(direction string) (new_p Point) {
	switch direction {
	case up:
		new_p = p.goUp()
	case down:
		new_p = p.goDown()
	case left:
		new_p = p.goLeft()
	case right:
		new_p = p.goRight()
	default:
		log.Fatalf("Unknow direction '%s'\n", direction)
	}
	return
}

func (p Point) updateTail(head Point) (new_p Point) {
	new_p = Point{p.X, p.Y}
	if p.isEqual(head) {
		return
	}

	if xor(p.X == head.X, p.Y == head.Y) {
		if p.X - head.X > 1 {
			new_p.X--
		} else if p.X - head.X < -1 {
			new_p.X++
		} else if p.Y - head.Y > 1 {
			new_p.Y--
		} else if p.Y - head.Y < -1 {
			new_p.Y++
		} 
	} else {
		if p.X - head.X > 1 {
			new_p.X--
			if p.Y < head.Y {
				new_p.Y++
			} else {
				new_p.Y--
			}
		} else if p.X - head.X < -1 {
			new_p.X++
			if p.Y < head.Y {
				new_p.Y++
			} else {
				new_p.Y--
			}
		} else if p.Y - head.Y > 1 {
			new_p.Y--
			if p.X < head.X {
				new_p.X++
			} else {
				new_p.X--
			}
		} else if p.Y - head.Y < -1 {
			new_p.Y++
			if p.X < head.X {
				new_p.X++
			} else {
				new_p.X--
			}
		} 
	}
	return 
}

func solveProblem(filePath string, solvePart1 bool) int {
	f, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	tailPositions := moveTail(scanner, solvePart1)
	//fmt.Println(grid)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return len(tailPositions)
}

func moveTail(scanner *bufio.Scanner, solvePart1 bool) map[string]int {
	var nbKnots int
	if solvePart1 {
		nbKnots = part1NbKnots
	} else {
		nbKnots = part2NbKnots
	}

	tailPositions := make(map[string]int)
	rope := make([]Point, nbKnots)
	for i := 0; i < nbKnots; i++ {
		rope[i] = Point{0, 0}
	}
	tailPositions[rope[nbKnots-1].toString()]++
	for scanner.Scan() {
		instruction := scanner.Text()
		split := strings.Split(instruction, sep)
		direction, nbSteps_string := split[0], split[1]
		nbSteps, err := strconv.Atoi(nbSteps_string)
		if err != nil {
			log.Fatal(err)
		}
		for nbSteps > 0 {
			//previousPosition := Point{rope[0].X, rope[0].Y}
			rope[0] = rope[0].updateHead(direction)
			//fmt.Printf("H: %v\t", rope[0].toString())

			for k := 1; k < nbKnots; k++ {
				rope[k] = rope[k].updateTail(rope[k-1])
				//fmt.Printf("%d: %v\t", k, rope[k].toString())
			}

			tailPositions[rope[nbKnots-1].toString()]++
			nbSteps--
			//fmt.Println()
		}
	}
	return tailPositions
}

func main() {
	inputArgPtr := flag.String("i", "./09_inputs.txt", "The input data (.txt file).")
	boolArgPtr := flag.Bool("solvePart1", false, "Whether to get the solution for Part 1 or for Part 2.")
	// Parse command line into the defined flags
	flag.Parse()

	numberPositions := solveProblem(*inputArgPtr, *boolArgPtr)

	fmt.Printf("The tail visited at least once %d positions.\n", numberPositions)
}
