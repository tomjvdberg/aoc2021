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
	grid := gridFromFile()
	pa := path{}
	var pathsToDiagonalPos = map[coord]path{
		createCoord(1, 0): *pa.withAddedStep(createCoord(1, 0), &grid),
		createCoord(0, 1): *pa.withAddedStep(createCoord(0, 1), &grid),
	}

	for {
		nextPathsToDiagonalPos := map[coord]path{}
		for c, p := range pathsToDiagonalPos {
			// find two neighbours
			rightNeighbour := coord{c.x + 1, c.y}
			bottomNeighbour := coord{c.x, c.y + 1}

			for _, neighbour := range []coord{rightNeighbour, bottomNeighbour} {
				if coordExistsInGrid(neighbour, grid) {
					newPath := *p.withAddedStep(neighbour, &grid)
					// add the new path if it is cheaper then an existing one
					existingPath, exists := nextPathsToDiagonalPos[neighbour]
					if !exists || newPath.totalRisk < existingPath.totalRisk {
						nextPathsToDiagonalPos[neighbour] = newPath
						continue
					}
					if newPath.totalRisk < existingPath.totalRisk {
						nextPathsToDiagonalPos[neighbour] = newPath
					}
				}
			}
		}
		pathsToDiagonalPos = nextPathsToDiagonalPos
		if len(pathsToDiagonalPos) == 1 {
			break // we reached the end
		}
	}

	fmt.Println("Least Risk", pathsToDiagonalPos[createCoord(len(grid[0])-1, len(grid)-1)].totalRisk)
	fmt.Println("End", time.Since(start))
}

func coordExistsInGrid(neighbour coord, grid [][]int) bool {
	return len(grid)-1 >= neighbour.y && len(grid[neighbour.y])-1 >= neighbour.x
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
