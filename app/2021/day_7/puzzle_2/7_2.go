package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	file, err := os.Open("input")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	var crabsGroupedByLevel map[int]int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		crabsGroupedByLevel = inputToLevelGroups(scanner.Text())
	}

	lowestLevel := 0 // I can see this from the input
	highestLevel := 0
	for level, _ := range crabsGroupedByLevel {
		if level > highestLevel {
			highestLevel = level
		}
	}

	maxLevelVariation := highestLevel - lowestLevel
	lowestCost := 0
	for i := 0; i < maxLevelVariation; i++ {
		cost := calculateCost(i, crabsGroupedByLevel)
		if i == 0 {
			lowestCost = cost
			continue
		}
		if cost < lowestCost {
			lowestCost = cost
		}
	}

	fmt.Println("lowestCost", lowestCost)
	fmt.Println("End", time.Since(start))
}

func calculateCost(nextLevel int, crabsGroupedByLevel map[int]int) int {
	cost := 0
	for level, crabCount := range crabsGroupedByLevel {
		if level == nextLevel {
			continue // no costs
		}

		cost += stepCostCalculator(abs(level-nextLevel)) * crabCount
	}

	return cost
}

func stepCostCalculator(steps int) int {
	cost := 0
	for i := 0; i < steps; i++ {
		cost = cost + (i + 1)
	}

	return cost
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func inputToLevelGroups(input string) map[int]int {
	crabsGroupedByLevel := map[int]int{}

	levelsAsString := strings.Split(input, ",")
	for _, levelAsString := range levelsAsString {
		level, _ := strconv.Atoi(levelAsString)

		levelGroupCount, exists := crabsGroupedByLevel[level]

		if exists {
			crabsGroupedByLevel[level] = levelGroupCount + 1
			continue
		}
		crabsGroupedByLevel[level] = 1
	}

	return crabsGroupedByLevel
}
