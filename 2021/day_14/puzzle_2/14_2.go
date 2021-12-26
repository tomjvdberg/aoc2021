package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type rule struct {
	inserts   string
	nextPairs []string
}

type pair string
type pairCounter map[pair]int
type pairGraph struct {
	pair        pair
	pairCounter pairCounter
	bottomPairs pairCounter
}
type pairCache map[pair][]pairGraph

func createPairCache(cacheSize int, rules map[pair]rule) pairCache {
	pairCache_ := map[pair][]pairGraph{}
	// init the cache
	for p := range rules {
		pairCache_[p] = []pairGraph{}
		for i := 0; i < cacheSize; i++ {
			pairCache_[p] = append(pairCache_[p], pairGraph{
				pair:        p,
				pairCounter: pairCounter{},
				bottomPairs: pairCounter{},
			})
		}
	}

	for p, r := range rules {
		recursivePairCounter(
			p,
			r,
			rules,
			pairCache_,
			0,
			cacheSize+1,
		)
	}

	for p, cache := range pairCache_ {
		pairCountSum := pairCounter{}
		for i, pairGraph_ := range cache {
			// add bottomPairs to pairCountSum
			for bottomPair, cnt := range pairGraph_.bottomPairs {
				pairCountSum[bottomPair] += cnt
			}
			// copy pairCountSum to this afterStep and move on
			pairCache_[p][i].pairCounter = *copyPairCounter(pairCountSum)
		}
	}
	return pairCache_
}

func recursivePairCounter(
	pair_ pair,
	rule rule,
	rules map[pair]rule,
	pairCache pairCache,
	currentStep int,
	maxSteps int,
) {
	if currentStep == maxSteps-1 {
		return
	}

	for _, p := range rule.nextPairs {
		pairCache[pair_][currentStep].bottomPairs[pair(p)]++
		applicableRule := rules[pair(p)]
		recursivePairCounter(
			pair_,
			applicableRule,
			rules,
			pairCache,
			currentStep+1,
			maxSteps,
		)
	}
}

func copyPairCounter(pairCountSum pairCounter) *pairCounter {
	pc := pairCounter{}
	for p, cnt := range pairCountSum {
		pc[p] = cnt
	}

	return &pc
}

const cacheSize = 20

func main() {
	start := time.Now()
	file, err := os.Open("input")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	startLetterCounter := map[string]int{}
	template := "KKOSPHCNOCHHHSPOBKVF"
	for i := 0; i < len(template); i++ {
		char := template[i]
		startLetterCounter[string(char)]++
	}

	rules := map[pair]rule{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pairId, insertionEl := lineToPairInsertionRule(scanner.Text())
		rules[pair(pairId)] = insertionEl
	}
	cache := createPairCache(cacheSize, rules)
	fmt.Println("cache ready", time.Since(start))

	// now create the initial pairs from template
	startPairs := templateToPairCount(template)
	var totalPairsCounter []pairCounter
	totalPairsCounter = append(totalPairsCounter, startPairs)
	desiredAfterStepNr := 39
	countPairsFromCache(startPairs, &totalPairsCounter, cache, desiredAfterStepNr, 1)

	ttCounter := mergePairCounters(totalPairsCounter...)
	printLetterCountFromPairCounter(ttCounter, rules, startLetterCounter)
	fmt.Println("End", time.Since(start))
}

func countPairsFromCache(
	pairs map[pair]int,
	totalPairsCounter *[]pairCounter,
	cache pairCache,
	stepsToGo int,
	mul int,
) {
	if stepsToGo == 0 {
		return
	}
	curStep := cacheSize
	remainingStepsToGo := 0

	if stepsToGo < cacheSize {
		curStep = stepsToGo
	}

	remainingStepsToGo = stepsToGo - curStep

	for pair_, cnt := range pairs {
		// first add this to the total
		pc := copyPairCounter(cache[pair_][curStep-1].pairCounter)
		for p, pcnt := range *pc {
			(*pc)[p] = pcnt * cnt * mul
		}
		*totalPairsCounter = append(*totalPairsCounter, *pc)
		countPairsFromCache(cache[pair_][curStep-1].bottomPairs, totalPairsCounter, cache, remainingStepsToGo, cnt)
	}
}

func printLetterCountFromPairCounter(totalPairCntr pairCounter, rules map[pair]rule, startLetterCounter map[string]int) {
	letterCounter := map[string]int{}
	for letter, cnt := range startLetterCounter {
		letterCounter[letter] = cnt
	}

	for p, cnt := range totalPairCntr {
		letterCounter[rules[p].inserts] += cnt
	}
	var letterCounts []int
	for _, cnt := range letterCounter {
		letterCounts = append(letterCounts, cnt)
	}
	sort.Ints(letterCounts)
	fmt.Println("Answer", letterCounts[len(letterCounts)-1]-letterCounts[0])
}

func mergePairCounters(counters ...pairCounter) pairCounter {
	res := pairCounter{}
	for _, counter := range counters {
		for p, cnt := range counter {
			res[p] += cnt
		}
	}

	return res
}

func lineToPairInsertionRule(lineAsText string) (string, rule) {
	split := strings.Split(lineAsText, " -> ")

	newPairs := []string{
		string(split[0][0]) + split[1],
		split[1] + string(split[0][1]),
	}

	return split[0], rule{split[1], newPairs}
}

func templateToPairCount(template string) pairCounter {
	var result []pair
	totalPairs := len(template) - 1

	for i := 0; i < totalPairs; i++ {
		p := pair(string(template[i]) + string(template[i+1]))
		result = append(result, p)
	}

	pc := pairCounter{}
	for _, pair_ := range result {
		pc[pair_] += 1
	}

	return pc
}
