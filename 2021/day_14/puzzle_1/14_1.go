package main

import (
	"bufio"
	"fmt"
	"os"
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

	template := "KKOSPHCNOCHHHSPOBKVF"
	rules := map[string]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pairId, insertionEl := lineToPairInsertionRule(scanner.Text())
		rules[pairId] = insertionEl
	}

	desiredSteps := 10
	nextTemplate := ""
	for i := 0; i < desiredSteps; i++ {
		pairs := templateToPairs(nextTemplate)
		if i == 0 {
			pairs = templateToPairs(template)
		}

		n := ""
		for _, pair := range pairs {
			n = appendToNextTemplate(n, ruleToResult(pair, rules[pair]))
		}
		nextTemplate = n
	}

	elementCount := map[string]int{}
	countElements(nextTemplate, elementCount)
	firstElPassed := false
	mcEl := 0
	lcEl := 0
	for _, cnt := range elementCount{
		if !firstElPassed {
			mcEl, lcEl = cnt, cnt
			firstElPassed = true
		}

		if cnt > mcEl{
			mcEl = cnt
		}
		if cnt < lcEl{
			lcEl = cnt
		}
	}

	fmt.Println("answer", mcEl - lcEl)
	fmt.Println("End", time.Since(start))
}

func countElements(template string, count map[string]int) {
	elLen := len(template)

	for i:=0; i< elLen; i++ {
		el := string(template[i])
		_, ex := count[el]
		if ex {
			count[el]++
			continue
		}
		count[el] = 1
	}
}

func appendToNextTemplate(s string, appendable string) string {
	if (len(s)) == 0 {
		// template is empty
		return appendable
	}
	ns := s[:len(s)-1]
	res := ns + appendable

	return res
}

func ruleToResult(pair string, insertion string) string {
	return string(pair[0]) + insertion + string(pair[1])
}

func templateToPairs(template string) []string {
	var result []string
	totalPairs := len(template) - 1

	for i := 0; i < totalPairs; i++ {
		result = append(result, string(template[i])+string(template[i+1]))
	}

	return result
}

func lineToPairInsertionRule(lineAsText string) (string, string) {
	split := strings.Split(lineAsText, " -> ")

	return split[0], split[1]
}
