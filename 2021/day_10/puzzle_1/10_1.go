package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type illegalCharCounter struct {
	a int // )
	b int // }
	c int // ]
	d int // >
}

func (counter *illegalCharCounter) add(char string) {
	if char == ")" {
		counter.a++
	}
	if char == "}" {
		counter.b++
	}
	if char == "]" {
		counter.c++
	}
	if char == ">" {
		counter.d++
	}
}

func main() {
	start := time.Now()
	file, err := os.Open("input")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	counter := illegalCharCounter{}
	// last opened should be first closed
	chunkStack := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		chunkChars := strings.Split(scanner.Text(), "")
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

			// illegal character found
			counter.add(char)
			break
		}

	}

	fmt.Println("answer", (counter.a*3)+(counter.c*57)+(counter.b*1197)+(counter.d*25137))
	fmt.Println("End", time.Since(start))
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
