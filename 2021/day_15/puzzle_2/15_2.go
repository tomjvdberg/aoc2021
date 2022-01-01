package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type coord struct {
	x int
	y int
}
type path struct {
	totalRisk int
	coords    []coord
}

func (p *path) withAddedStep(c coord, g *[][]int) *path {
	newP := path{}
	for _, coordInOldPath := range p.coords {
		newC := coord{}
		newC.x = coordInOldPath.x
		newC.y = coordInOldPath.y

		newP.coords = append(newP.coords, newC)
		newP.totalRisk += (*g)[newC.y][newC.x]
	}
	newP.coords = append(newP.coords, c)
	newP.totalRisk += (*g)[c.y][c.x]

	return &newP
}

func main() {
	start := time.Now()
	grid := applyPuzzle2RulesToGrid(gridFromFile())

	pa := path{}
	lowestRiskPathsCoords := map[coord]path{
		createCoord(1, 0): *pa.withAddedStep(createCoord(1, 0), &grid),
		createCoord(0, 1): *pa.withAddedStep(createCoord(0, 1), &grid),
	}
	coordsThatNeedToBeVisited := map[coord]path{
		createCoord(1, 0): *pa.withAddedStep(createCoord(1, 0), &grid),
		createCoord(0, 1): *pa.withAddedStep(createCoord(0, 1), &grid),
	}

	for {
		nextCoordsThatNeedToBeVisited := map[coord]path{}
		for c, p := range coordsThatNeedToBeVisited {
			// find two neighbours
			leftNeighbour := coord{c.x - 1, c.y}
			topNeighbour := coord{c.x, c.y - 1}
			rightNeighbour := coord{c.x + 1, c.y}
			bottomNeighbour := coord{c.x, c.y + 1}

			for _, neighbour := range []coord{leftNeighbour, topNeighbour, bottomNeighbour, rightNeighbour} {
				if !coordExistsInGrid(neighbour, grid) {
					continue
				}

				newPath := *p.withAddedStep(neighbour, &grid)
				// add the new path if it is cheaper then an existing one
				existingPath, exists := lowestRiskPathsCoords[neighbour]
				if !exists || newPath.totalRisk < existingPath.totalRisk {
					lowestRiskPathsCoords[neighbour] = newPath
					nextCoordsThatNeedToBeVisited[neighbour] = newPath
					continue
				}
			}
		}
		if len(nextCoordsThatNeedToBeVisited) == 0 {
			break // we reached the end
		}
		coordsThatNeedToBeVisited = nextCoordsThatNeedToBeVisited
	}

	fmt.Println("Least Risk", lowestRiskPathsCoords[createCoord(len(grid[0])-1, len(grid)-1)].totalRisk)
	fmt.Println("End", time.Since(start))
}

func coordExistsInGrid(neighbour coord, grid [][]int) bool {
	return neighbour.y > -1 && neighbour.x > -1 && len(grid)-1 >= neighbour.y && len(grid[neighbour.y])-1 >= neighbour.x
}

func createCoord(x int, y int) coord {
	return coord{
		x,
		y,
	}
}

func gridFromFile() [][]int {
	file, err := os.Open("input")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	var grid [][]int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var row []int
		positionsAsString := strings.Split(scanner.Text(), "")
		for _, positionAsString := range positionsAsString {
			posInt, _ := strconv.Atoi(positionAsString)
			row = append(row, posInt)
		}
		grid = append(grid, row)
	}

	return grid
}

func applyPuzzle2RulesToGrid(grid [][]int) [][]int {
	// first copy to right
	newGridRight := make([][]int, len(grid))
	copy(newGridRight, grid)

	for i := 1; i < 5; i++ {
		for rowI, row := range grid {
			for _, risk := range row {
				newRisk := risk + i
				if newRisk > 9 {
					newRisk = newRisk - 9
				}

				newGridRight[rowI] = append(newGridRight[rowI], newRisk)
			}
		}
	}

	// then to bottom
	newGrid := make([][]int, len(newGridRight))
	copy(newGrid, newGridRight)
	for i := 1; i < 5; i++ {
		for _, row := range newGridRight {
			newRow := []int{}

			for _, risk := range row {
				newRisk := risk + i
				if newRisk > 9 {
					newRisk = newRisk - 9
				}

				newRow = append(newRow, newRisk)
			}

			newGrid = append(newGrid, newRow)
		}
	}

	return newGrid
}
