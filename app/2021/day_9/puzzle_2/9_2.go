package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Grid struct {
	rows   map[int]map[int]int
	points map[string]int
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

	grid := Grid{
		map[int]map[int]int{},
		map[string]int{},
	}
	scanner := bufio.NewScanner(file)

	row := 0
	for scanner.Scan() {
		grid.rows[row] = parseInputLineToRow(scanner.Text())
		row++
	}

	lowPoints := map[string]Point{}

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
				point := Point{rowIndex, columnIndex, height}
				lowPoints[pointIndex(point.row, point.column)] = point
			}
		}
	}

	// each low point is in ONE basin
	// all points belong to ONE basin
	// 9's are the delimiters of the basins
	var basins []int // basins by size
	for pointId, lowPoint := range lowPoints {
		unvisitedPoints := map[string]Point{}
		visitedPoints := map[string]Point{}
		// Always add the lowpoint to visitedPoints
		visitedPoints[pointId] = lowPoint
		// scan horizontally until nine or border detected
		visiblePoints := getNeighboursOf(lowPoint, grid)
		// add points to unvisitedPoints unless they are visited
		pointsDiff(&visiblePoints, &visitedPoints, &unvisitedPoints)
		visiteRecursive(&visiblePoints, &visitedPoints, &unvisitedPoints, grid)

		basins = append(basins, len(visitedPoints))
	}
	sort.Ints(basins)

	// find the three largest basins = sort by size. Pick last three.
	sum := 0
	largestBasins := basins[len(basins)-3:]
	sum = largestBasins[0] * largestBasins[1] * largestBasins[2]
	fmt.Println("sum", sum)
	fmt.Println("End", time.Since(start))
}

func visiteRecursive(
	visiblePoints *map[string]Point,
	visitedPoints *map[string]Point,
	unvisitedPoints *map[string]Point,
	grid Grid,
) {
	if len(*unvisitedPoints) == 0 {
		return
	}

	workingUnvisitedPoints := map[string]Point{}
	for id, unvisitedPoint := range *unvisitedPoints {
		(*visitedPoints)[id] = unvisitedPoint // I now visit this point
		visiblePointsT := getNeighboursOf(unvisitedPoint, grid)
		// add points to unvisitedPoints unless they are visited
		pointsDiff(&visiblePointsT, visitedPoints, &workingUnvisitedPoints)
	}
	visiteRecursive(visiblePoints, visitedPoints, &workingUnvisitedPoints, grid)
}

func pointsDiff(
	visiblePoints *map[string]Point,
	visitedPoints *map[string]Point,
	unvisitedPoints *map[string]Point,
) {
	for id, visiblePoint := range *visiblePoints {
		_, isVisited := (*visitedPoints)[id]
		if !isVisited {
			(*unvisitedPoints)[id] = visiblePoint
		}
	}
}

func getNeighboursOf(point Point, grid Grid) map[string]Point {
	visiblePoints := map[string]Point{}
	// watch left
	columnToVisit := point.column - 1
	nextPointHeight, exists := grid.rows[point.row][columnToVisit]
	if exists && nextPointHeight != 9 {
		visiblePoints[pointIndex(point.row, columnToVisit)] = Point{point.row, columnToVisit, nextPointHeight}
	}
	// watch right
	columnToVisit = point.column + 1
	nextPointHeight, exists = grid.rows[point.row][columnToVisit]
	if exists && nextPointHeight != 9 {
		visiblePoints[pointIndex(point.row, columnToVisit)] = Point{point.row, columnToVisit, nextPointHeight}
	}
	// Move up
	rowToVisit := point.row - 1
	nextPointHeight, exists = grid.rows[rowToVisit][point.column]
	if exists && nextPointHeight != 9 {
		visiblePoints[pointIndex(rowToVisit, point.column)] = Point{rowToVisit, point.column, nextPointHeight}
	}
	// Move down
	rowToVisit = point.row + 1
	nextPointHeight, exists = grid.rows[rowToVisit][point.column]
	if exists && nextPointHeight != 9 {
		visiblePoints[pointIndex(rowToVisit, point.column)] = Point{rowToVisit, point.column, nextPointHeight}
	}

	return visiblePoints
}

func getVisiblePointsFrom(point Point, grid Grid) map[string]Point {
	visiblePoints := map[string]Point{}
	// move left
	for i := 1; ; i++ {
		columnToVisit := point.column - i
		nextPointHeight, exists := grid.rows[point.row][columnToVisit]
		if !exists || nextPointHeight == 9 {
			break
		}
		visiblePoints[pointIndex(point.row, columnToVisit)] = Point{point.row, columnToVisit, nextPointHeight}
	}
	// Move right
	for i := 1; ; i++ {
		columnToVisit := point.column + i
		nextPointHeight, exists := grid.rows[point.row][columnToVisit]
		if !exists || nextPointHeight == 9 {
			break
		}
		visiblePoints[pointIndex(point.row, columnToVisit)] = Point{point.row, columnToVisit, nextPointHeight}
	}
	// Move up
	for i := 1; ; i++ {
		rowToVisit := point.row - i
		nextPointHeight, exists := grid.rows[rowToVisit][point.column]
		if !exists || nextPointHeight == 9 {
			break
		}
		visiblePoints[pointIndex(rowToVisit, point.column)] = Point{rowToVisit, point.column, nextPointHeight}
	}
	// Move down
	for i := 1; ; i++ {
		rowToVisit := point.row + i
		nextPointHeight, exists := grid.rows[rowToVisit][point.column]
		if !exists || nextPointHeight == 9 {
			break
		}
		visiblePoints[pointIndex(rowToVisit, point.column)] = Point{rowToVisit, point.column, nextPointHeight}
	}

	return visiblePoints
}

func pointIndex(row int, column int) string {
	return strconv.Itoa(row) + "," + strconv.Itoa(column)
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
