package main

import (
	"errors"
	"strconv"
	"testing"
)

func TestExtractAction(t *testing.T) {
	cases := []struct {
		input     string
		direction string
		distance  int
		err       error
	}{
		{"forward 8", "forward", 8, nil},
		{"down 1", "down", 1, nil},
		{"forward8", "", 0, ExtractionError{}},
		{"down 6 7", "", 0, ExtractionError{}},
		{"down -6", "down", -6, nil},
		{"up someStringText", "", 0, strconv.ErrSyntax},
	}

	for _, c := range cases {
		direction, distance, err := extractAction(c.input)

		if !errors.Is(err, c.err) {
			t.Log("Error expected to be | ", c.err, " | but got", err)
			t.Fail()
		}

		if direction != c.direction {
			t.Log("Direction should be", c.direction, "but got:", direction)
			t.Fail()
		}

		if distance != c.distance {
			t.Log("Distance should be", c.distance, "but got", distance)
			t.Fail()
		}
	}
}
