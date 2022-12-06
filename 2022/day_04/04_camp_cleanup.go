package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const rangeDelimiter string = "-"

func solveProblem(filePath string, solvePart1 bool) int {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	nbPairs := countPairs(csvReader, solvePart1)

	return nbPairs
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

func makeRange(min, max int) []int {
	// makes a range of integers from a min and a max value
	if min > max {
		var empty_range []int
		return empty_range
	}
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func parseRange(section, delimiter string) [2]int {
	// parse "int1-int2" into a [int1, int2]
	range_ := strings.Split(section, delimiter)
	range_start, err := strconv.Atoi(range_[0])
	if err != nil {
		log.Fatal(err)
	}
	range_end, err := strconv.Atoi(range_[1])
	if err != nil {
		log.Fatal(err)
	}
	return [2]int{range_start, range_end}
}

func intersect(range1, range2 [2]int) []int {
	// gets intersection of two range of integers
	return makeRange(max(range1[0], range2[0]), min(range1[1], range2[1]))
}

func isFullyContained(section1, section2 string) bool {
	// whether a range is fully contained in another
	range1 := parseRange(section1, rangeDelimiter)
	range2 := parseRange(section2, rangeDelimiter)
	return len(intersect(range1, range2)) == min(len(makeRange(range1[0], range1[1])), len(makeRange(range2[0], range2[1])))
}

func areOverlapping(section1, section2 string) bool {
	// whether the intersection of two range is empty
	range1 := parseRange(section1, rangeDelimiter)
	range2 := parseRange(section2, rangeDelimiter)
	return len(intersect(range1, range2)) > 0
}

func countPairs(data *csv.Reader, solvePart1 bool) int {
	// initialization
	count_pairs := 0

	for {
		rec, err := data.Read()
		// deals with errors
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		var sections_elf1 string = rec[0]
		var sections_elf2 string = rec[1]

		if solvePart1 && isFullyContained(sections_elf1, sections_elf2) {
			count_pairs++
		}

		if !solvePart1 && areOverlapping(sections_elf1, sections_elf2) {
			count_pairs++
		}
	}
	return count_pairs
}

func main() {
	inputArgPtr := flag.String("i", "./04_inputs.csv", "The input data (.csv file).")
	boolArgPtr := flag.Bool("solvePart1", false, "Whether to get the solution for Part 1 or for Part 2.")
	// Parse command line into the defined flags
	flag.Parse()

	nbPairs := solveProblem(*inputArgPtr, *boolArgPtr)
	fmt.Printf("There are %d assignment pairs in which one range fully contains the other.\n", nbPairs)
}
