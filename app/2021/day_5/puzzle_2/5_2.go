package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Grid struct {
	points map[string]Point
}

func (grid *Grid) addLine(line Line) {
	var coordinates []Coordinate
	if line.isVertical() || line.isHorizontal() {
		coordinates = line.getStraightCoordinates()
	} else {
		coordinates = line.getDiagonalCoordinates()
	}

	for _, coordinate := range coordinates {
		grid.markPoint(coordinate)
	}
}

func (grid *Grid) markPoint(coordinate Coordinate) {
	pointIndex := createPointIndex(coordinate)
	point, exists := grid.points[pointIndex]
	updatePoint := Point{coordinate.x, coordinate.y, 1}

	if exists {
		updatePoint.timesCovered = point.timesCovered + 1
		grid.points[pointIndex] = updatePoint
		return
	}
	grid.points[pointIndex] = updatePoint
}

func (grid *Grid) getCountOfPointsThatAreCoveredAtLeastTwice() int {
	countOfPointsThatWereCoveredTwiceOrMore := 0

	for _, point := range grid.points {
		if point.timesCovered >= 2 {
			countOfPointsThatWereCoveredTwiceOrMore++
		}
	}

	return countOfPointsThatWereCoveredTwiceOrMore
}

type Point struct {
	x            int
	y            int
	timesCovered int
}

type Coordinate struct {
	x int
	y int
}

type Line struct {
	start Coordinate
	end   Coordinate
}

func (line Line) getStraightCoordinates() []Coordinate {
	var coordinates []Coordinate

	var xValues []int
	startX := line.start.x
	endX := line.end.x

	switch {
	case startX == endX:
		xValues = append(xValues, startX)
	case endX > startX:
		xValues = append(xValues, endX)
		xDiff := endX - startX

		for i := 1; i < xDiff; i++ {
			xValues = append(xValues, endX-i)
		}
		xValues = append(xValues, startX)
	case startX > endX:
		xValues = append(xValues, startX)
		xDiff := startX - endX

		for i := 1; i < xDiff; i++ {
			xValues = append(xValues, startX-i)
		}
		xValues = append(xValues, endX)
	}

	var yValues []int
	startY := line.start.y
	endY := line.end.y

	switch {
	case startY == endY:
		yValues = append(yValues, startY)
	case endY > startY:
		yValues = append(yValues, endY)
		xDiff := endY - startY

		for i := 1; i < xDiff; i++ {
			yValues = append(yValues, endY-i)
		}
		yValues = append(yValues, startY)
	case startY > endY:
		yValues = append(yValues, startY)
		xDiff := startY - endY

		for i := 1; i < xDiff; i++ {
			yValues = append(yValues, startY-i)
		}
		yValues = append(yValues, endY)
	}

	for _, xValue := range xValues {
		for _, yValue := range yValues {
			coordinates = append(coordinates, Coordinate{
				x: xValue,
				y: yValue,
			})
		}
	}

	return coordinates
}

func (line Line) getDiagonalCoordinates() []Coordinate {
	var coordinates []Coordinate
	coordinatesBetweenStartAndEnd := Abs(line.start.x-line.end.x) - 1
	xPositiveDirection := line.start.x < line.end.x
	yPositiveDirection := line.start.y < line.end.y

	// always add the start coordinate
	coordinates = append(coordinates, line.start)
	for i := 0; i < coordinatesBetweenStartAndEnd; i++ {
		newCoordinate := Coordinate{}

		if xPositiveDirection {
			newCoordinate.x = line.start.x + 1 + i
		} else {
			newCoordinate.x = line.start.x - 1 - i
		}

		if yPositiveDirection {
			newCoordinate.y = line.start.y + 1 + i
		} else {
			newCoordinate.y = line.start.y - 1 - i
		}
		coordinates = append(coordinates, newCoordinate)
	}
	// always add the end coordinate
	coordinates = append(coordinates, line.end)

	return coordinates
}

func (line Line) isHorizontal() bool {
	return line.start.y == line.end.y
}

func (line Line) isVertical() bool {
	return line.start.x == line.end.x
}

func (line Line) isDiagonal() bool {
	xDiff := Abs(line.start.x - line.end.x)
	yDiff := Abs(line.start.y - line.end.y)

	return xDiff == yDiff
}

func main() {
	file, err := os.Open("input")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	var lines []Line

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// line to Lines
		lines = append(lines, inputToLine(scanner.Text()))
	}
	grid := Grid{map[string]Point{}}

	for _, line := range lines {
		if !(line.isHorizontal() || line.isVertical() || line.isDiagonal()) {
			continue // skip non vertical, horizontal or 45deg diagonal lines
		}
		grid.addLine(line)
	}
	fmt.Println("answer", grid.getCountOfPointsThatAreCoveredAtLeastTwice())
}

func inputToLine(input string) Line {
	startAndEnd := strings.Split(input, " -> ")

	startString := strings.Split(startAndEnd[0], ",")
	endString := strings.Split(startAndEnd[1], ",")

	startX, _ := strconv.Atoi(startString[0])
	startY, _ := strconv.Atoi(startString[1])

	endX, _ := strconv.Atoi(endString[0])
	endY, _ := strconv.Atoi(endString[1])

	startCoordinate := Coordinate{
		x: startX,
		y: startY,
	}
	endCoordinate := Coordinate{
		x: endX,
		y: endY,
	}

	return Line{
		start: startCoordinate,
		end:   endCoordinate,
	}
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func createPointIndex(coordinate Coordinate) string {
	return strconv.Itoa(coordinate.x) + "," + strconv.Itoa(coordinate.y)
}
