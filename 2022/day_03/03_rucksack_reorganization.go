package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var Priority = map[string]int{
	"a": 1, "A": 27,
	"b": 2, "B": 28,
	"c": 3, "C": 29,
	"d": 4, "D": 30,
	"e": 5, "E": 31,
	"f": 6, "F": 32,
	"g": 7, "G": 33,
	"h": 8, "H": 34,
	"i": 9, "I": 35,
	"j": 10, "J": 36,
	"k": 11, "K": 37,
	"l": 12, "L": 38,
	"m": 13, "M": 39,
	"n": 14, "N": 40,
	"o": 15, "O": 41,
	"p": 16, "P": 42,
	"q": 17, "Q": 43,
	"r": 18, "R": 44,
	"s": 19, "S": 45,
	"t": 20, "T": 46,
	"u": 21, "U": 47,
	"v": 22, "V": 48,
	"w": 23, "W": 49,
	"x": 24, "X": 50,
	"y": 25, "Y": 51,
	"z": 26, "Z": 52,
}

const nbChar int = 52
const nbElfsInGroups int = 3

func solveProblem(filePath string, solvePart1 bool) int {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	sumPriorities := 0
	if solvePart1 {
		sumPriorities = getRucksackSumPriorities(csvReader)
	} else {
		sumPriorities = getBadgeSumPriorities(csvReader)
	}

	return sumPriorities
}

func getRucksackSumPriorities(data *csv.Reader) int {
	// initialization
	sumPriorities := 0

	for {
		rec, err := data.Read()
		// deals with errors
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		var rucksack string = rec[0]
		// split compartments in half
		rucksack_len := len(rucksack)
		first_compartment := rucksack[:rucksack_len/2]
		second_compartment := rucksack[rucksack_len/2:]

		for _, char := range first_compartment {
			if strings.Contains(second_compartment, string(char)) {
				sumPriorities += Priority[string(char)]
				//fmt.Printf("The common element in both compartments is '%c' with priority %d.\n", char, Priority[string(char)])
				break
			}
		}

	}
	//fmt.Println()
	return sumPriorities
}

func mapKey(m map[string]int, value int) (key string, ok bool) {
	// retrieve key associated to value in map
	for k, v := range m {
		if v == value {
			key = k
			ok = true
			return
		}
	}
	return
}

func IntersectStrings(str1, str2 string) string {
	// Gets intersection (=common characters) of two strings
	var inStr1 [nbChar]bool
	for k, v := range Priority {
		inStr1[v-1] = strings.Contains(str1, k)
	}

	var inStr2 [nbChar]bool
	for k, v := range Priority {
		inStr2[v-1] = strings.Contains(str2, k)
	}

	commonString := ""
	for i := 0; i < nbChar; i++ {
		if inStr1[i] && inStr2[i] {
			key, ok := mapKey(Priority, i+1)
			if !ok {
				log.Fatal("Value does not exist in map")
			}
			commonString += key
		}
	}
	return commonString
}

func getBadgeSumPriorities(data *csv.Reader) int {
	// initialization
	sumPriorities := 0

	rec, err := data.ReadAll()
	// deals with errors
	if err != nil {
		log.Fatal(err)
	}

	var nb_groups int = len(rec) / nbElfsInGroups
	//fmt.Printf("There are %d groups.\n", nb_groups)
	for i := 0; i < nb_groups; i++ {
		var commonElems string = rec[i*nbElfsInGroups][0]
		for j := 1; j < nbElfsInGroups; j++ {
			commonElems = IntersectStrings(commonElems, rec[i*nbElfsInGroups+j][0])
		}
		//fmt.Printf("Common elements in group #%d: %v.\n", i, commonElems)
		sumPriorities += Priority[string(commonElems[0])]
	}

	//fmt.Println()
	return sumPriorities
}

func main() {
	inputArgPtr := flag.String("i", "./03_inputs.csv", "The input data (.csv file).")
	boolArgPtr := flag.Bool("solvePart1", false, "Whether to get the solution for Part 1 or for Part 2.")
	// Parse command line into the defined flags
	flag.Parse()

	sumPriorities := solveProblem(*inputArgPtr, *boolArgPtr)
	fmt.Printf("The sum of priorities is %d.\n", sumPriorities)
}
