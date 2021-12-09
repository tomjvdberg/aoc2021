package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Entry struct {
	signalPatterns []string
	outputValues   []string
}

type LetterCollection struct {
	letters map[string]int
}

func (collection LetterCollection) minusLettersOf(minusCollection LetterCollection) LetterCollection {
	returnCollection := LetterCollection{map[string]int{}}
	for letter, _ := range collection.letters {
		// if the letter appears in the minusCollection don't add it to the return collection
		_, existsInMinusCollection := minusCollection.letters[letter]
		if !existsInMinusCollection {
			returnCollection.letters[letter] = 1
		}
	}

	return returnCollection
}

func (collection LetterCollection) containsLettersOf(testCollection LetterCollection) bool {
	containsAll := true
	for letter, _ := range testCollection.letters {
		// if the letter appears in the testCollection don't add it to the return collection
		_, existsInMinusCollection := collection.letters[letter]
		if !existsInMinusCollection {
			containsAll = false
		}
	}

	return containsAll
}

func (collection LetterCollection) hasExactlyOneOf(testCollection LetterCollection) bool {
	matchCount := 0
	for letter, _ := range testCollection.letters {
		// if the letter appears in the testCollection don't add it to the return collection
		_, existsInMinusCollection := collection.letters[letter]
		if !existsInMinusCollection {
			matchCount++
		}
	}

	return matchCount == 1
}

func newLetterCollection(letters string) LetterCollection {
	collection := LetterCollection{map[string]int{}}
	for _, letter := range letters {
		collection.letters[string(letter)] = 1
	}

	return collection
}

func main() {
	start := time.Now()
	file, err := os.Open("input")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	var entries []Entry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entries = append(entries, parseInputLineToEntry(scanner.Text()))
	}

	var outputNrs []int

	for _, entry := range entries {
		var zero LetterCollection
		var one LetterCollection
		var two LetterCollection
		var three LetterCollection
		var four LetterCollection
		var five LetterCollection
		var six LetterCollection
		var seven LetterCollection
		var eight LetterCollection
		var nine LetterCollection

		var A LetterCollection
		var B LetterCollection
		var C LetterCollection
		var D LetterCollection

		for _, signalPattern := range entry.signalPatterns {
			switch len(signalPattern) {
			case 2:
				one = newLetterCollection(signalPattern)
			case 3:
				seven = newLetterCollection(signalPattern)
			case 4:
				four = newLetterCollection(signalPattern)
			case 7:
				eight = newLetterCollection(signalPattern)
			}
		}

		// Explanation
		//		  __B__
		//		 |     |
		//		 C     A
		//		 |     |
		//		  __C__
		//		 |     |
		//		 D     A
		//		 |     |
		//		  __D__
		// A = letters of digit 1 (the signal pattern with only two letters)
		// B = letters of digit 7 (the one with three letters) minus the letters of A (or one as they are the same)
		// C = letters of digit 4 (the one with four letters) minus the letters of A
		// D = letters of digit 8 (the one with seven letters) minus letters of digit 4 and minus letters of digit 7

		// So we have 1,4,7 and 8. All other numbers are found in the signals with a length of 5 or 6.
		// Five will have B for sure. There is no other digit with 5 letters that has both letters of C.
		// So a signal with a length of 5 that contains the letter B and both of C must be five.
		//
		// Three will have B for sure. There is no other digit with 5 letters that has both letters of A.
		// So a signal with a length of 5 that contains the letter B and both of A must be three.
		//
		// Two will have B for sure. It is the signal with the letter of B and one of A (but not both), and one of C (but not both).

		// Then the six length signals. They can be either 0, 6 or 9.
		// 0 can only be a signal that has only one of C (since the middle element would not be active)
		// 6 can only be a signal that has only one of A (since the top right element would not be active)
		// 9 can only be a signal that has only one of D (since the bottom left element would not be active)

		// Now we can decode the output patterns

		A = one
		B = seven.minusLettersOf(one)
		C = four.minusLettersOf(one)
		D = eight.minusLettersOf(four).minusLettersOf(seven)

		for _, signalPattern := range entry.signalPatterns {
			spColl := newLetterCollection(signalPattern)
			switch len(signalPattern) {
			case 5:
				if spColl.containsLettersOf(B) && spColl.containsLettersOf(C) {
					five = spColl
					continue
				}
				if spColl.containsLettersOf(B) && spColl.containsLettersOf(A) {
					three = spColl
					continue
				}
				if spColl.containsLettersOf(B) && (spColl.hasExactlyOneOf(A) && spColl.hasExactlyOneOf(C)) {
					two = spColl
					continue
				}
			case 6:
				if spColl.hasExactlyOneOf(C) {
					zero = spColl
					continue
				}
				if spColl.hasExactlyOneOf(A) {
					six = spColl
					continue
				}
				if spColl.hasExactlyOneOf(D) {
					nine = spColl
					continue
				}
			}
		}

		var outputValues []int
		for _, outputPattern := range entry.outputValues {
			opColl := newLetterCollection(outputPattern)
			if reflect.DeepEqual(opColl, zero) {
				outputValues = append(outputValues, 0)
				continue
			}
			if reflect.DeepEqual(opColl, one) {
				outputValues = append(outputValues, 1)
				continue
			}
			if reflect.DeepEqual(opColl, two) {
				outputValues = append(outputValues, 2)
				continue
			}
			if reflect.DeepEqual(opColl, three) {
				outputValues = append(outputValues, 3)
				continue
			}
			if reflect.DeepEqual(opColl, four) {
				outputValues = append(outputValues, 4)
				continue
			}
			if reflect.DeepEqual(opColl, five) {
				outputValues = append(outputValues, 5)
				continue
			}
			if reflect.DeepEqual(opColl, six) {
				outputValues = append(outputValues, 6)
				continue
			}
			if reflect.DeepEqual(opColl, seven) {
				outputValues = append(outputValues, 7)
				continue
			}
			if reflect.DeepEqual(opColl, eight) {
				outputValues = append(outputValues, 8)
				continue
			}
			if reflect.DeepEqual(opColl, nine) {
				outputValues = append(outputValues, 9)
				continue
			}

		}

		output := fmt.Sprintf("%d%d%d%d", outputValues[0], outputValues[1], outputValues[2], outputValues[3])
		intVal, _ := strconv.Atoi(output)
		outputNrs = append(outputNrs, intVal)
	}

	sum := 0
	for _, outputNr := range outputNrs {
		sum += outputNr
	}

	fmt.Println("sum", sum)
	fmt.Println("End", time.Since(start))
}

func parseInputLineToEntry(input string) Entry {
	inAndOutput := strings.Split(input, " | ")

	inputValues := strings.Split(inAndOutput[0], " ")
	outputValues := strings.Split(inAndOutput[1], " ")

	return Entry{inputValues, outputValues}
}
