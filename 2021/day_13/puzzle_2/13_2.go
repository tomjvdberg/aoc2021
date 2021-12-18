package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type dotCount int
type DotMap struct {
	rows [][]dotCount
}

type axis string
type FoldInstruction struct {
	axis axis
	line int
}

func (dotMap *DotMap) fold(instr FoldInstruction) {
	// set all dots on folding line to zero
	if instr.axis == "y" {
		for i := range dotMap.rows[instr.line] {
			dotMap.rows[instr.line][i] = 0
		}
	}
	if instr.axis == "x" {
		for iY := range dotMap.rows {
			dotMap.rows[iY][instr.line] = 0
		}
	}
	// folding line is removed

	// now copy higher lines values to lower lines
	for prev, next := instr.line-1, instr.line+1; prev > -1; prev, next = prev-1, next+1 {
		if instr.axis == "y" {
			if next > (len(dotMap.rows) - 1) {
				continue
			}
			for i, cnt := range dotMap.rows[next] {
				dotMap.rows[prev][i] += cnt
				dotMap.rows[next][i] = 0 // remove the dot on the higher line
			}
		}
		if instr.axis == "x" {
			for iR, _ := range dotMap.rows {
				dotMap.rows[iR][prev] += dotMap.rows[iR][next]
				dotMap.rows[iR][next] = 0 // remove the dot in the higher x line
			}
		}
	}
}

func main() {
	start := time.Now()
	dotMap := retrieveDots()
	foldInstructions := retrieveFoldInstructions()

	for _, instr := range foldInstructions {
		dotMap.fold(instr)
	}
	draw(dotMap)
	fmt.Println("End", time.Since(start))
}

func draw(dotMap DotMap) {
	fn := "./output/code.txt"

	f, err := os.Create(fn)
	if err != nil {
		fmt.Println(err)
	}
	// close the file with defer
	defer f.Close()

	for _, row := range dotMap.rows {
		output := ""
		for _, cnt := range row {
			if cnt > 0 {
				output += "#"
				continue
			}
			output += " "
		}
		f.WriteString(output + "\n")
	}
}

func retrieveDots() DotMap {
	dotFile, err := os.Open("dotinput")
	if err != nil {
		fmt.Println(err)
	}
	defer dotFile.Close()

	dotMap := DotMap{[][]dotCount{}}

	highestX := 0
	highestY := 0
	scanner := bufio.NewScanner(dotFile)
	for scanner.Scan() {
		x, y := LineToCoords(scanner.Text())
		if x > highestX {
			highestX = x
		}
		if y > highestY {
			highestY = y
		}
	}

	fmt.Println(highestX, highestY)
	// first add rows and cols
	for iY := 0; iY < highestY+1; iY++ {
		var xPoints []dotCount
		for iX := 0; iX < highestX+1; iX++ {
			xPoints = append(xPoints, dotCount(0))
		}
		dotMap.rows = append(dotMap.rows, xPoints)
	}

	dotFile2, err := os.Open("dotinput")
	if err != nil {
		fmt.Println(err)
	}
	defer dotFile2.Close()

	scanner = bufio.NewScanner(dotFile2)
	for scanner.Scan() {
		x, y := LineToCoords(scanner.Text())
		dotMap.rows[y][x]++
	}

	return dotMap
}

func retrieveFoldInstructions() []FoldInstruction {
	foldfile, err := os.Open("foldinput")
	if err != nil {
		fmt.Println(err)
	}
	defer foldfile.Close()

	var instructions []FoldInstruction

	scanner := bufio.NewScanner(foldfile)
	for scanner.Scan() {
		foldStripped := strings.Split(scanner.Text(), "fold along ")
		axisDistanceString := strings.Split(foldStripped[1], "=")
		ax := axis(axisDistanceString[0])
		line, _ := strconv.Atoi(axisDistanceString[1])
		instructions = append(instructions, FoldInstruction{ax, line})
	}

	return instructions
}

func LineToCoords(lineAsText string) (int, int) {
	coordsAsString := strings.Split(lineAsText, ",")

	x, _ := strconv.Atoi(coordsAsString[0])
	y, _ := strconv.Atoi(coordsAsString[1])

	return x, y
}
