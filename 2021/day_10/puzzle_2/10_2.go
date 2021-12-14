package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"math"
	"sort"
)

func main() {
	start := time.Now()
	file, err := os.Open("input")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var incompleteLines [][]string
	for scanner.Scan() {
		chunkStack := []string{}
		chunkChars := strings.Split(scanner.Text(), "")
		isIllegalLine := false
		for _, char := range chunkChars {
			if isOpeningChar(char) {
				chunkStack = append(chunkStack, char)
				continue
			}
			lastOpenedChunkChar := chunkStack[len(chunkStack)-1]
			if chunkCanBeClosedWith(lastOpenedChunkChar, char) {
				chunkStack = chunkStack[:len(chunkStack)-1] // remove from stack
				continue
			}
			isIllegalLine = true
		}
		if !isIllegalLine {
			incompleteLines = append(incompleteLines, chunkChars)
		}

	}
	var scores []int
	for _, incompleteLineChars := range incompleteLines {
		expectingClosingChars := []string{}
		for _, char := range incompleteLineChars {
			if isOpeningChar(char) {
				expectingClosingChars = append(expectingClosingChars, closingCharOf(char))
				continue
			}
			expectingClosingChars = expectingClosingChars[:len(expectingClosingChars)-1]
		}
		reverse := reverseSlice(expectingClosingChars)
		score := scoreLine(reverse)
		scores = append(scores, score)
	}

	sort.Ints(scores)
	middle := int(math.Round(float64(len(scores)) / 2.0))
	fmt.Println("answer", scores[int(middle)-1])
	fmt.Println("End", time.Since(start))
}

func scoreLine(closingChars []string) int {
	score := 0
	for _, char := range closingChars {
		score = score * 5
		score += valueOf(char)
	}
	return score
}

func valueOf(char string) int {
	switch char {
	case ")":
		return 1
	case "]":
		return 2
	case "}":
		return 3
	case ">":
		return 4
	}
	panic("dont be here")
}

func reverseSlice(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func closingCharOf(openeningChar string) string {
	if openeningChar == "(" {
		return ")"
	}
	if openeningChar == "{" {
		return "}"
	}
	if openeningChar == "[" {
		return "]"
	}
	if openeningChar == "<" {
		return ">"
	}

	panic("You should never come here.")
}

func chunkCanBeClosedWith(openeningChar string, closingChar string) bool {
	if openeningChar == "(" {
		return closingChar == ")"
	}
	if openeningChar == "{" {
		return closingChar == "}"
	}
	if openeningChar == "[" {
		return closingChar == "]"
	}
	if openeningChar == "<" {
		return closingChar == ">"
	}

	panic("You should never come here.")
}

func isOpeningChar(char string) bool {
	return char == "(" || char == "[" || char == "{" || char == "<"
}
