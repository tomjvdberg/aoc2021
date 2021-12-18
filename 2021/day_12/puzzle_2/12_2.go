package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"
)

type caveName string

type Cave struct {
	name        caveName
	isBigCave   bool
	connections map[caveName]Cave
}

type Connection struct {
	enpointA caveName
	enpointB caveName
}

type Hop struct {
	from       Cave
	to         Cave
	nextHops   []Hop
}

func main() {
	start := time.Now()
	file, err := os.Open("input")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	caves := map[caveName]Cave{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		route := lineToRoute(scanner.Text())
		// make sure the caves exist
		for _, connectedCave := range []caveName{route.enpointA, route.enpointB} {
			_, exists := caves[connectedCave]
			if !exists {
				caves[connectedCave] = Cave{
					connectedCave,
					isBigCave(connectedCave),
					map[caveName]Cave{}}
				continue
			}
		}
		caves[route.enpointA].connections[route.enpointB] = caves[route.enpointB]
		caves[route.enpointB].connections[route.enpointA] = caves[route.enpointA]
	}

	var hops []Hop
	for _, connectedCave := range caves["start"].connections {
		hop := Hop{
			caves["start"],
			connectedCave,
			[]Hop{},
		}
		hop.nextHops = nextHopsFromHops([]Hop{hop})

		hops = append(hops, hop)
	}

	// transform hops to paths
	var paths []string
	for _, startHop := range hops {
		p := string(startHop.from.name)
		handlePaths(startHop.nextHops, p, &paths)
	}
	fmt.Println("paths", len(paths))
	fmt.Println("End", time.Since(start))
}

func handlePaths(hops []Hop, p string, paths *[]string) {
	for _, hop := range hops {
		hopP := p + "," + string(hop.from.name)

		if len(hop.nextHops) == 0 {
			if hop.to.name == "end" {
				//fmt.Println(hopP + "," + string(hop.to.name))
				*paths = append(*paths, hopP+","+string(hop.to.name))
			}
			continue
		}
		handlePaths(hop.nextHops, hopP, paths)
	}
}

func nextHopsFromHops(previousHops []Hop) []Hop {
	lastHop := previousHops[len(previousHops)-1]

	var nextHops []Hop
	currentCave := lastHop.to

	for _, destCave := range lastHop.to.connections {
		if destCave.name == "start" {
			continue
		}

		if !isAllowedToVisitCave(destCave, previousHops) {
			continue
		}
		if currentCave.name == "end"{
			// make no hops anymore
			continue
		}

		hop := Hop{
			from:       currentCave,
			to:         destCave,
			nextHops:   []Hop{},
		}
		// MAKE THE HOP
		hop.nextHops = nextHopsFromHops(append(previousHops,hop))
		nextHops = append(nextHops, hop)
	}

	return nextHops
}

func isAllowedToVisitCave(destCave Cave, previousHops []Hop) bool {
	// you are allowed to visit when
	// dest is not start
	if destCave.name == "start" {
		return false
	}
	// dest is a big cave
	if destCave.isBigCave {
		return true
	}
	// dest is end
	if destCave.name == "end" {
		return true
	}

	smallCavesVisited := map[string]int{}
	for _, previousHop := range previousHops {
		if previousHop.to.isBigCave {
			continue
		}
		_, exists := smallCavesVisited[string(previousHop.to.name)]
		if !exists {
			smallCavesVisited[string(previousHop.to.name)] = 1
		} else {
			smallCavesVisited[string(previousHop.to.name)]++
		}
	}

	_, exists := smallCavesVisited[string(destCave.name)]
	if !exists {
		// dest is small and not visited yet
		return true
	}

	smallCaveVisitedTwice := false
	for _, visitCount := range smallCavesVisited {
		if visitCount > 1 {
			smallCaveVisitedTwice = true
		}
	}

	if smallCaveVisitedTwice {
		return false
	}

	return true
}

func isBigCave(name caveName) bool {
	for _, r := range name {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func lineToRoute(lineAsText string) Connection {
	caves := strings.Split(lineAsText, "-")

	return Connection{caveName(caves[0]), caveName(caves[1])}
}
