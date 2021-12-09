package main

import (
	"reflect"
	"testing"
)

func TestLettersToCollection(t *testing.T) {
	letterCollection := newLetterCollection("abcde")

	expectedCollection := LetterCollection{map[string]int{"b": 1, "a": 1, "c": 1, "d": 1, "e": 1}}

	if !reflect.DeepEqual(letterCollection, expectedCollection) {
		t.Log("Expected collections to be similar but got", letterCollection)
		t.Fail()
	}
}

func TestMinusLetters(t *testing.T) {
	letterCollection := newLetterCollection("abcde")
	minusCollection := LetterCollection{map[string]int{"d": 1, "e": 1}}

	expectedCollection := LetterCollection{map[string]int{"b": 1, "a": 1, "c": 1}}

	if !reflect.DeepEqual(letterCollection.minusLettersOf(minusCollection), expectedCollection) {
		t.Log("Expected collections to be similar but got", letterCollection)
		t.Fail()
	}
}
