package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

const RockPoints int = 1

var RockSigns = [2]string{"A", "X"}

const PaperPoints int = 2

var PaperSigns = [2]string{"B", "Y"}

const ScissorsPoints int = 3

var ScissorsSigns = [2]string{"C", "Z"}

const LossPoints int = 0
const Loss string = "X"
const DrawPoints int = 3
const Draw string = "Y"
const WinPoints int = 6
const Win string = "Z"

func solveProblem(filePath string, solvePart1 bool) int {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	totalScore := 0
	if solvePart1 {
		totalScore = getTotalScorePart1(csvReader)
	} else {
		totalScore = getTotalScorePart2(csvReader)
	}

	return totalScore
}

func getTotalScorePart1(data *csv.Reader) int {
	// initialization
	totalScore := 0

	for {
		rec, err := data.Read()
		// deals with errors
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		curr_score := 0
		// switch data[0]: the opponent sign
		switch rec[0] {
		case RockSigns[0]:
			switch rec[1] {
			case RockSigns[1]:
				curr_score += DrawPoints + RockPoints
				//fmt.Printf("This is a draw (%d points) of rocks (%d points). Round score is %d points.\n", DrawPoints, RockPoints, curr_score)
			case PaperSigns[1]:
				curr_score += WinPoints + PaperPoints
				//fmt.Printf("This is a win (%d points) of paper (%d points). Round score is %d points.\n", WinPoints, PaperPoints, curr_score)
			case ScissorsSigns[1]:
				curr_score += LossPoints + ScissorsPoints
				//fmt.Printf("This is a loss (%d points) of scissors (%d points). Round score is %d points.\n", LossPoints, ScissorsPoints, curr_score)
			}
		case PaperSigns[0]:
			switch rec[1] {
			case RockSigns[1]:
				curr_score += LossPoints + RockPoints
				//fmt.Printf("This is a loss (%d points) of rock (%d points). Round score is %d points.\n", LossPoints, RockPoints, curr_score)
			case PaperSigns[1]:
				curr_score += DrawPoints + PaperPoints
				//fmt.Printf("This is a draw (%d points) of papers (%d points). Round score is %d points.\n", DrawPoints, PaperPoints, curr_score)
			case ScissorsSigns[1]:
				curr_score += WinPoints + ScissorsPoints
				//fmt.Printf("This is a win (%d points) of scissors (%d points). Round score is %d points.\n", WinPoints, ScissorsPoints, curr_score)
			}
		case ScissorsSigns[0]:
			switch rec[1] {
			case RockSigns[1]:
				curr_score += WinPoints + RockPoints
				//fmt.Printf("This is a win (%d points) of rock (%d points). Round score is %d points.\n", WinPoints, RockPoints, curr_score)
			case PaperSigns[1]:
				curr_score += LossPoints + PaperPoints
				//fmt.Printf("This is a loss (%d points) of paper (%d points). Round score is %d points.\n", LossPoints, PaperPoints, curr_score)
			case ScissorsSigns[1]:
				curr_score += DrawPoints + ScissorsPoints
				//fmt.Printf("This is a draw (%d points) of scissors (%d points). Round score is %d points.\n", DrawPoints, ScissorsPoints, curr_score)
			}
		}
		totalScore += curr_score
	}
	//fmt.Println()
	return totalScore
}

func getTotalScorePart2(data *csv.Reader) int {
	// initialization
	totalScore := 0

	for {
		rec, err := data.Read()
		// deals with errors
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		curr_score := 0
		// switch data[0]: the opponent sign
		switch rec[0] {
		case RockSigns[0]:
			switch rec[1] {
			case Loss:
				curr_score += LossPoints + ScissorsPoints
				//fmt.Printf("This is a loss (%d points) of scissors (%d points). Round score is %d points.\n", LossPoints, ScissorsPoints, curr_score)
			case Draw:
				curr_score += DrawPoints + RockPoints
				//fmt.Printf("This is a draw (%d points) of rocks (%d points). Round score is %d points.\n", DrawPoints, RockPoints, curr_score)
			case Win:
				curr_score += WinPoints + PaperPoints
				//fmt.Printf("This is a win (%d points) of paper (%d points). Round score is %d points.\n", WinPoints, PaperPoints, curr_score)
			}
		case PaperSigns[0]:
			switch rec[1] {
			case Loss:
				curr_score += LossPoints + RockPoints
				//fmt.Printf("This is a loss (%d points) of rock (%d points). Round score is %d points.\n", LossPoints, RockPoints, curr_score)
			case Draw:
				curr_score += DrawPoints + PaperPoints
				//fmt.Printf("This is a draw (%d points) of papers (%d points). Round score is %d points.\n", DrawPoints, PaperPoints, curr_score)
			case Win:
				curr_score += WinPoints + ScissorsPoints
				//fmt.Printf("This is a win (%d points) of scissors (%d points). Round score is %d points.\n", WinPoints, ScissorsPoints, curr_score)
			}
		case ScissorsSigns[0]:
			switch rec[1] {
			case Loss:
				curr_score += LossPoints + PaperPoints
				//fmt.Printf("This is a loss (%d points) of paper (%d points). Round score is %d points.\n", LossPoints, PaperPoints, curr_score)
			case Draw:
				curr_score += DrawPoints + ScissorsPoints
				//fmt.Printf("This is a draw (%d points) of scissors (%d points). Round score is %d points.\n", DrawPoints, ScissorsPoints, curr_score)
			case Win:
				curr_score += WinPoints + RockPoints
				//fmt.Printf("This is a win (%d points) of rock (%d points). Round score is %d points.\n", WinPoints, RockPoints, curr_score)
			}
		}
		totalScore += curr_score
	}
	//fmt.Println()
	return totalScore
}

func main() {
	inputArgPtr := flag.String("i", "./02_inputs.csv", "The input data (.csv file).")
	boolArgPtr := flag.Bool("solvePart1", false, "Whether to get the solution for Part 1 or for Part 2.")
	// Parse command line into the defined flags
	flag.Parse()

	totalScore := solveProblem(*inputArgPtr, *boolArgPtr)
	fmt.Printf("Your total score for this game is %d points.\n", totalScore)
}
