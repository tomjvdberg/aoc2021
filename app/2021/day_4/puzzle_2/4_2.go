package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Board struct {
	id        int
	positions []BoardPosition
}

func (board *Board) addPosition(position BoardPosition) {
	board.positions = append(board.positions, position)
}

func (board *Board) markNumber(numberToBeMarked int) {
	for i := 0; i < len(board.positions); i++ {
		position := &board.positions[i]
		if (*position).number == numberToBeMarked {
			(*position).isMarked = true
		}
	}
}

func (board *Board) isComplete() bool {
	isComplete := false
	// slice each row and check if it is complete
	var rows []BoardRow
	for _, rowNr := range [5]int{1, 2, 3, 4, 5} {
		boardRow, err := board.getRow(rowNr)
		if err != nil {
			continue // skip the errors
		}
		rows = append(rows, boardRow)
	}
	// then slice each column and check if it is complete
	var columns []BoardColumn
	for _, columnNr := range [5]int{1, 2, 3, 4, 5} {
		boardColumn, err := board.getColumn(columnNr)
		if err != nil {
			continue // skip the errors
		}
		columns = append(columns, boardColumn)
	}

	for _, row := range rows {
		if row.isComplete() {
			isComplete = true
		}
	}

	for _, column := range columns {
		if column.isComplete() {
			isComplete = true
		}
	}

	return isComplete
}

func (board *Board) getRow(rowNr int) (BoardRow, error) {
	var positions []*BoardPosition
	for i := 0; i < len(board.positions); i++ { // for each position
		position := &board.positions[i]
		if (*position).row == rowNr {
			positions = append(positions, position)
		}
	}

	if len(positions) == 0 {
		return BoardRow{}, RowNotAvailable{}
	}

	if len(positions) < 5 {
		return BoardRow{}, InvalidRowLength{}
	}

	return BoardRow{positions}, nil
}

func (board *Board) getColumn(columnNr int) (BoardColumn, error) {
	var positions []*BoardPosition
	for i := 0; i < len(board.positions); i++ { // for each position
		position := &board.positions[i]
		if (*position).column == columnNr {
			positions = append(positions, position)
		}
	}

	if len(positions) == 0 {
		return BoardColumn{}, ColumnNotAvailable{}
	}

	if len(positions) < 5 {
		return BoardColumn{}, InvalidColumnLength{}
	}

	return BoardColumn{positions}, nil
}

type BoardPosition struct {
	row      int
	column   int
	number   int
	isMarked bool
}

type BoardRow struct {
	positions []*BoardPosition
}

func (row *BoardRow) isComplete() bool {
	isComplete := true
	for i := 0; i < len(row.positions); i++ {
		position := &row.positions[i]
		if !(*position).isMarked {
			isComplete = false
		}
	}

	return isComplete
}

type InvalidRowLength struct{}

func (InvalidRowLength) Error() string { return "invalid row length encountered" }

type RowNotAvailable struct{}

func (RowNotAvailable) Error() string { return "row not available" }

type BoardColumn struct {
	positions []*BoardPosition
}

func (column *BoardColumn) isComplete() bool {
	isComplete := true
	for i := 0; i < len(column.positions); i++ {
		position := &column.positions[i]
		if !(*position).isMarked {
			isComplete = false
		}
	}

	return isComplete
}

type InvalidColumnLength struct{}

func (InvalidColumnLength) Error() string { return "invalid column length encountered" }

type ColumnNotAvailable struct{}

func (ColumnNotAvailable) Error() string { return "column not available" }

func main() {
	file, err := os.Open("input")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	var numbersDrawn []int
	var boards []Board

	rowCounter := 1 // start with row 1
	bufferBoard := Board{1, []BoardPosition{}}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Read the first line that will be the number drawer
		if numbersDrawn == nil {
			stringNumbersDrawn := strings.Split(scanner.Text(), ",")
			for _, stringNumberDrawn := range stringNumbersDrawn {
				numberDrawn, _ := strconv.Atoi(stringNumberDrawn)
				numbersDrawn = append(numbersDrawn, numberDrawn)
			}
			continue
		}

		// Read the other input lines. Ignore empty lines.
		if len(scanner.Bytes()) < 1 {
			continue
		}

		// now starts the reading of the boards
		for i, number := range extractNumbersFromLine(scanner.Text()) {
			bufferBoard.addPosition(BoardPosition{
				row:      rowCounter,
				column:   i + 1,
				number:   number,
				isMarked: false,
			})
		}

		rowCounter++

		if rowCounter > 5 {
			// board is full
			boards = append(boards, bufferBoard)
			rowCounter = 1 // reset row counter
			bufferBoard.positions = []BoardPosition{}
			bufferBoard.id = len(boards) + 1
		}
	}

	var lastDrawnNumber int
	lastWinningBoardFound := false

	var lastWinningBoard Board
	for _, numberDrawn := range numbersDrawn {
		if lastWinningBoardFound {
			break
		}
		var boardIdsToBeRemoved []int

		lastDrawnNumber = numberDrawn
		boardsCount := len(boards)
		for i := 0; i < boardsCount; i++ {
			board := &boards[i]
			(*board).markNumber(numberDrawn)

			if board.isComplete() {
				if len(boards) == 1 {
					// we found the last winning board
					lastWinningBoardFound = true
					lastWinningBoard = *board
					break
				}
				boardIdsToBeRemoved = append(boardIdsToBeRemoved, board.id)
			}
		}

		// clean up completed boards
		for _, boardId := range boardIdsToBeRemoved {
			boards = removeBoardById(boards, boardId)
		}
	}

	fmt.Println("boards", boards)
	fmt.Println("lastWinningBoard", lastWinningBoard)
	fmt.Println("lastDrawnNumber", lastDrawnNumber)
	unMarkedNumberSum := 0

	for _, position := range lastWinningBoard.positions {
		if position.isMarked {
			continue
		}
		unMarkedNumberSum += position.number
	}

	fmt.Println("unMarkedNumberSum", unMarkedNumberSum)
	fmt.Println("score", unMarkedNumberSum*lastDrawnNumber)
}

func removeBoardById(boards []Board, id int) []Board {
	var indexFound bool
	var indexOfBoardWithId int

	for i := 0; i < len(boards); i++ {
		if boards[i].id == id {
			indexOfBoardWithId = i
			indexFound = true
		}
	}

	if !indexFound {
		fmt.Println("No board found for id", id)

		return boards
	}

	boards[indexOfBoardWithId] = boards[len(boards)-1]

	return boards[:len(boards)-1]
}

func extractNumbersFromLine(line string) []int {
	numbersOnRowAsString := strings.Split(line, " ")

	var sanitizedOutput []int
	for _, numberAsString := range numbersOnRowAsString {
		if numberAsString == "" {
			continue
		}
		number, _ := strconv.Atoi(numberAsString)
		sanitizedOutput = append(sanitizedOutput, number)
	}

	return sanitizedOutput
}
