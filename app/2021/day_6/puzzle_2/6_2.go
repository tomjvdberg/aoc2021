package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type FishGroup struct {
	count int
}

func main() {
	start := time.Now()
	file, err := os.Open("input")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	var fishGroups map[int]FishGroup

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fishGroups = inputToFish(scanner.Text())
	}

	daysToCalculate := 256
	for i := 0; i < daysToCalculate; i++ {
		nextDaysFishGroup := map[int]FishGroup{}

		fishToBeAddedToGroup6 := 0
		for timer, fishGroup := range fishGroups {
			if timer == 0 {
				nextDaysFishGroup[8] = FishGroup{fishGroup.count}
				// remember the sum of fishes for next days fishGroup 6
				fishToBeAddedToGroup6 = fishGroup.count
				continue
			}
			// move this fish to a lower group
			nextDaysFishGroup[timer-1] = fishGroup
		}

		group6, _ := nextDaysFishGroup[6]
		nextDaysFishGroup[6] = FishGroup{group6.count + fishToBeAddedToGroup6}
		fishGroups = nextDaysFishGroup
	}

	fmt.Printf("Ending with %d fishGroups \n", len(fishGroups))

	totalFish := 0
	for _, group := range fishGroups {
		totalFish += group.count
	}
	fmt.Printf("Ending with %d fish \n", totalFish)
	fmt.Println("End", time.Since(start))
}

func inputToFish(input string) map[int]FishGroup {
	fishGroups := map[int]FishGroup{}

	for i, _ := range [9]int{} {
		fishGroups[i] = FishGroup{0}
	}

	fishesAsString := strings.Split(input, ",")

	for _, fishAsString := range fishesAsString {
		fishTimer, _ := strconv.Atoi(fishAsString)

		fishGroup, exists := fishGroups[fishTimer]
		if !exists {
			panic("all groups should be present by now")
		}
		fishGroups[fishTimer] = FishGroup{fishGroup.count + 1}
	}

	return fishGroups
}
