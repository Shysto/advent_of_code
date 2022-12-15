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

const comma string = ", "
const sep string = ": "
const xIs string = "x="
const yIs string = "y="
const rowY int = 2000000
const globalMin int = 0
const globalMax int = 4000000
const tuningFreqFactor int = 4000000

var minX, minY, maxX, maxY int = int(^uint(0) >> 1), int(^uint(0) >> 1), -int(^uint(0)>>1) - 1, -int(^uint(0)>>1) - 1

type Point struct {
	X, Y int
}

type Sensor struct {
	Coordinates Point
	Beacon      Point
	Radius      int
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func manhattanDistance(p, q Point) int {
	return abs(q.X-p.X) + abs(q.Y-p.Y)
}

func (p *Point) equals(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

func (s *Sensor) setRadius() {
	s.Radius = manhattanDistance(s.Coordinates, s.Beacon)
}

func (s *Sensor) isInsideCoverage(p Point) int {
	if manhattanDistance(s.Coordinates, p) < s.Radius {
		return 1
	} else if manhattanDistance(s.Coordinates, p) == s.Radius {
		return 0
	}
	return -1
}

func (s *Sensor) getPointsTouchingContour() []Point {
	cx, cy := s.Coordinates.X, s.Coordinates.Y
	r := s.Radius + 1
	points := []Point{}
	curr_pt := Point{cx - r, cy}
	for curr_pt.Y < cy-r {
		curr_pt.X++
		curr_pt.Y--
		points = append(points, curr_pt)
	}
	for curr_pt.X < cx+r {
		curr_pt.X++
		curr_pt.Y++
		points = append(points, curr_pt)
	}
	for curr_pt.Y < cy+r {
		curr_pt.X--
		curr_pt.Y++
		points = append(points, curr_pt)
	}
	for curr_pt.X < cx-r {
		curr_pt.X--
		curr_pt.Y--
		points = append(points, curr_pt)
	}
	return points
}

func solveProblem(filePath string, solvePart1 bool) int {
	f, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	sensors := parseSensors(scanner)
	if solvePart1 {
		return countImpossiblePositionsInRow(rowY, sensors)
	}

	possiblePosition := findPossiblePosition(sensors)
	return possiblePosition.X*tuningFreqFactor + possiblePosition.Y
}

func parseSensors(scanner *bufio.Scanner) []Sensor {
	sensors := []Sensor{}
	for scanner.Scan() {
		data := scanner.Text()
		firstSplit := strings.Split(data, xIs)
		sensor_str := strings.Split(firstSplit[1], comma+yIs)
		xSensor_str := sensor_str[0]
		ySensor_str := strings.Split(sensor_str[1], sep)[0]
		xSensor, _ := strconv.Atoi(xSensor_str)
		ySensor, _ := strconv.Atoi(ySensor_str)
		//fmt.Printf("Sensor is at (%s,%s)\n", xSensor_str, ySensor_str)
		beacon_str := strings.Split(firstSplit[2], comma+yIs)
		xBeacon_str := beacon_str[0]
		yBeacon_str := beacon_str[1]
		xBeacon, _ := strconv.Atoi(xBeacon_str)
		yBeacon, _ := strconv.Atoi(yBeacon_str)
		//fmt.Printf("Beacon is at (%s,%s)\n", xBeacon_str, yBeacon_str)

		sensor := Sensor{Coordinates: Point{xSensor, ySensor}, Beacon: Point{xBeacon, yBeacon}}
		sensor.setRadius()
		//fmt.Printf("Sensor: %+v\n", sensor)
		sensors = append(sensors, sensor)

		// Get min/max coordinates
		tmp_minX := min(xSensor-sensor.Radius, xBeacon)
		tmp_maxX := max(xSensor+sensor.Radius, xBeacon)
		tmp_minY := min(ySensor-sensor.Radius, yBeacon)
		tmp_maxY := max(ySensor+sensor.Radius, yBeacon)
		if tmp_minX < minX {
			minX = tmp_minX
		}
		if tmp_minY < minY {
			minY = tmp_minY
		}
		if tmp_maxX > maxX {
			maxX = tmp_maxX
		}
		if tmp_maxY > maxY {
			maxY = tmp_maxY
		}
	}
	//fmt.Printf("minX=%d, maxX=%d, minY=%d, maxY=%d\n", minX, maxX, minY, maxY)
	return sensors
}

func canBeaconBePresent(p Point, sensors []Sensor) bool {
	for _, s := range sensors {
		if s.isInsideCoverage(p) == 0 && s.Beacon.equals(p) {
			//fmt.Printf("Point %+v is the beacon of Sensor %+v.\n", p, s)
			return true
		} else if s.isInsideCoverage(p) > -1 {
			return false
		}
	}
	return true
}

func countImpossiblePositionsInRow(row int, sensors []Sensor) int {
	point := Point{X: minX, Y: row}
	count := 0
	for point.X <= maxX {
		if !canBeaconBePresent(point, sensors) {
			count++
		}
		point.X++
	}
	return count
}

func findPossiblePosition(sensors []Sensor) Point {
	for i := 0; i < len(sensors); i++ {
		// A possible solution must be outside the edges of the sensor coverage area
		for _, point := range sensors[i].getPointsTouchingContour() {
			// Test if the point coordinates are in [0,4000000]
			if point.X >= globalMin && point.X <= globalMax && point.Y >= globalMin && point.Y <= globalMax {
				isSolution := true
				for j := 0; j < len(sensors); j++ {
					// If a point is inside the coverage of at least one sensor, it is not a possible position
					if i != j && sensors[j].isInsideCoverage(point) > -1 {
						isSolution = false
						break
					}
				}
				// This point is outside the coverage of all sensors, it is a possible position
				if isSolution {
					return point
				}
			}
		}
	}
	return Point{}
}

func main() {
	inputArgPtr := flag.String("i", "./15_inputs.txt", "The input data (.txt file).")
	boolArgPtr := flag.Bool("solvePart1", false, "Whether to get the solution for Part 1 or for Part 2.")
	// Parse command line into the defined flags
	flag.Parse()

	answer := solveProblem(*inputArgPtr, *boolArgPtr)

	if *boolArgPtr {
		fmt.Printf("In the row where y=%d, %d positions cannot contain a beacon.\n", rowY, answer)
	} else {
		fmt.Printf("The tuning frequency is %d.\n", answer)
	}
}
