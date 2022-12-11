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

const addx string = "addx "
const noop string = "noop"

const addx_cycle int = 2

const firstStrengthMeasure int = 20
const measureEveryCycle int = 40

var X int = 1

func sumSignalStrengths(slice []int) (sum int) {
	sum = 0
	for i, elem := range slice {
		sum += elem * (firstStrengthMeasure + i*measureEveryCycle)
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
	signalStrengths := measureSignalStrength(scanner, solvePart1)
	//fmt.Println(grid)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return sumSignalStrengths(signalStrengths)
}

func measureSignalStrength(scanner *bufio.Scanner, solvePart1 bool) []int {
	signalStrengths := []int{}
	i := 1
	insideAddxForNCycles := -1
	addToX := 0
	for {
		if (i == firstStrengthMeasure) || (i-firstStrengthMeasure)%measureEveryCycle == 0 {
			fmt.Printf("Append X=%d ath the %dth cycle\n", X, i)
			signalStrengths = append(signalStrengths, X)
		}
		if insideAddxForNCycles > 0 {
			if insideAddxForNCycles < addx_cycle-1 {
				fmt.Printf("Continue instruction addx %d (cycle %d)\n", addToX, i)
				insideAddxForNCycles++
				i++
				continue
			} else {
				X += addToX
				insideAddxForNCycles = -1
				fmt.Printf("Execute instruction addx %d (cycle %d). New X is %d\n", addToX, i, X)
				i++
				continue
			}
		}
		if !scanner.Scan() {
			break
		}
		instruction := scanner.Text()
		if strings.Contains(instruction, addx) {
			insideAddxForNCycles = 1
			addToX, _ = strconv.Atoi(instruction[len(addx):])
			fmt.Printf("Instruction addx %d (cycle %d)\n", addToX, i)
			i++
			continue
		} else if strings.Contains(instruction, noop) {
			fmt.Printf("Instruction noop (cycle %d)\n", i)
			i++
			continue
		}
	}
	return signalStrengths
}

func main() {
	inputArgPtr := flag.String("i", "./10_inputs.txt", "The input data (.txt file).")
	boolArgPtr := flag.Bool("solvePart1", false, "Whether to get the solution for Part 1 or for Part 2.")
	// Parse command line into the defined flags
	flag.Parse()

	sum := solveProblem(*inputArgPtr, *boolArgPtr)

	fmt.Printf("The sum of the signal strengths is %d.\n", sum)
}
