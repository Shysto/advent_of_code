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

const prompt string = "$ "
const cd string = "cd "
const ls string = "ls"
const dir string = "dir "
const root string = "/"
const previous string = ".."
const sep string = " "

const maxSize int = 100000
const totalAvailableSpace int = 70000000
const needSpace int = 30000000

type File struct {
	Name string
	Size int
}

type Node struct {
	Name   string
	Size   int
	Files  []*File
	Parent *Node
	Childs map[string]*Node
}

func solveProblem(filePath string, solvePart1 bool) int {
	f, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	folders := getHierarchy(scanner)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var totalSize int
	for _, folder := range folders {
		size := getSize(*folder)
		folder.Size = size
		if size <= maxSize {
			totalSize += size
		}
	}

	if solvePart1 {
		return totalSize
	}

	unusedSpace := totalAvailableSpace - folders[0].Size
	toBeFreed := needSpace - unusedSpace
	if toBeFreed <= 0 {
		return 0
	}
	//fmt.Printf("unusedSpace = %d, toBeFrees = %d\n", unusedSpace, toBeFreed)

	smallestEnoughSize := folders[0].Size
	for _, folder := range folders {
		size := folder.Size
		//fmt.Printf("dir %s size = %d\n", folder.Name, size)
		if (size > toBeFreed) && (size < smallestEnoughSize) {
			smallestEnoughSize = size
		}
	}
	return smallestEnoughSize
}

func getHierarchy(lines *bufio.Scanner) []*Node {
	folders := []*Node{}
	var currDir *Node
	for lines.Scan() {
		line := lines.Text()
		//fmt.Println(line)

		if strings.Contains(line, prompt) {
			if strings.Contains(line, cd) {
				if strings.Contains(line, previous) {
					currDir = currDir.Parent
				} else if strings.Contains(line, root) {
					currDir = &Node{Name: root, Files: []*File{}, Childs: make(map[string]*Node)}
					folders = append(folders, currDir)
				} else {
					currDir = currDir.Childs[line[len(prompt)+len(cd):]]
				}
			} else if strings.Contains(line, ls) {
				continue
			} else {
				log.Fatalf("Unknown command: '%s'", line)
			}
		} else {
			if strings.Contains(line, dir) {
				currDir.Childs[line[len(dir):]] = &Node{Name: line[len(dir):], Parent: currDir, Files: []*File{}, Childs: make(map[string]*Node)}
				folders = append(folders, currDir.Childs[line[len(dir):]])
			} else {
				split := strings.Split(line, sep)
				size, err := strconv.Atoi(split[0])
				if err != nil {
					log.Fatal(err)
				}
				currDir.Files = append(currDir.Files, &File{split[1], size})
			}
		}
		//fmt.Println(currDir)
	}
	return folders
}

func getSize(root Node) int {
	size := 0
	for _, file := range root.Files {
		size += file.Size
	}
	for _, folder := range root.Childs {
		size += getSize(*folder)
	}
	return size
}

func main() {
	inputArgPtr := flag.String("i", "./07_inputs.txt", "The input data (.csv file).")
	boolArgPtr := flag.Bool("solvePart1", false, "Whether to get the solution for Part 1 or for Part 2.")
	// Parse command line into the defined flags
	flag.Parse()

	size := solveProblem(*inputArgPtr, *boolArgPtr)
	if *boolArgPtr {
		fmt.Printf("The sum of the total sizes of directories with a total size of at most %d is %d\n", maxSize, size)
	} else {
		fmt.Printf("The size of the directory to delete in order to have enough free space to run the update is %d\n", size)
	}
}
