package main

import (
	"strconv"
	"testing"
)

func TestBitChecking(t *testing.T) {
	binaryValue, err := strconv.ParseInt("0001", 2, 64)
	if err != nil {
		panic(err)
	}

	isClear := bitIsClearAtPosition(1, binaryValue)

	if !isClear {
		t.Log("expected bit to be clear but it was not")
		t.Fail()
	}

	isSet := bitIsSetAtPosition(0, binaryValue)

	if !isSet {
		t.Log("expected bit to be set but it was not")
		t.Fail()
	}
}
