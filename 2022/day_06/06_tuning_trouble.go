package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
)

const startOfPacketMarkerLength int = 4
const startOfMessageMarkerLength int = 14

type void struct{}

var void_value void

func solveProblem(filePath string, solvePart1 bool) int {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	nbChar := countNbChar(csvReader, solvePart1)

	return nbChar
}

func hasDifferentChars(s string) bool {
	set := make(map[string]void)
	for _, char := range s {
		if _, ok := set[string(char)]; ok {
			return false
		} else {
			set[string(char)] = void_value
		}
	}
	return true
}

func countNbChar(data *csv.Reader, solvePart1 bool) int {
	line, err := data.Read()
	// deals with errors
	if err != nil {
		log.Fatal(err)
	}

	var markerLength int
	if solvePart1 {
		markerLength = startOfPacketMarkerLength
	} else {
		markerLength = startOfMessageMarkerLength
	}

	for i := 0; i < len(line[0])-markerLength; i++ {
		curr_marker := line[0][i : i+markerLength]
		if hasDifferentChars(curr_marker) {
			return i + markerLength
		}
	}
	return -1
}

func main() {
	inputArgPtr := flag.String("i", "./06_inputs.txt", "The input data (.csv file).")
	boolArgPtr := flag.Bool("solvePart1", false, "Whether to get the solution for Part 1 or for Part 2.")
	// Parse command line into the defined flags
	flag.Parse()

	nbChar := solveProblem(*inputArgPtr, *boolArgPtr)
	fmt.Printf("%d characters need to be processed before the first start-of-packet marker is detected.\n", nbChar)
}
