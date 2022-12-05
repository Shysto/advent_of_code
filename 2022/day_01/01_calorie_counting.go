package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
)

type Elf struct {
	Id          int
	MaxCalories int
}

const delimiter string = "."

func solveProblem(filePath string) []Elf {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	elfs := getMaxCalories(csvReader)

	return elfs
}

func getMaxCalories(data *csv.Reader) []Elf {
	// initialization
	elfs := make([]Elf, 0)
	curr_elf_id := 1
	curr_elf_calories := 0

	for {
		rec, err := data.Read()
		// deals with errors
		if err == io.EOF {
			elfs = append(elfs, Elf{curr_elf_id, curr_elf_calories})
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// switch data: '.' means we go to another elf
		if rec[0] == delimiter {
			elfs = append(elfs, Elf{curr_elf_id, curr_elf_calories})
			curr_elf_id++
			curr_elf_calories = 0
		} else {
			calories, err := strconv.Atoi(rec[0])
			if err != nil {
				log.Fatal(err)
			}
			curr_elf_calories += calories
			//fmt.Printf("Current calories: %d.\n", curr_elf_calories)
		}
	}
	//fmt.Println()
	return elfs
}

func sumCalories(elfs []Elf) int {
	result := 0
	for _, v := range elfs {
		result += v.MaxCalories
	}
	return result
}

func main() {
	inputArgPtr := flag.String("i", "./01_inputs.csv", "The input data (.csv file).")
	topKArgPtr := flag.Int("topK", 3, "The top-k Elfs to retrieve.")
	// Parse command line into the defined flags
	flag.Parse()

	elfs := solveProblem(*inputArgPtr)
	for i, v := range elfs {
		fmt.Printf("The %dth Elf carries %d Calories.\n", i, v)
	}
	top_k := *topKArgPtr
	// sorts a slice of Elf in decreasing order of their MaxCalories
	sort.Slice(elfs, func(i, j int) bool {
		return elfs[i].MaxCalories > elfs[j].MaxCalories
	})
	fmt.Printf("The top-%d Elf(s) carrying the most Calories are %+v.\n", top_k, elfs[:top_k])
	fmt.Printf("The top-%d Elf(s) carrying the most Calories have a total of %d Calories.\n", top_k, sumCalories(elfs[:top_k]))
}
