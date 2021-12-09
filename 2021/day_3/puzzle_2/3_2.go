package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	file, err := os.Open("input")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	var binaryValues []int64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		binaryValue, err := strconv.ParseInt(text, 2, 64)
		if err != nil {
			panic(err)
		}
		binaryValues = append(binaryValues, binaryValue)
	}

	counter := [12]int64{11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}

	// Oxygen rating
	var oxygenGeneratorRating int64
	filteredOxygenRatingBinaryValues := append(binaryValues)

	for _, bitPosition := range counter {
		falseBitCount := countClearBitsAtPosition(filteredOxygenRatingBinaryValues, bitPosition)
		trueBitCount := countSetBitsAtPosition(filteredOxygenRatingBinaryValues, bitPosition)

		mostCommonBitValueForThisPosition := true // true has priority
		if falseBitCount > trueBitCount {
			mostCommonBitValueForThisPosition = false
		}

		filtered := keepOnlyBinaryValuesWithValueAtPosition(filteredOxygenRatingBinaryValues, mostCommonBitValueForThisPosition, bitPosition)
		if len(filtered) < 1 {
			continue // there were no results so skip this step
		}
		filteredOxygenRatingBinaryValues = filtered

		if len(filteredOxygenRatingBinaryValues) == 1 {
			fmt.Println("Found the oxygen generator rating")
			fmt.Printf("%012b, %d\n", filteredOxygenRatingBinaryValues[0], filteredOxygenRatingBinaryValues[0])
			oxygenGeneratorRating = filteredOxygenRatingBinaryValues[0]
			break
		}
	}

	// CO2 scrubber rating
	var co2ScrubberRating int64
	filteredCo2ScrubberRatingBinaryValues := append(binaryValues)

	for _, bitPosition := range counter {
		falseBitCount := countClearBitsAtPosition(filteredCo2ScrubberRatingBinaryValues, bitPosition)
		trueBitCount := countSetBitsAtPosition(filteredCo2ScrubberRatingBinaryValues, bitPosition)

		leastCommonBitValueForThisPosition := false // false has priority
		if falseBitCount > trueBitCount {
			leastCommonBitValueForThisPosition = true
		}

		filtered := keepOnlyBinaryValuesWithValueAtPosition(filteredCo2ScrubberRatingBinaryValues, leastCommonBitValueForThisPosition, bitPosition)
		if len(filtered) < 1 {
			continue // there were no results so skip this step
		}
		filteredCo2ScrubberRatingBinaryValues = filtered

		if len(filteredCo2ScrubberRatingBinaryValues) == 1 {
			fmt.Println("Found the CO2 scrubber rating")
			fmt.Printf("%012b, %d\n", filteredCo2ScrubberRatingBinaryValues[0], filteredCo2ScrubberRatingBinaryValues[0])
			co2ScrubberRating = filteredCo2ScrubberRatingBinaryValues[0]
			break
		}
	}
	fmt.Println("Life support rating", oxygenGeneratorRating*co2ScrubberRating)
}

func countSetBitsAtPosition(binaryValues []int64, pos int64) int {
	countOfBitsSet := 0
	for _, binaryValue := range binaryValues {
		if bitIsSetAtPosition(pos, binaryValue) {
			countOfBitsSet++
		}
	}

	return countOfBitsSet
}

func countClearBitsAtPosition(binaryValues []int64, pos int64) int {
	countOfBitsClear := 0
	for _, binaryValue := range binaryValues {
		if bitIsClearAtPosition(pos, binaryValue) {
			countOfBitsClear++
		}
	}

	return countOfBitsClear
}

func keepOnlyBinaryValuesWithValueAtPosition(binaryValues []int64, mostCommonBitValueForThisPosition bool, binaryPosition int64) []int64 {
	var filtered []int64

	for _, binaryValue := range binaryValues {
		if bitIsSetAtPosition(binaryPosition, binaryValue) == mostCommonBitValueForThisPosition {
			filtered = append(filtered, binaryValue)
		}
	}

	return filtered
}

func bitIsSetAtPosition(bitPosition int64, binaryValue int64) bool {
	return 1 == ((binaryValue >> bitPosition) & 1)
}

func bitIsClearAtPosition(bitPosition int64, binaryValue int64) bool {
	return 0 == ((binaryValue >> bitPosition) & 1)
}
