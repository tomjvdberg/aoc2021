package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Grid struct {
	rows map[int]map[int]int
}

type Point struct {
	row    int
	column int
	height int
}

func main() {
	start := time.Now()
	file, err := os.Open("input")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	grid := Grid{map[int]map[int]int{}}
	scanner := bufio.NewScanner(file)

	i := 0
	for scanner.Scan() {
		grid.rows[i] = parseInputLineToRow(scanner.Text())
		i++
	}

	var lowPoints []Point

	for rowIndex, row := range grid.rows {
		for columnIndex, height := range row {
			if height == 9 {
				continue // can never have neighbours that have a higher value as 9 is the highest number
			}

			leftNeighbourHeight, leftNeighbourExists := row[columnIndex-1]
			rightNeighbourHeight, rightNeighbourExists := row[columnIndex+1]
			topNeighbourHeight, topNeighbourExists := grid.rows[rowIndex-1][columnIndex]
			bottomNeighbourHeight, bottomNeighbourExists := grid.rows[rowIndex+1][columnIndex]

			leftHigher := true // default true to be able handling when no neighbour is available
			if leftNeighbourExists && height > leftNeighbourHeight {
				leftHigher = false
			}
			rightHigher := true // default true to be able handling when no neighbour is available
			if rightNeighbourExists && height > rightNeighbourHeight {
				rightHigher = false
			}
			topHigher := true // default true to be able handling when no neighbour is available
			if topNeighbourExists && height > topNeighbourHeight {
				topHigher = false
			}
			bottomHigher := true // default true to be able handling when no neighbour is available
			if bottomNeighbourExists && height > bottomNeighbourHeight {
				bottomHigher = false
			}

			if leftHigher && rightHigher && topHigher && bottomHigher {
				lowPoints = append(lowPoints, Point{rowIndex, columnIndex, height})
			}
		}
	}

	sumOfRiskPoints := 0
	for _, point := range lowPoints {
		sumOfRiskPoints += point.height + 1
	}
	fmt.Println("count of lowPoints", len(lowPoints))
	fmt.Println("answer", sumOfRiskPoints)
	fmt.Println("End", time.Since(start))
}

func parseInputLineToRow(input string) map[int]int {
	heights := map[int]int{}
	heightsAsString := strings.Split(input, "")

	for i, heightAsString := range heightsAsString {
		intVal, _ := strconv.Atoi(heightAsString)
		heights[i] = intVal
	}

	return heights
}
