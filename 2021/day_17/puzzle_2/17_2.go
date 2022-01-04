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
	hits := map[coord]bool{}

	for _, yDir := range []bool{true, false} {
		cnt := 0
		stepPos := coord{}
		initialVelocity := coord{}
		for {
			if initialVelocity.x > targetArea[1].x {
				break
			}

			if initialVelocity.y > 10000 || initialVelocity.y < -10000 {
				initialVelocity.y = 0
				initialVelocity.x++
			}

			step := 0
			for {
				red := initialVelocity.x - step
				if red < 0 {
					red = 0
				}
				stepPos.x += red

				stepPos.y += initialVelocity.y
				stepPos.y -= step

				if isOverShoot(stepPos, targetArea) {
					break
				}

				if isInTargetArea(stepPos, targetArea) {
					hits[coord{initialVelocity.x, initialVelocity.y}] = true
				}
				step++
			}

			stepPos.x = 0
			stepPos.y = 0

			if yDir {
				initialVelocity.y++
			} else {
				initialVelocity.y--
			}

			cnt++
		}
	}

	fmt.Println(len(hits))
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
