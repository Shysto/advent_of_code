package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"strings"
)

const leftBracket string = "["
const rightBracket string = "]"
const sep string = ","

var dividerPackets = []string{"[[2]]", "[[6]]"}

type Node struct {
	Value    int
	Children []*Node
	Parent   *Node
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solveProblem(filePath string, solvePart1 bool) int {
	content, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}
	var answer int
	if solvePart1 {
		answer = countOrderedPackets(string(content))
	} else {
		answer = sortPackets(string(content))
	}

	return answer
}

func parsePackets(packet_str string) Node {
	packet := Node{-1, []*Node{}, nil}
	tmp := &packet

	var curr_int string
	for _, char := range packet_str {
		switch string(char) {
		case leftBracket: // We create a new list
			newList := Node{-1, []*Node{}, tmp}
			tmp.Children = append(tmp.Children, &newList)
			tmp = &newList
		case rightBracket: // We finish the current list, so we go up to the parent node
			if len(curr_int) > 0 { // We append the current value if needed
				n, _ := strconv.Atoi(curr_int)
				tmp.Value = n
				curr_int = ""
			}
			tmp = tmp.Parent // We go up to the parent node
		case sep:
			if len(curr_int) > 0 { // We append the current value if needed
				n, _ := strconv.Atoi(curr_int)
				tmp.Value = n
				curr_int = ""
			}
			tmp = tmp.Parent                    // We go up to the parent node
			newList := Node{-1, []*Node{}, tmp} // We create a new list
			tmp.Children = append(tmp.Children, &newList)
			tmp = &newList
		default: // Multiple digit number
			curr_int += string(char)
		}
	}
	return packet
}

func areRightOrdered(left, right Node) int {
	switch {
	case len(left.Children) == 0 && len(right.Children) == 0: // Integers comparison
		if left.Value > right.Value {
			return -1
		} else if left.Value == right.Value {
			return 0
		}
		return 1

	case left.Value >= 0: // Integer with list comparison
		return areRightOrdered(Node{-1, []*Node{&left}, nil}, right)

	case right.Value >= 0: // List with integer comparison
		return areRightOrdered(left, Node{-1, []*Node{&right}, nil})

	default: // List comparison
		for i := 0; i < min(len(left.Children), len(right.Children)); i++ { // List comparison element by element
			areOrdered := areRightOrdered(*left.Children[i], *right.Children[i])
			if areOrdered != 0 {
				return areOrdered
			}
		}
		// List size comparison
		if len(left.Children) < len(right.Children) {
			return 1
		} else if len(left.Children) > len(right.Children) {
			return -1
		}
	}
	return 0
}

func countOrderedPackets(raw_data string) int {
	packets := strings.Split(raw_data, "\n")
	k, i := 1, 0
	sum := 0
	for i < len(packets) {
		first_packet := packets[i]
		second_packet := packets[i+1]
		//fmt.Printf("Packet %d: %v (index %d)\n", k, first_packet, i)
		//fmt.Printf("Packet %d: %v (index %d)\n", k+1, second_packet, i)

		left := parsePackets(first_packet)
		right := parsePackets(second_packet)

		if areRightOrdered(left, right) == 1 {
			sum += k
		}

		i += 3
		k++
	}
	return sum
}

func sortPackets(raw_data string) int {
	lines := strings.Split(raw_data, "\n")
	var packets []Node
	// List all packets
	for i := 0; i < len(lines); i += 3 {
		packet := lines[i]
		packets = append(packets, parsePackets(packet))
		packet = lines[i+1]
		packets = append(packets, parsePackets(packet))
	}
	// Add divider packets
	for i := 0; i < len(dividerPackets); i++ {
		packets = append(packets, parsePackets(dividerPackets[i]))
	}

	// Sort list of packets (in place)
	sort.Slice(packets, func(i, j int) bool {
		return areRightOrdered(packets[i], packets[j]) == 1
	})

	// Find indices of divider packets
	decoder_key := 1
	for _, divPackStr := range dividerPackets {
		divPack := parsePackets(divPackStr)
		for i, packet := range packets {
			if areRightOrdered(packet, divPack) == 0 {
				decoder_key *= i + 1
				break
			}
		}
	}
	return decoder_key
}

func main() {
	inputArgPtr := flag.String("i", "./13_inputs.txt", "The input data (.txt file).")
	boolArgPtr := flag.Bool("solvePart1", false, "Whether to get the solution for Part 1 or for Part 2.")
	// Parse command line into the defined flags
	flag.Parse()

	answer := solveProblem(*inputArgPtr, *boolArgPtr)

	if *boolArgPtr {
		fmt.Printf("The sum of the indices of pairs of packets already in the right order is %d.\n", answer)
	} else {
		fmt.Printf("The decoder key for the distress signal is %d.\n", answer)
	}
}
