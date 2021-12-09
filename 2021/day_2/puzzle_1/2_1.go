package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ExtractionError struct{}

func (ExtractionError) Error() string { return "error extracting values" }

func main() {
	file, err := os.Open("input")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	horizontalPosition := 0
	aim := 0
	depth := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		direction, distance, err := extractAction(scanner.Text())
		if err != nil {
			panic(err)
		}

		if distance < 0 {
			panic("distance should not be negative. Got:" + strconv.Itoa(distance))
		}

		switch direction {
		case "forward":
			horizontalPosition += distance
			depth += distance * aim
		case "up":
			aim -= distance
		case "down":
			aim += distance
		default:
			panic("Direction not valid: " + direction)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	fmt.Println(horizontalPosition, depth)
	fmt.Println(horizontalPosition * depth)
}

func extractAction(action string) (string, int, error) {
	actionParts := strings.Split(action, " ")

	if len(actionParts) != 2 {
		return "", 0, ExtractionError{}
	}

	direction := actionParts[0]
	distance, err := strconv.Atoi(actionParts[1])

	if err != nil {
		return "", 0, err
	}

	return direction, distance, err
}
