package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PointNotAvailable struct{}

func (PointNotAvailable) Error() string { return "point not available" }

type Grid struct {
	points []Point
}

func (grid *Grid) addLine(line Line) {
	coordinates := line.getCoordinates()

	for _, coordinate := range coordinates {
		grid.markPoint(coordinate)
	}
}

func (grid *Grid) markPoint(coordinate Coordinate) {
	point, err := grid.getPoint(coordinate)

	if err != nil {
		// point does not exist, add it
		grid.points = append(grid.points, Point{coordinate.x, coordinate.y, 1})

		return
	}

	(*point).timesCovered++
}

func (grid *Grid) getPoint(coordinate Coordinate) (*Point, error) {
	for i := 0; i < len((*grid).points); i++ {
		point := &(grid.points[i])
		if (*point).x == coordinate.x && (*point).y == coordinate.y {
			return point, nil
		}
	}

	return &Point{}, PointNotAvailable{}
}

func (grid *Grid) getCountOfPointsThatAreCoveredAtLeastTwice() int {
	countOfPointsThatWereCoveredTwiceOrMore := 0

	for i := 0; i < len(grid.points); i++ {
		if grid.points[i].timesCovered >= 2 {
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

func (line Line) getCoordinates() []Coordinate {
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

func (line Line) isHorizontal() bool {
	return line.start.y == line.end.y
}

func (line Line) isVertical() bool {
	return line.start.x == line.end.x
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
	var grid Grid

	addedLines := 0
	for _, line := range lines {
		if !(line.isHorizontal() || line.isVertical()) {
			continue // skip non vertical or horizontal lines
		}
		grid.addLine(line)
		addedLines++
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
