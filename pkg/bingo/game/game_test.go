package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGame_AddText(t *testing.T) {
	game := NewGame("")

	game.AddText("hello")
	assert.Equal(t, []string{"hello"}, game.Texts)
}

func TestGame_RemoveText(t *testing.T) {
	game := NewGame("")
	game.AddText("hello")
	game.AddText("world")
	game.AddText("test")

	game.RemoveText(1)
	assert.Equal(t, []string{"hello", "test"}, game.Texts)
}
