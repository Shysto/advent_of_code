package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const noteLength int = 7
const topActiveMonkeys int = 2
const add string = "+"
const mult string = "*"

var divideWorryLevelBy int
var nbRounds int
var limit int = 1

type Queue []int

func (q *Queue) isEmpty() bool {
	return len(*q) == 0
}

func (q *Queue) push(element int) {
	*q = append(*q, element) // Simply append to q.
}

func (q *Queue) pop() (int, bool) {
	if q.isEmpty() {
		return -1, false
	}
	element := (*q)[0] // The first element is the one to be dequeued.
	// Slice off the element once it is dequeued.
	if len(*q) == 1 {
		*q = Queue{}
	} else {
		*q = (*q)[1:]
	}
	return element, true
}

type Monkey struct {
	Id               int
	Items            Queue
	Operation        func(*Queue, int)
	DivisibilityTest func(int) bool
	TrueMonkeyId     int
	FalseMonkeyId    int
}

var monkeys []Monkey

func getMonkeyBusiness(nbInspections []int, getTopK int) (monkeyBusiness int) {
	monkeyBusiness = 1
	sort.Slice(nbInspections, func(i, j int) bool {
		return nbInspections[i] > nbInspections[j]
	})
	for _, nb := range nbInspections[:getTopK] {
		monkeyBusiness *= nb
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
	parseNotes(scanner, solvePart1)
	//fmt.Println(grid)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	nbInspections := runRounds()

	return getMonkeyBusiness(nbInspections, topActiveMonkeys)
}

func parseNotes(scanner *bufio.Scanner, solvePart1 bool) {
	i := 1
	var monkey Monkey
	for scanner.Scan() {
		note := scanner.Text()
		// Create new monkey
		if i%noteLength == 1 {
			monkey = Monkey{Id: i / noteLength}
			// Starting indexes
		} else if i%noteLength == 2 {
			splits := strings.Split(note, ":")
			//display := "Monkey " + strconv.Itoa(monkey.Id) + ":\n  Starting items: "
			for _, elem := range strings.Split(splits[1], ",") {
				worryLevel, _ := strconv.Atoi(strings.TrimSpace(elem))
				//display += strings.TrimSpace(elem) + ", "
				monkey.Items.push(worryLevel)
			}
			//fmt.Println(display)
			// Operation
		} else if i%noteLength == 3 {
			//fmt.Println(note)
			if strings.Contains(note, add) {
				factor, _ := strconv.Atoi(strings.TrimSpace(strings.Split(note, add)[1]))
				monkey.Operation = func(q *Queue, l int) {
					for i := 0; i < len(*q); i++ {
						(*q)[i] += factor
						if l > 1 {
							(*q)[i] %= l
						}
					}
				}
			} else if strings.Contains(note, mult) {
				factor_str := strings.TrimSpace(strings.Split(note, mult)[1])
				if factor_str == "old" {
					monkey.Operation = func(q *Queue, l int) {
						for i := 0; i < len(*q); i++ {
							(*q)[i] *= (*q)[i]
							if l > 1 {
								(*q)[i] %= l
							}
						}
					}
				} else {
					factor, _ := strconv.Atoi(factor_str)

					monkey.Operation = func(q *Queue, l int) {
						for i := 0; i < len(*q); i++ {
							(*q)[i] *= factor
							if l > 1 {
								(*q)[i] %= l
							}
						}
					}
				}
			} else {
				log.Fatalf("Unknown operation: %v", note)
			}
			//monkey.Operation(&monkey.Items)
			//fmt.Printf("  Test operation: %v\n", monkey.Items)
			// Divisibility test
		} else if i%noteLength == 4 {
			//fmt.Println(note)
			divisibilityFactor, _ := strconv.Atoi(strings.TrimSpace(strings.Split(note, "by")[1]))
			if !solvePart1 {
				limit *= divisibilityFactor
			}

			monkey.DivisibilityTest = func(worryLevel int) bool {
				return worryLevel%divisibilityFactor == 0
			}
			//fmt.Printf("  Test divisibility: 7 is %v, 2 is %v, 19 is %v\n", monkey.DivisibilityTest(7), monkey.DivisibilityTest(2), monkey.DivisibilityTest(19))
			// Send to monkey if test divisibility is true
		} else if i%noteLength == 5 {
			trueMonkey, _ := strconv.Atoi(strings.TrimSpace(strings.Split(note, "monkey")[1]))
			monkey.TrueMonkeyId = trueMonkey
			//fmt.Printf("   If true: throw to monkey %d\n", monkey.TrueMonkeyId)
			// Send to monkey if test divisibility is false
		} else if i%noteLength == 6 {
			falseMonkey, _ := strconv.Atoi(strings.TrimSpace(strings.Split(note, "monkey")[1]))
			monkey.FalseMonkeyId = falseMonkey
			//fmt.Printf("   If false: throw to monkey %d\n", monkey.FalseMonkeyId)

			// add monkey
			monkeys = append(monkeys, monkey)
		}
		i++
	}
}

func runRounds() []int {
	nbInspections := make([]int, len(monkeys))
	for i := 1; i <= nbRounds; i++ {
		//fmt.Println()
		for m := 0; m < len(monkeys); m++ {
			//fmt.Printf("Monkey %d:\n", monkeys[m].Id)
			//cpy := make([]int, len(monkeys[m].Items))
			//copy(cpy, monkeys[m].Items)
			monkeys[m].Operation(&monkeys[m].Items, limit)
			k := 0
			for !monkeys[m].Items.isEmpty() {
				nbInspections[m]++
				worryLevel, err := monkeys[m].Items.pop()
				if !err {
					continue
				}
				//fmt.Printf("  Monkey inspects an item with a worry level of %d.\n", cpy[k])
				//fmt.Printf("    Worry level is changed to %d.\n", worryLevel)
				worryLevel /= divideWorryLevelBy
				//fmt.Printf("    Monkey gets bored with item. Worry level is divided by %d to %d.\n", divideWorryLevelBy, worryLevel)
				if monkeys[m].DivisibilityTest(worryLevel) {
					//fmt.Println("    Current worry level is divisible.")
					monkeys[monkeys[m].TrueMonkeyId].Items.push(worryLevel)
					//fmt.Printf("    Item with worry level %d is thrown to monkey %d.\n", worryLevel, monkeys[m].TrueMonkeyId)
				} else {
					//fmt.Println("    Current worry level is not divisible.")
					monkeys[monkeys[m].FalseMonkeyId].Items.push(worryLevel)
					//fmt.Printf("    Item with worry level %d is thrown to monkey %d.\n", worryLevel, monkeys[m].FalseMonkeyId)
				}
				k++
			}
		}
		//fmt.Printf("After round %d, the monkeys are holding items with these worry levels:\n", i)
		//for _, monkey := range monkeys {
		//fmt.Printf("Monkey %d: %v\n", monkey.Id, monkey.Items)
		//}
	}
	return nbInspections
}

func main() {
	inputArgPtr := flag.String("i", "./11_inputs.txt", "The input data (.txt file).")
	boolArgPtr := flag.Bool("solvePart1", false, "Whether to get the solution for Part 1 or for Part 2.")
	// Parse command line into the defined flags
	flag.Parse()

	if *boolArgPtr {
		divideWorryLevelBy = 3
		nbRounds = 20
	} else {
		divideWorryLevelBy = 1
		nbRounds = 10000
	}

	monkeyBusiness := solveProblem(*inputArgPtr, *boolArgPtr)

	fmt.Printf("The level of monkey business after %d rounds is %d.\n", nbRounds, monkeyBusiness)
}
