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

const sep string = ","
const arrow string = " -> "
const rock string = "#"
const sand string = "o"
const air string = "."
const origin string = "+"
const infinity int = 1000

var offsetX int

type Point struct {
	X, Y int
}

func (p *Point) equals(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

var sandOrigin Point = Point{500, 0}

type Scan [][]string

func (s *Scan) toString() string {
	repr := ""
	for i := 0; i < len(*s); i++ {
		for j := 0; j < len((*s)[0]); j++ {
			repr += (*s)[i][j]
		}
		repr += "\n"
	}
	return repr
}

func (s *Scan) canMove(p Point) (int, Point) {
	if p.Y >= len(*s)-1 { // flows out
		return -1, Point{}
	}
	if (*s)[p.Y+1][p.X-offsetX] == air { // falls vertically
		return 1, Point{p.X, p.Y + 1}
	}

	if p.X-offsetX < 1 { // flows out
		return -1, Point{}
	}
	if (*s)[p.Y+1][p.X-offsetX-1] == air { // falls diagonally on the left
		return 1, Point{p.X - 1, p.Y + 1}
	}

	if p.X-offsetX >= len((*s)[0])-1 { // flows out
		return -1, Point{}
	}
	if (*s)[p.Y+1][p.X-offsetX+1] == air { // falls diagonally on the right
		return 1, Point{p.X + 1, p.Y + 1}
	}

	return 0, p // comes to rest
}

func solveProblem(filePath string, solvePart1 bool) int {
	f, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	if !solvePart1 {
		scan := buildInfiniteScan(scanner)
		//fmt.Println(scan.toString())
		sandUnit := fillInfiniteSand(scan)
		return sandUnit
	} else {
		scan := initScan(scanner)
		//fmt.Println(scan.toString())
		sandUnit := fillSand(scan)
		return sandUnit
	}
}

func initScan(scanner *bufio.Scanner) Scan {
	scan := Scan{}
	lines := [][]Point{}
	minX, maxX, maxY := int(^uint(0)>>1), -1, -1

	// Get min and max coordinates
	i := 0
	for scanner.Scan() {
		lines = append(lines, []Point{})
		data := scanner.Text()
		points := strings.Split(data, arrow)
		var prev_point Point
		for _, p := range points {
			coord := strings.Split(p, sep)
			x, _ := strconv.Atoi(coord[0])
			y, _ := strconv.Atoi(coord[1])
			if x < minX {
				minX = x
			} else if x > maxX {
				maxX = x
			}
			if y > maxY {
				maxY = y
			}

			// Get all points in line
			new_point := Point{x, y}
			if prev_point != (Point{}) {
				for prev_point.X < new_point.X {
					prev_point.X++
					lines[i] = append(lines[i], prev_point)
				}
				for prev_point.X > new_point.X {
					prev_point.X--
					lines[i] = append(lines[i], prev_point)
				}
				for prev_point.Y < new_point.Y {
					prev_point.Y++
					lines[i] = append(lines[i], prev_point)
				}
				for prev_point.Y > new_point.Y {
					prev_point.Y--
					lines[i] = append(lines[i], prev_point)
				}
			}
			lines[i] = append(lines[i], new_point)
			prev_point = Point{new_point.X, new_point.Y}
		}
	}

	// Initialize scan with air
	offsetX = minX
	for i := 0; i < maxY+1; i++ {
		row := []string{}
		for j := 0; j < maxX-minX+1; j++ {
			row = append(row, air)
		}
		scan = append(scan, row)
	}

	// Initialize lines
	for _, line := range lines {
		for _, point := range line {
			scan[point.Y][point.X-offsetX] = rock
		}
	}

	// Initialize sand origin
	scan[sandOrigin.Y][sandOrigin.X-offsetX] = origin

	//fmt.Printf("X in [%d, %d], Y in [0, %d]\n", minX, maxX, maxY)

	return scan
}

func buildInfiniteScan(scanner *bufio.Scanner) map[Point]string {
	scan := make(map[Point]string)
	minX, maxX, maxY := int(^uint(0)>>1), -1, -1

	for scanner.Scan() {
		data := scanner.Text()
		points := strings.Split(data, arrow)
		var prev_point Point
		for _, p := range points {
			coord := strings.Split(p, sep)
			x, _ := strconv.Atoi(coord[0])
			y, _ := strconv.Atoi(coord[1])
			if x < minX {
				minX = x
			} else if x > maxX {
				maxX = x
			}
			if y > maxY {
				maxY = y
			}

			// Get all points in line
			new_point := Point{x, y}
			if prev_point != (Point{}) {
				for prev_point.X < new_point.X {
					prev_point.X++
					scan[prev_point] = rock
				}
				for prev_point.X > new_point.X {
					prev_point.X--
					scan[prev_point] = rock
				}
				for prev_point.Y < new_point.Y {
					prev_point.Y++
					scan[prev_point] = rock
				}
				for prev_point.Y > new_point.Y {
					prev_point.Y--
					scan[prev_point] = rock
				}
			}
			scan[new_point] = rock
			prev_point = Point{new_point.X, new_point.Y}
		}
	}

	// Add infinite ground
	for i := minX - infinity/2; i < maxX+infinity/2+1; i++ {
		scan[Point{i, maxY + 2}] = rock
	}

	return scan
}

func fillSand(s Scan) int {
	sandUnit := 0
	sandCanMove, sandNextPosition := s.canMove(sandOrigin)
	for sandCanMove != -1 || !sandNextPosition.equals(sandOrigin) {
		if sandCanMove == 1 {
			sandCanMove, sandNextPosition = s.canMove(sandNextPosition)
		} else if sandCanMove == 0 {
			//if sandUnit%10 == 0 {
			//	fmt.Println(s.toString())
			//}
			s[sandNextPosition.Y][sandNextPosition.X-offsetX] = sand
			sandUnit++
			sandCanMove, sandNextPosition = s.canMove(sandOrigin)
		} else {
			break
		}
	}
	fmt.Println(s.toString())
	return sandUnit
}

func canMove(p Point, s map[Point]string) (int, Point) {
	below_point, ok := s[Point{p.X, p.Y + 1}]
	if !ok || below_point == air { // falls vertically
		return 1, Point{p.X, p.Y + 1}
	}

	below_left_point, ok := s[Point{p.X - 1, p.Y + 1}]
	if !ok || below_left_point == air { // falls diagonally on the left
		return 1, Point{p.X - 1, p.Y + 1}
	}

	below_right_point, ok := s[Point{p.X + 1, p.Y + 1}]
	if !ok || below_right_point == air { // falls diagonally on the right
		return 1, Point{p.X + 1, p.Y + 1}
	}

	return 0, p // comes to rest
}

func fillInfiniteSand(s map[Point]string) int {
	sandUnit := 0
	sandCanMove, sandNextPosition := canMove(sandOrigin, s)
	for !sandNextPosition.equals(sandOrigin) {
		if sandCanMove == 1 {
			sandCanMove, sandNextPosition = canMove(sandNextPosition, s)
		} else {
			//if sandUnit%10 == 0 {
			//	fmt.Println(s.toString())
			//}
			s[sandNextPosition] = sand
			sandUnit++
			sandCanMove, sandNextPosition = canMove(sandOrigin, s)
		}
	}
	return sandUnit + 1
}

func main() {
	inputArgPtr := flag.String("i", "./14_inputs.txt", "The input data (.txt file).")
	boolArgPtr := flag.Bool("solvePart1", false, "Whether to get the solution for Part 1 or for Part 2.")
	// Parse command line into the defined flags
	flag.Parse()

	answer := solveProblem(*inputArgPtr, *boolArgPtr)

	if *boolArgPtr {
		fmt.Printf("%d units of sand come to rest before sand starts flowing into the abyss below.\n", answer)
	} else {
		fmt.Printf("%d units of sand come to rest before the source of the sand becomes blocked.\n", answer)
	}
}
