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

const crateLength int = 4
const crateLetterBegin string = "A"
const crateLetterEnd string = "Z"

const instructionMove string = "move "
const instructionFrom string = " from "
const instructionTo string = " to "

type Stack []string

// IsEmpty: check if stack is empty
func (stack *Stack) isEmpty() bool {
	return len(*stack) == 0
}

// Push a new value onto the stack (LIFO)
func (stack *Stack) push(str string) {
	*stack = append(*stack, str)
}

// Remove and return top element of stack (LIFO). Return false if stack is empty.
func (stack *Stack) pop() (string, bool) {
	if stack.isEmpty() {
		return "", false
	} else {
		last_idx := len(*stack) - 1   // Get the index of the top most element.
		element := (*stack)[last_idx] // Index into the slice and obtain the element.
		*stack = (*stack)[:last_idx]  // Remove it from the stack by slicing it off.
		return element, true
	}
}

func (stack *Stack) inverse() {
	for i, j := 0, len(*stack)-1; i < j; i, j = i+1, j-1 {
		(*stack)[i], (*stack)[j] = (*stack)[j], (*stack)[i]
	}
}

func initializeStacks(nbStacks int) []Stack {
	stacks := make([]Stack, nbStacks)
	for i := 0; i < nbStacks; i++ {
		var s Stack
		stacks[i] = s
	}
	return stacks
}

func solveProblem(filePath string, solvePart1 bool) string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	topCrates := getTopCrates(csvReader, solvePart1)

	return topCrates
}

func getInitialStacks(stacks []Stack, line string) {
	for i, char := range line {
		if string(char) >= crateLetterBegin && string(char) <= crateLetterEnd {
			var stackID int = i / crateLength
			stacks[stackID].push(string(char))
		}
	}
}

func moveCratesOnebyOne(nbCrates, fromStackIdx, toStackIdx int, stacks []Stack) {
	for i := 0; i < nbCrates; i++ {
		// unstack 1 crate
		crate, ok := stacks[fromStackIdx-1].pop()
		if !ok {
			log.Fatal("Cannot unstack an empty stack.")
		}
		// stack 1 crate
		stacks[toStackIdx-1].push(crate)
	}
}

func moveCratesTogether(nbCrates, fromStackIdx, toStackIdx int, stacks []Stack) {
	var tmp_stack Stack
	for i := 0; i < nbCrates; i++ {
		// unstack 1 crate
		crate, ok := stacks[fromStackIdx-1].pop()
		if !ok {
			log.Fatal("Cannot unstack an empty stack.")
		}
		// stack 1 crate
		tmp_stack.push(crate)
	}
	for !tmp_stack.isEmpty() {
		crate, ok := tmp_stack.pop()
		if !ok {
			log.Fatal("Cannot unstack an empty stack.")
		}
		stacks[toStackIdx-1].push(crate)
	}
}

func getTopCrates(data *csv.Reader, solvePart1 bool) string {
	count_pairs := ""
	// stacks initialization
	line, err := data.Read()
	// deals with errors
	if err != nil {
		log.Fatal(err)
	}
	var nbStacks int = (len(line[0]) + 1) / crateLength
	stacks := initializeStacks(nbStacks)
	delimiter := " "
	for i := 1; i < nbStacks; i++ {
		delimiter += fmt.Sprint(i) + strings.Repeat(" ", crateLength-1)
	}
	delimiter += fmt.Sprint(nbStacks) + " "
	//fmt.Println(delimiter)
	// stacks initial configuration
	for {
		if line[0] == delimiter {
			break
		}
		getInitialStacks(stacks, line[0])

		line, err = data.Read()
		// deals with errors
		if err != nil {
			log.Fatal(err)
		}
	}
	// inverse stacks
	for i := 0; i < len(stacks); i++ {
		stacks[i].inverse()
	}
	//fmt.Printf("Stacks '%v'.\n", stacks)

	// follow instructions
	for {
		line, err = data.Read()
		// deals with errors
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		split_elems := strings.Split(line[0], instructionMove)
		split_elems = strings.Split(split_elems[1], instructionFrom)
		nbCrates, err := strconv.Atoi(split_elems[0])
		if err != nil {
			log.Fatal(err)
		}
		split_elems = strings.Split(split_elems[1], instructionTo)
		fromStackIdx, err := strconv.Atoi(split_elems[0])
		if err != nil {
			log.Fatal(err)
		}
		toStackIdx, err := strconv.Atoi(split_elems[1])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("move %d from %d to %d\n", nbCrates, fromStackIdx, toStackIdx)
		if solvePart1 {
			moveCratesOnebyOne(nbCrates, fromStackIdx, toStackIdx, stacks)
		} else {
			moveCratesTogether(nbCrates, fromStackIdx, toStackIdx, stacks)
		}
		//fmt.Printf("Stacks '%v'.\n", stacks)
	}

	// get top crates
	for i := 0; i < len(stacks); i++ {
		top_crate, ok := stacks[i].pop()
		if ok {
			count_pairs += top_crate
		}
	}
	return count_pairs
}

func main() {
	inputArgPtr := flag.String("i", "./05_inputs.txt", "The input data (.csv file).")
	boolArgPtr := flag.Bool("solvePart1", false, "Whether to get the solution for Part 1 or for Part 2.")
	// Parse command line into the defined flags
	flag.Parse()

	topCrates := solveProblem(*inputArgPtr, *boolArgPtr)
	fmt.Printf("Top crates are '%s'.\n", topCrates)
}
