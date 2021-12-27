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

func copyPairCounter(pairCountSum pairCounter) *pairCounter {
	pc := pairCounter{}
	for p, cnt := range pairCountSum {
		pc[p] = cnt
	}

	return &pc
}

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

	totalPairCounter := pairCounter{}
	bottomPairsCounter := templateToPairCount(template)
	for i := 0; i < 40; i++ {
		nextBottomPairs := pairCounter{}
		for p, cnt := range bottomPairsCounter {
			totalPairCounter[p] += cnt
			for _, stringP := range rules[p].nextPairs {
				nextBottomPairs[pair(stringP)] += cnt
			}
		}
		bottomPairsCounter = *copyPairCounter(nextBottomPairs)
	}
	printLetterCountFromPairCounter(totalPairCounter, rules, startLetterCounter)
	fmt.Println("End", time.Since(start))
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
