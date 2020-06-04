package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func dummyBoard() *Board {
	board, _ := NewBoard(
		[][]string{
			make([]string, 5),
			make([]string, 5),
			make([]string, 5),
			make([]string, 5),
			make([]string, 5),
		},
	)
	return board
}

func TestBoard_MarkField(t *testing.T) {
	board := dummyBoard()

	assert.Equal(t, uint8(0), board.Ticked[0][0])
	board.Mark(1)
	assert.Equal(t, uint8(1), board.Ticked[0][0])

	assert.Equal(t, uint8(0), board.Ticked[0][2])
	board.Mark(3)
	assert.Equal(t, uint8(1), board.Ticked[0][2])

	assert.Equal(t, uint8(0), board.Ticked[2][0])
	board.Mark(11)
	assert.Equal(t, uint8(1), board.Ticked[2][0])
}

func TestBoard_MarkingMarkedField(t *testing.T) {
	board := dummyBoard()

	marked, _ := board.Mark(1)
	assert.Equal(t, true, marked)
	marked, _ = board.Mark(1)
	assert.Equal(t, false, marked)
}

func TestBoard_MarkWrongFields(t *testing.T) {
	board := dummyBoard()

	_, err := board.Mark(0)
	assert.NotNil(t, err)

	_, err = board.Mark(26)
	assert.NotNil(t, err)
}

func TestBoard_Bingo(t *testing.T) {
	board := dummyBoard()

	assert.False(t, board.IsBingo())
	board.Mark(1)
	board.Mark(7)
	board.Mark(13)
	board.Mark(19)
	board.Mark(25)
	assert.True(t, board.IsBingo())
}
