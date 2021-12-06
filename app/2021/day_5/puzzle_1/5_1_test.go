package main

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestInputToLine(t *testing.T) {
	input := "301,306 -> 301,935"

	line := inputToLine(input)

	expectedStartCoordinate := Coordinate{
		x: 301,
		y: 306,
	}
	expectedEndCoordinate := Coordinate{
		x: 301,
		y: 935,
	}

	if line.start.x != expectedStartCoordinate.x || line.start.y != expectedStartCoordinate.y {
		fmt.Println(line)
		t.Log("Expected start coordinates to be the correct")
		t.Fail()
	}

	if line.end.x != expectedEndCoordinate.x || line.end.y != expectedEndCoordinate.y {
		t.Log("Expected end coordinates to be the correct")
		t.Fail()
	}
}

func TestDetectLineIsVertical(t *testing.T) {
	startCoordinate := Coordinate{
		x: 301,
		y: 306,
	}
	endCoordinate := Coordinate{
		x: 301,
		y: 935,
	}

	line := Line{
		start: startCoordinate,
		end:   endCoordinate,
	}

	if !line.isVertical() {
		t.Log("Expected to detect that the line is vertical")
		t.Fail()
	}
}

func TestDetectLineIsHorizontal(t *testing.T) {
	startCoordinate := Coordinate{
		x: 304,
		y: 306,
	}
	endCoordinate := Coordinate{
		x: 301,
		y: 306,
	}

	line := Line{
		start: startCoordinate,
		end:   endCoordinate,
	}

	if !line.isHorizontal() {
		t.Log("Expected to detect that the line is horizontal")
		t.Fail()
	}
}

func TestDetectInvalidLine(t *testing.T) {
	startCoordinate := Coordinate{
		x: 929,
		y: 976,
	}
	endCoordinate := Coordinate{
		x: 62,
		y: 109,
	}

	line := Line{
		start: startCoordinate,
		end:   endCoordinate,
	}

	if line.isHorizontal() {
		t.Log("Expected to detect that the line is not horizontal")
		t.Fail()
	}

	if line.isVertical() {
		t.Log("Expected to detect that the line is not vertical")
		t.Fail()
	}
}

func TestGetCoordinatesOfLine(t *testing.T) {
	startCoordinate := Coordinate{
		x: 304,
		y: 306,
	}
	endCoordinate := Coordinate{
		x: 301,
		y: 306,
	}

	line := Line{
		start: startCoordinate,
		end:   endCoordinate,
	}

	expectedCoordinates := []Coordinate{
		{304, 306},
		{303, 306},
		{302, 306},
		{301, 306},
	}

	if !reflect.DeepEqual(expectedCoordinates, line.getCoordinates()) {
		t.Log("Got unexpected results from the coordinates")
		t.Fail()
	}
}

func TestGetPoint(t *testing.T) {
	grid := Grid{points: []Point{{
		x:            2,
		y:            5,
		timesCovered: 1,
	}}}

	point, err := grid.getPoint(Coordinate{2, 5})

	if err != nil {
		t.Log("expected point but got error", err)
		t.Fail()
	}

	if !reflect.DeepEqual(*point, Point{2, 5, 1}) {
		t.Log("Points to be as expected but got", point)
		t.Fail()
	}
}

func TestGetPointWhenNotAvailable(t *testing.T) {
	grid := Grid{}

	_, err := grid.getPoint(Coordinate{})

	if !errors.Is(PointNotAvailable{}, err) {
		t.Log("expected PointNotAvailable error but got", err)
		t.Fail()
	}
}

func TestMarkPointWhenPointDoesNotExist(t *testing.T) {
	grid := Grid{}

	grid.markPoint(Coordinate{2, 5})

	if len(grid.points) == 0 {
		t.Log("expected point to be added but it was not")
		t.Fail()
	}

	if !reflect.DeepEqual(grid.points[0], Point{2, 5, 1}) {
		t.Log("Expected point to be in the grid and 1 time covered")
		t.Fail()
	}
}

func TestMarkPointThatAlreadyExists(t *testing.T) {
	grid := Grid{[]Point{{2, 5, 1}}}

	grid.markPoint(Coordinate{2, 5})

	if (grid.points[0]).timesCovered != 2 {
		t.Log("Expected point to be two times covered but it was not", grid.points[0])
		t.Fail()
	}
}

