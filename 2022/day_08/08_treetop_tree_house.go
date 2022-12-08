package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Grid [][]int

func (g Grid) isEmpty() bool {
	return len(g) < 1
}
func (g Grid) isBorder(i, j int) bool {
	return i == 0 || j == 0 || i == len(g)-1 || j == len(g[0])-1
}

func (g Grid) getTopViewingDistance(i, j int) int {
	height := g[i][j]
	viewingDistance := 0
	for k := i - 1; k >= 0; k-- {
		viewingDistance++
		if g[k][j] >= height {
			return viewingDistance
		}
	}
	return viewingDistance
}

func (g Grid) isVisibleFromTop(i, j int) bool {
	if g.isBorder(i, j) {
		return true
	}
	return (g[0][j] < g[i][j]) && (g.getTopViewingDistance(i, j) == i)
}

func (g Grid) getBottomViewingDistance(i, j int) int {
	height := g[i][j]
	viewingDistance := 0
	for k := i + 1; k < len(g); k++ {
		viewingDistance++
		if g[k][j] >= height {
			return viewingDistance
		}
	}
	return viewingDistance
}

func (g Grid) isVisibleFromBottom(i, j int) bool {
	if g.isBorder(i, j) {
		return true
	}
	return (g[len(g)-1][j] < g[i][j]) && (g.getBottomViewingDistance(i, j) == len(g)-i-1)
}

func (g Grid) getLeftViewingDistance(i, j int) int {
	height := g[i][j]
	viewingDistance := 0
	for k := j - 1; k >= 0; k-- {
		viewingDistance++
		if g[i][k] >= height {
			return viewingDistance
		}
	}
	return viewingDistance
}

func (g Grid) isVisibleFromLeft(i, j int) bool {
	if g.isBorder(i, j) {
		return true
	}
	return (g[i][0] < g[i][j]) && (g.getLeftViewingDistance(i, j) == j)
}

func (g Grid) getRightViewingDistance(i, j int) int {
	height := g[i][j]
	viewingDistance := 0
	for k := j + 1; k < len(g[0]); k++ {
		viewingDistance++
		if g[i][k] >= height {
			return viewingDistance
		}
	}
	return viewingDistance
}

func (g Grid) isVisibleFromRight(i, j int) bool {
	if g.isBorder(i, j) {
		return true
	}
	return (g[i][len(g[0])-1] < g[i][j]) && (g.getRightViewingDistance(i, j) == len(g[0])-j-1)
}

func (g Grid) isVisible(i, j int) bool {
	return g.isBorder(i, j) || g.isVisibleFromLeft(i, j) || g.isVisibleFromTop(i, j) || g.isVisibleFromRight(i, j) || g.isVisibleFromBottom(i, j)
}

func (g Grid) countVisible() int {
	if g.isEmpty() {
		return 0
	}
	nbVisibleTrees := 0
	for i := 0; i < len(g); i++ {
		for j := 0; j < len(g[0]); j++ {
			if g.isVisible(i, j) {
				nbVisibleTrees++
			}
		}
	}
	return nbVisibleTrees
}

func (g Grid) computeScenicScore(i, j int) int {
	return g.getLeftViewingDistance(i, j) * g.getTopViewingDistance(i, j) * g.getRightViewingDistance(i, j) * g.getBottomViewingDistance(i, j)
}

func (g Grid) getBestScenicScore() int {
	if g.isEmpty() {
		return 0
	}

	bestScenicScore := 0
	for i := 0; i < len(g); i++ {
		for j := 0; j < len(g[0]); j++ {
			currScenicScore := g.computeScenicScore(i, j)
			if currScenicScore > bestScenicScore {
				bestScenicScore = currScenicScore
			}
		}
	}
	return bestScenicScore
}

func solveProblem(filePath string, solvePart1 bool) int {
	f, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	grid := makeGrid(scanner)
	//fmt.Println(grid)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if solvePart1 {
		return grid.countVisible()
	} else {
		return grid.getBestScenicScore()
	}
}

func makeGrid(scanner *bufio.Scanner) Grid {
	grid := Grid{}
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]int, len(line))
		for j, tree := range line {
			height, err := strconv.Atoi(string(tree))
			if err != nil {
				log.Fatal(err)
			}
			row[j] = height
		}
		grid = append(grid, row)
	}
	return grid
}

func main() {
	inputArgPtr := flag.String("i", "./08_inputs.txt", "The input data (.txt file).")
	boolArgPtr := flag.Bool("solvePart1", false, "Whether to get the solution for Part 1 or for Part 2.")
	// Parse command line into the defined flags
	flag.Parse()

	number := solveProblem(*inputArgPtr, *boolArgPtr)
	if *boolArgPtr {
		fmt.Printf("There are %d trees visible from outside the grid.\n", number)
	} else {
		fmt.Printf("The highest scenic score is %d.\n", number)
	}

}
