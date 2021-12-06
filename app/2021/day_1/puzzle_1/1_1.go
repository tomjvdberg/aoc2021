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

	lastNumberSet := false
	lastNumber := 0

	greaterThanPreviousCount := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentNumber, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println(err)
		}

		// check if lastNumber is already set
		if !lastNumberSet {
			lastNumber = currentNumber
			lastNumberSet = true
			continue
		}

		if currentNumber > lastNumber {
			greaterThanPreviousCount++
		}

		lastNumber = currentNumber
	}

	fmt.Println(greaterThanPreviousCount)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}
