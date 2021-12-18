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
	allowReUse bool
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
			false,
			[]Hop{},
		}
		hop.nextHops = nextHopsFromHops(hop, map[string]Hop{})

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
		hopP := p + "-" + string(hop.from.name)

		if len(hop.nextHops) == 0 {
			if hop.to.name == "end" {
				*paths = append(*paths, hopP+"-"+string(hop.to.name))
			}
			continue
		}
		handlePaths(hop.nextHops, hopP, paths)
	}
}

func nextHopsFromHops(currentHop Hop, previousHops map[string]Hop) []Hop {
	copiedPreviousHops := *copyHops(previousHops)
	var nextHops []Hop
	previousCave := currentHop.from
	currentCave := currentHop.to

	for _, connectedCave := range currentHop.to.connections {
		if connectedCave.name == "start" {
			continue
		}

		if connectedCave.name == previousCave.name && !previousCave.isBigCave {
			continue
		}

		hop := Hop{
			from:       currentCave,
			to:         connectedCave,
			allowReUse: connectedCave.isBigCave,
			nextHops:   []Hop{},
		}
		// check if hop already exists.
		_, hopAlreadyExists := copiedPreviousHops[hopId(hop)]
		if hopAlreadyExists {
			continue
		}

		if !isAllowedToRevisitCave(connectedCave, previousHops) {
			continue
		}

		// add current hop to previous hops
		copiedPreviousHops[hopId(currentHop)] = currentHop
		
		if hop.to.name != "end" {
			hop.nextHops = nextHopsFromHops(hop, copiedPreviousHops)
		}

		nextHops = append(nextHops, hop)
	}

	return nextHops
}

func isAllowedToRevisitCave(connectedCave Cave, previousHops map[string]Hop) bool {
	for _, previousHop := range previousHops {
		if previousHop.to.name == connectedCave.name && !connectedCave.isBigCave {
			return false
		}
	}

	return true
}

func copyHops(m map[string]Hop) *map[string]Hop {
	cm := map[string]Hop{}
	for id, e := range m {
		cm[id] = e
	}

	return &cm
}
func hopId(hop Hop) string {
	return string(hop.from.name + "-" + hop.to.name)
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
