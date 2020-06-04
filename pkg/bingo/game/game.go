package game

import (
	"fmt"
	"math/rand"
)

type Game struct {
	ID    string
	Texts []string
	Board *Board
}

func NewGame(ID string) *Game {
	return &Game{
		ID:    ID,
		Texts: make([]string, 0),
	}
}

func (g *Game) AddText(text string) error {
	if len(g.Texts) == 24 {
		return fmt.Errorf("there are already 24 texts")
	}

	g.Texts = append(g.Texts, text)
	return nil
}

func (g *Game) RemoveText(pos int) error {
	if len(g.Texts) <= pos {
		return fmt.Errorf("index overflow")
	}

	g.Texts = append(g.Texts[:pos], g.Texts[pos+1:]...)
	return nil
}

func (g *Game) NewBoard() error {
	if len(g.Texts) != 24 {
		return fmt.Errorf("you dont have 24 texts")
	}

	randomTexts := make([]string, len(g.Texts))
	copy(randomTexts, g.Texts)

	rand.Shuffle(len(randomTexts), func(i, j int) {
		randomTexts[i], randomTexts[j] = randomTexts[j], randomTexts[i]
	})

	pos := 0
	fields := make([][]string, 5)
	ticked := make([][]uint8, 5)
	for y := 0; y < 5; y++ {
		fields[y] = make([]string, 5)
		ticked[y] = make([]uint8, 5)

		for x := 0; x < 5; x++ {
			if x == 2 && y == 2 {
				fields[y][x] = "BINGO"
				ticked[y][x] = 1
				continue
			}
			fields[y][x] = randomTexts[pos]
			pos++
		}
	}

	g.Board = &Board{
		FieldTexts: fields,
		Ticked:     ticked,
	}
	return nil
}
