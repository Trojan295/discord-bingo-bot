package game

import (
	"fmt"
)

type Board struct {
	FieldTexts [][]string
	Ticked     [][]uint8
}

func NewBoard(fields [][]string) (*Board, error) {
	board := Board{
		FieldTexts: fields,
		Ticked: [][]uint8{
			make([]uint8, 5),
			make([]uint8, 5),
			make([]uint8, 5),
			make([]uint8, 5),
			make([]uint8, 5),
		},
	}

	return &board, nil
}

func (board *Board) IsBingo() bool {
	for x := 0; x < 5; x++ {
		var sum uint8 = 0
		for y := 0; y < 5; y++ {
			sum += board.Ticked[y][x]
		}

		if sum == 5 {
			return true
		}
	}

	for y := 0; y < 5; y++ {
		var sum uint8 = 0
		for x := 0; x < 5; x++ {
			sum += board.Ticked[y][x]
		}

		if sum == 5 {
			return true
		}
	}

	if board.Ticked[0][0]+board.Ticked[1][1]+board.Ticked[2][2]+board.Ticked[3][3]+board.Ticked[4][4] == 5 {
		return true
	}
	if board.Ticked[0][4]+board.Ticked[1][3]+board.Ticked[2][2]+board.Ticked[3][1]+board.Ticked[4][1] == 5 {
		return true
	}

	return false
}

func (board *Board) Mark(pos int) (bool, error) {
	if pos < 0 || pos > 24 {
		return false, fmt.Errorf("wrong field number")
	}

	y := pos / 5
	x := pos % 5

	if board.Ticked[y][x] == 1 {
		return false, nil
	}

	board.Ticked[y][x] = 1
	return true, nil
}
