package main

import (
	"testing"
)

func TestLineIsDiagonal(t *testing.T) {
	cases := []Line{
		{Coordinate{1, 4}, Coordinate{3, 6}},
		{Coordinate{3, 4}, Coordinate{1, 6}},
		{Coordinate{1, 6}, Coordinate{3, 4}},
		{Coordinate{461, 425}, Coordinate{929, 893}},
		{Coordinate{893, 738}, Coordinate{412, 257}},
	}

	for _, line := range cases {
		if !line.isDiagonal() {
			t.Log("Expected line to be diagonal but it was not", line)
			t.Fail()
		}
	}
}

func TestGetDiagonalCoordinates(t *testing.T) {
	line := Line{Coordinate{3, 4}, Coordinate{7, 8}}
	if len(line.getDiagonalCoordinates()) != 5 {
		t.Log("Expected line to have 3 coordinates but there were not", line.getDiagonalCoordinates())
		t.Fail()
	}
}
