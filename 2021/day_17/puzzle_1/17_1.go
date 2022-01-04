package main

import (
	"fmt"
	"time"
)

type coord struct {
	x int
	y int
}

func main() {
	start := time.Now()
	targetArea := [2]coord{{211, -69}, {232, -124}}
	cnt := 0
	hits := []coord{}
	stepPos := coord{}
	initialVelocity := coord{}
	highestY := 0
	for {
		if cnt > 10100000 {
			break
		}

		if initialVelocity.y > 10000 {
			initialVelocity.y = 0
			initialVelocity.x++
		}

		step := 0
		highY := 0
		for {
			red := initialVelocity.x - step
			if red < 0 {
				red = 0
			}
			stepPos.x += red

			stepPos.y += initialVelocity.y
			stepPos.y -= step

			if stepPos.y > highY {
				highY = stepPos.y
			}

			if isOverShoot(stepPos, targetArea) {
				break
			}

			if isInTargetArea(stepPos, targetArea) {
				if highY > highestY {
					highestY = highY
				}
				hits = append(hits, coord{initialVelocity.x, initialVelocity.y})
			}
			step++
		}

		stepPos.x = 0
		stepPos.y = 0
		initialVelocity.y++

		cnt++
	}
	fmt.Println(highestY)
	fmt.Println("End", time.Since(start))
}

func isInTargetArea(c coord, targetArea [2]coord) bool {
	return c.x >= targetArea[0].x && c.x <= targetArea[1].x && c.y <= targetArea[0].y && c.y >= targetArea[1].y
}
func isOverShoot(c coord, targetArea [2]coord) bool {
	isTooLow := c.y < targetArea[1].y
	isTooFar := c.x > targetArea[1].x

	return isTooLow || isTooFar
}
