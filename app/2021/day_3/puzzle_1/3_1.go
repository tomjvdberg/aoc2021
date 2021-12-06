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

	totalValuesCount := 0
	var bitCounter [12]int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		binaryValue, err := strconv.ParseInt(text, 2, 64)
		if err != nil {
			panic(err)
		}
		totalValuesCount++

		for index, _ := range bitCounter {
			if bitIsSetAtPosition(index, binaryValue) {
				bitCounter[indexToBinaryPosition(index, 12)] = bitCounter[indexToBinaryPosition(index, 12)] + 1
			}
		}
	}

	fmt.Println(totalValuesCount)
	fmt.Println(bitCounter)

	majorityCount := (totalValuesCount / 2) + 1

	// now create a new byte
	gammaRate := 0b000000000000
	for index, trueBitCount := range bitCounter {
		if trueBitCount >= majorityCount {
			// there are more 1's than 0's
			gammaRate = setBit(gammaRate, indexToBinaryPosition(index, 12))
		}
	}

	epsilonRate := 0b111111111111 ^ gammaRate

	fmt.Printf("%012b\n", gammaRate)
	fmt.Printf("%012b\n", epsilonRate)
	fmt.Printf("Gamma %d\n", gammaRate)
	fmt.Printf("Epsilon %d\n", epsilonRate)
	fmt.Printf("Power consumption %d\n", gammaRate*epsilonRate)
}

func indexToBinaryPosition(index int, arrayLength int) int {
	highestArrayPosition := arrayLength - 1

	return highestArrayPosition - index
}

func bitIsSetAtPosition(bitPosition int, binaryValue int64) bool {
	return 1 == ((binaryValue >> bitPosition) & 1)
}

func setBit(binaryValue int, bitPosition int) int {
	binaryValue |= 1 << bitPosition

	return binaryValue
}
