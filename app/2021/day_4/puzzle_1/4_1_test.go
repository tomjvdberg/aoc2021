package main

import (
	"errors"
	"testing"
)

func TestBoardIsCompleteWhenRowIsComplete(t *testing.T) {
	board := Board{[]BoardPosition{
		{1, 1, 22, false},
		{1, 2, 13, false},
		{1, 3, 17, false},
		{1, 4, 11, false},
		{1, 5, 0, false},
	}}

	board.markNumber(22)
	board.markNumber(13)
	board.markNumber(17)
	board.markNumber(11)
	board.markNumber(0)

	if !board.isComplete() {
		t.Log("expected row to be complete but it was not")
		t.Fail()
	}
}

func TestBoardIsCompleteWhenColumnIsComplete(t *testing.T) {
	board := Board{[]BoardPosition{
		{1, 1, 22, false},
		{2, 1, 13, false},
		{3, 1, 17, false},
		{4, 1, 11, false},
		{5, 1, 0, false},
	}}

	board.markNumber(22)
	board.markNumber(13)
	board.markNumber(17)
	board.markNumber(11)
	board.markNumber(0)

	if !board.isComplete() {
		t.Log("expected row to be complete but it was not")
		t.Fail()
	}
}

func TestErrorWhenQueryingForNonAvailableColumn(t *testing.T) {
	board := Board{[]BoardPosition{
		{1, 1, 22, false},
		{2, 1, 13, false},
		{3, 1, 17, false},
		{4, 1, 11, false},
		{5, 1, 0, false},
	}}

	_, err := board.getColumn(2)

	if !errors.Is(ColumnNotAvailable{}, err) {
		t.Log("expected ColumnNotAvailable error but got", err)
		t.Fail()
	}
}

func TestErrorWhenQueryingColumnWithInvalidLength(t *testing.T) {
	board := Board{[]BoardPosition{
		{1, 1, 22, false},
		{2, 1, 13, false},
		{3, 1, 17, false},
		{4, 1, 11, false},
		{5, 2, 0, false},
	}}

	_, err := board.getColumn(1)

	if !errors.Is(InvalidColumnLength{}, err) {
		t.Log("expected InvalidColumnLength error but got", err)
		t.Fail()
	}
}
