package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type number struct {
	level int
	value int
}

func main() {
	start := time.Now()
	file, err := os.Open("input")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var nrs []number
	for scanner.Scan() {
		nrs = reduce(add(nrs, parseInput(scanner.Text())))
	}

	fmt.Println("magnitude:", magnitudeOf(nrs))
	fmt.Println("End", time.Since(start))
}

func add(nrs []number, input []number) []number {
	if len(nrs) == 0 {
		return input
	}
	newNrs := append(nrs, input...)

	for i, nr := range newNrs {
		newNrs[i].level = nr.level + 1
	}

	return newNrs
}

func reduce(nrs []number) []number {
	for {
		pairAtLevel4Found := false
		level4PairIdx := 0

		for nrIdx, nr := range nrs {
			if nr.level == 4 {
				pairAtLevel4Found = true
				level4PairIdx = nrIdx
				break
			}
		}

		if pairAtLevel4Found {
			nrs = explodePair(level4PairIdx, nrs)
			continue
		}

		splitFound := false
		splitNrIdx := 0
		for nrIdx, nr := range nrs {
			if nr.value > 9 {
				splitFound = true
				splitNrIdx = nrIdx
				break
			}
		}

		if !splitFound {
			break
		}
		nrs = splitPair(splitNrIdx, nrs)
	}
	return nrs
}

func splitPair(splitNrIdx int, nrs []number) []number {
	valA := nrs[splitNrIdx].value / 2
	valB := nrs[splitNrIdx].value - (nrs[splitNrIdx].value / 2)

	A := number{nrs[splitNrIdx].level + 1, valA}
	B := number{nrs[splitNrIdx].level + 1, valB}

	bef := make([]number, len(nrs[:splitNrIdx+1]))
	copy(bef, nrs[:splitNrIdx+1])
	aft := make([]number, len(nrs[splitNrIdx+1:]))
	copy(aft, nrs[splitNrIdx+1:])
	bef = append(nrs[:splitNrIdx+1], B)

	nrs[splitNrIdx] = A
	nrs = append(bef, aft...)

	return nrs
}

func explodePair(level4PairIdx int, nrs []number) []number {
	// increment left neighbor
	leftNeighborIdx := level4PairIdx - 1
	if leftNeighborIdx >= 0 {
		nrs[leftNeighborIdx].value += nrs[level4PairIdx].value
	}
	// increment right neighbor with value of second level4 element
	rightNeighborIdx := level4PairIdx + 2
	if rightNeighborIdx <= (len(nrs) - 1) {
		nrs[rightNeighborIdx].value += nrs[level4PairIdx+1].value
	}

	// remove the second element of the L4 pair
	// set the first element to be a level higher and value 0
	replacementNr := number{3, 0}
	replacePairWith(replacementNr, level4PairIdx, &nrs)

	return nrs
}

func replacePairWith(replacementNr number, idx int, nrs *[]number) {
	*nrs = append((*nrs)[:idx], (*nrs)[idx+1:]...)
	(*nrs)[idx] = replacementNr
}

func magnitudeOf(nrs []number) int {
	inc0 := 3
	inc1 := 2

	for {
		if len(nrs) == 2 {
			break
		}

		firstElementPassed := false
		lastNrLevel := 0
		pairFound := false
		pairStartIdx := -1

		for idx, nr := range nrs {
			if !firstElementPassed {
				lastNrLevel = nr.level
				firstElementPassed = true
				continue
			}
			// pair found starting at idx -1
			if nr.level == lastNrLevel {
				pairFound = true
				pairStartIdx = idx - 1
				break
			}
			lastNrLevel = nr.level
		}
		if pairFound {
			repl := number{
				lastNrLevel - 1,
				(inc0 * nrs[pairStartIdx].value) + (inc1 * nrs[pairStartIdx+1].value),
			}
			replacePairWith(repl, pairStartIdx, &nrs)
		}
	}

	return (inc0 * nrs[0].value) + (inc1 * nrs[1].value)
}

func parseInput(input string) []number {
	var nrs []number
	currentLevel := -1
	parts := strings.Split(input, "")
	for _, part := range parts {
		switch part {
		case "[":
			currentLevel++
		case ",":
			// noop
		case "]":
			currentLevel--
		default: // it is a number
			nrVal, _ := strconv.Atoi(part)
			nrs = append(nrs, number{currentLevel, nrVal})
		}
	}

	return nrs
}