func TestAddLine(t *testing.T) {
	grid := Grid{}

	startCoordinate := Coordinate{
		x: 304,
		y: 306,
	}
	endCoordinate := Coordinate{
		x: 301,
		y: 306,
	}

	line := Line{
		start: startCoordinate,
		end:   endCoordinate,
	}

	grid.addLine(line)

	if len(grid.points) != 4 {
		t.Log("Expected 4 points on grid but got", len(grid.points))
		t.Fail()
	}
}

func TestAddOverlappingLine(t *testing.T) {
	grid := Grid{}

	horizontalLine := Line{Coordinate{4, 6}, Coordinate{1, 6}}
	verticalLine := Line{Coordinate{4, 4}, Coordinate{4, 6}}

	grid.addLine(horizontalLine)
	grid.addLine(verticalLine)

	if len(grid.points) != 6 {
		t.Log("Expected 6 points on grid but got", len(grid.points))
		t.Fail()
	}

	overlappingPoint, err := grid.getPoint(Coordinate{4, 6})
	if err != nil {
		t.Log("expected to find point but got error", err)
		t.Fail()
	}

	if overlappingPoint.timesCovered != 2 {
		t.Log("Expected overlappingPoint to be covered twice but got", overlappingPoint.timesCovered)
		t.Fail()
	}

	expectedCoveredMoreThanTwiceCount := 1
	countOfMoreThanTwiceCoveredPoints := grid.getCountOfPointsThatAreCoveredAtLeastTwice()

	if countOfMoreThanTwiceCoveredPoints != expectedCoveredMoreThanTwiceCount {
		t.Log("Expected count to be twice but got", countOfMoreThanTwiceCoveredPoints)
		t.Fail()
	}
}

func TestAddOverlappingHorizontalLineInverted(t *testing.T) {
	grid := Grid{}

	horizontalLine := Line{Coordinate{4, 6}, Coordinate{1, 6}}
	horizontalLineInverted := Line{Coordinate{1, 6}, Coordinate{4, 6}}

	grid.addLine(horizontalLine)
	grid.addLine(horizontalLineInverted)

	if len(grid.points) != 4 {
		t.Log("Expected 4 points on grid but got", len(grid.points))
		t.Fail()
	}

	overlappingPoint, err := grid.getPoint(Coordinate{4, 6})
	if err != nil {
		t.Log("expected to find point but got error", err)
		t.Fail()
	}

	if overlappingPoint.timesCovered != 2 {
		t.Log("Expected overlappingPoint to be covered twice but got", overlappingPoint.timesCovered)
		t.Fail()
	}

	expectedCoveredMoreThanTwiceCount := 4
	countOfMoreThanTwiceCoveredPoints := grid.getCountOfPointsThatAreCoveredAtLeastTwice()

	if countOfMoreThanTwiceCoveredPoints != expectedCoveredMoreThanTwiceCount {
		t.Log("Expected count to be four but got", countOfMoreThanTwiceCoveredPoints)
		t.Fail()
	}
}

func TestAddOverlappingVerticalLineInverted(t *testing.T) {
	grid := Grid{}

	verticalLine := Line{Coordinate{1, 1}, Coordinate{1, 4}}
	verticalLineInverted := Line{Coordinate{1, 4}, Coordinate{1, 1}}

	grid.addLine(verticalLine)
	grid.addLine(verticalLineInverted)

	if len(grid.points) != 4 {
		t.Log("Expected 4 points on grid but got", len(grid.points))
		t.Fail()
	}

	overlappingPoint, err := grid.getPoint(Coordinate{1, 4})
	if err != nil {
		t.Log("expected to find point but got error", err)
		t.Fail()
	}

	if overlappingPoint.timesCovered != 2 {
		t.Log("Expected overlappingPoint to be covered twice but got", overlappingPoint.timesCovered)
		t.Fail()
	}

	expectedCoveredMoreThanTwiceCount := 4
	countOfMoreThanTwiceCoveredPoints := grid.getCountOfPointsThatAreCoveredAtLeastTwice()

	if countOfMoreThanTwiceCoveredPoints != expectedCoveredMoreThanTwiceCount {
		t.Log("Expected count to be four but got", countOfMoreThanTwiceCoveredPoints)
		t.Fail()
	}
}
