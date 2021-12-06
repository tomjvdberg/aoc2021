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

	bufferFilled := false
	bufferCounter := 0
	var intArray [4]int

	greaterThanPreviousCount := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentNumber, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println(err)
		}

		if !bufferFilled {
			intArray[bufferCounter] = currentNumber
			bufferCounter++

			if bufferCounter < 4 {
				continue // buffer is not full yet
			}
			bufferCounter = 0 // reset the buffer counter
			bufferFilled = true
		} else {
			// first move 1->0 etc
			intArray[0] = intArray[1]
			intArray[1] = intArray[2]
			intArray[2] = intArray[3]
			// current number should now be placed at intArray[3]
			intArray[3] = currentNumber
		}

		// now compare ints 0,1,2 with ints 1,2,3
		firstWindow := intArray[0] + intArray[1] + intArray[2]
		secondWindow := intArray[1] + intArray[2] + intArray[3]

		if secondWindow > firstWindow {
			greaterThanPreviousCount++
		}
	}

	fmt.Println(greaterThanPreviousCount)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}
