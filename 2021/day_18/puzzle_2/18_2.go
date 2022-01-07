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

	highestMagnitude := 0
	nrList := [][]number{}
	for scanner.Scan() {
		nrList = append(nrList, parseInput(scanner.Text()))
	}

	for fi, firstNr := range nrList {
		for si, secondNr := range nrList {
			if fi == si {
				continue
			}
			added := add(firstNr, secondNr)
			reduced := reduce(added)
			mag := magnitudeOf(reduced)
			if mag > highestMagnitude {
				highestMagnitude = mag
			}
		}
	}

	fmt.Println("highestMagnitude:", highestMagnitude)
	fmt.Println("End", time.Since(start))
}

func add(A []number, B []number) []number {
	newNrs := make([]number, len(A)+len(B))
	copy(newNrs, append(A, B...))

	if len(A) == 0 {
		return newNrs
	}

	for i, nr := range newNrs {
		newNrs[i].level = nr.level + 1
	}

	return newNrs
}

func reduce(nrs []number) []number {
	nrsCopy := *copyNrs(nrs)

	for {
		pairAtLevel4Found := false
		level4PairIdx := 0

		for nrIdx, nr := range nrsCopy {
			if nr.level == 4 {
				pairAtLevel4Found = true
				level4PairIdx = nrIdx
				break
			}
		}

		if pairAtLevel4Found {
			nrsCopy = explodePair(level4PairIdx, nrsCopy)
			continue
		}

		splitFound := false
		splitNrIdx := 0
		for nrIdx, nr := range nrsCopy {
			if nr.value > 9 {
				splitFound = true
				splitNrIdx = nrIdx
				break
			}
		}

		if !splitFound {
			break
		}
		nrsCopy = splitNr(splitNrIdx, nrsCopy)
	}

	return nrsCopy
}

func splitNr(splitNrIdx int, nrs []number) []number {
	nrsCopy := *copyNrs(nrs)

	valA := nrsCopy[splitNrIdx].value / 2
	valB := nrsCopy[splitNrIdx].value - (nrsCopy[splitNrIdx].value / 2)

	A := number{nrsCopy[splitNrIdx].level + 1, valA}
	B := number{nrsCopy[splitNrIdx].level + 1, valB}

	bef := make([]number, len(nrsCopy[:splitNrIdx+1]))
	copy(bef, nrsCopy[:splitNrIdx+1])

	aft := make([]number, len(nrsCopy[splitNrIdx+1:]))
	copy(aft, nrsCopy[splitNrIdx+1:])

	bef = append(nrsCopy[:splitNrIdx+1], B)
	nrsCopy = append(bef, aft...)
	nrsCopy[splitNrIdx] = A

	return nrsCopy
}

func explodePair(level4PairIdx int, nrs []number) []number {
	nrsCopy := *copyNrs(nrs)
	// increment left neighbor
	leftNeighborIdx := level4PairIdx - 1
	if leftNeighborIdx >= 0 {
		nrsCopy[leftNeighborIdx].value += nrsCopy[level4PairIdx].value
	}
	// increment right neighbor with value of second level4 element
	rightNeighborIdx := level4PairIdx + 2
	if rightNeighborIdx <= (len(nrsCopy) - 1) {
		nrsCopy[rightNeighborIdx].value += nrsCopy[level4PairIdx+1].value
	}

	// remove the second element of the L4 pair
	// set the first element to be a level higher and value 0
	replacementNr := number{3, 0}

	return replacePairWith(replacementNr, level4PairIdx, nrsCopy)
}

func replacePairWith(replacementNr number, idx int, nrs []number) []number {
	nrsCopy := *copyNrs(nrs)
	nrsCopy = append(nrsCopy[:idx], nrsCopy[idx+1:]...)
	nrsCopy[idx] = replacementNr

	return nrsCopy
}

func magnitudeOf(nrs []number) int {
	nrsCopy := *copyNrs(nrs)

	inc0 := 3
	inc1 := 2

	for {
		if len(nrsCopy) == 2 {
			break
		}

		firstElementPassed := false
		lastNrLevel := 0
		pairFound := false
		pairStartIdx := -1

		for idx, nr := range nrsCopy {
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
				(inc0 * nrsCopy[pairStartIdx].value) + (inc1 * nrsCopy[pairStartIdx+1].value),
			}
			nrsCopy = replacePairWith(repl, pairStartIdx, nrsCopy)
		}
	}

	return (inc0 * nrsCopy[0].value) + (inc1 * nrsCopy[1].value)
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

func copyNrs(nrs []number) *[]number {
	cp := make([]number, len(nrs))
	copy(cp, nrs)

	return &cp
}
