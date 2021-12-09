package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Fish struct {
	timer int
}

func main() {
	start := time.Now()
	file, err := os.Open("input")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	var fishes []Fish

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fishes = inputToFish(scanner.Text())
	}
	fmt.Printf("Starting with %d fishes \n", len(fishes))

	daysToCalculate := 80

	for i := 0; i < daysToCalculate; i++ {
		var fishesToBeAdded []Fish
		for i := 0; i < len(fishes); i++ {
			currentFish := &fishes[i]
			if (*currentFish).timer == 0 {
				fishesToBeAdded = append(fishesToBeAdded, createNewFish())
				(*currentFish).timer = 6
				continue
			}
			(*currentFish).timer--
		}
		fishes = append(fishes, fishesToBeAdded...)
	}
	fmt.Printf("Ending with %d fishes \n", len(fishes))
	fmt.Println("End", time.Since(start))
}

func createNewFish() Fish {
	return Fish{8}
}

func inputToFish(input string) []Fish {
	var fishes []Fish
	fishesAsString := strings.Split(input, ",")

	for _, fishAsString := range fishesAsString {
		fish, _ := strconv.Atoi(fishAsString)
		fishes = append(fishes, Fish{fish})
	}

	return fishes
}
