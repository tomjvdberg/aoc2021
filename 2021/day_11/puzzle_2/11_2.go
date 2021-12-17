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

	octopusMap := map[string]int{}
	row := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		octopusesOnLine := lineToOctopusses(scanner.Text())
		for column, octopus := range octopusesOnLine {
			octopusMap[octopusIndex(row, column)] = octopus
		}
		row++
	}

	//desiredSteps := 1000000000
	desiredSteps := 2000
	flashesCount := 0
	step := 1
	for ; step < desiredSteps; step++ {
		flashRegister := map[string]bool{}
		increaseEnergyOfAllByOne(octopusMap)

		flashables := findFlashables(octopusMap)

		handleFlashes(&octopusMap, &flashRegister, flashables)

		if len(flashRegister) == 100 {
			break
		}
		flashesCount += len(flashRegister)
		flashRegister = map[string]bool{} // reset
	}
	fmt.Println("answer", step)
	fmt.Println("End", time.Since(start))
}

func findFlashables(octopusMap map[string]int) []string {
	var flashables []string
	for id, energy := range octopusMap {
		if energy > 9 {
			flashables = append(flashables, id)
		}
	}
	return flashables
}

func handleFlashes(m *map[string]int, fr *map[string]bool, flashables []string) {
	if len(flashables) < 1 {
		return
	}
	// foreach flasher increase the neighbor and set itself to 0
	// skip where is 0. That is only possible with ones that flashed.
	for _, flashableId := range flashables {
		(*m)[flashableId] = 0 // reset to zero
		(*fr)[flashableId] = true
		increaseNeighborsSkipFlashed(flashableId, m, fr)
	}

	newFlashables := findFlashables(*m)

	// now find new ones that are greater than 9
	handleFlashes(m, fr, newFlashables)
}

func increaseNeighborsSkipFlashed(index string, m *map[string]int, flashRegister *map[string]bool) {
	row, column := reverseIndex(index)
	neighbors := []string{
		octopusIndex(row-1, column-1), // tl
		octopusIndex(row-1, column),   // tm
		octopusIndex(row-1, column+1), // tr
		octopusIndex(row, column-1),   // l
		octopusIndex(row, column+1),   // r
		octopusIndex(row+1, column-1), // bl
		octopusIndex(row+1, column),   // bm
		octopusIndex(row+1, column+1), // br
	}

	for _, neighborId := range neighbors {
		_, alreadyFlashed := (*flashRegister)[neighborId]
		if alreadyFlashed {
			continue
		}

		nb, exists := (*m)[neighborId]
		if exists {
			(*m)[neighborId] = nb + 1
		}
	}
}

func increaseEnergyOfAllByOne(m map[string]int) {
	for octopusIndex, octopusEnergy := range m {
		(m)[octopusIndex] = octopusEnergy + 1
	}
}

func lineToOctopusses(lineAsText string) []int {
	var octopuses []int
	octopusesAsString := strings.Split(lineAsText, "")

	for _, octopusAsString := range octopusesAsString {
		intVal, _ := strconv.Atoi(octopusAsString)
		octopuses = append(octopuses, intVal)
	}

	return octopuses
}

func octopusIndex(row int, column int) string {
	return strconv.Itoa(row) + "," + strconv.Itoa(column)
}

func reverseIndex(index string) (row int, column int) {
	coords := strings.Split(index, ",")
	row, _ = strconv.Atoi(coords[0])
	column, _ = strconv.Atoi(coords[1])

	return
}
