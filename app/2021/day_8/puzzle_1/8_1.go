package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Entry struct {
	signalPatterns []string
	outputValues   []string
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

	uniqueSegmentsCount := 0
	for _, entry := range entries {
		for _, outputSegment := range entry.outputValues {
			switch len(outputSegment) {
			case 2: // digit 1
				fallthrough
			case 4: // digit 4
				fallthrough
			case 3: // digit 7
				fallthrough
			case 7: // digit 8
				uniqueSegmentsCount++
			}
		}
	}

	fmt.Println("uniqueSegmentsCount", uniqueSegmentsCount)
	fmt.Println("End", time.Since(start))
}

func parseInputLineToEntry(input string) Entry {
	inAndOutput := strings.Split(input, " | ")

	inputValues := strings.Split(inAndOutput[0], " ")
	outputValues := strings.Split(inAndOutput[1], " ")

	return Entry{inputValues, outputValues}
}
