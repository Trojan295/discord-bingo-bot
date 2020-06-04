package discord

import (
	"fmt"
	"strings"

	"github.com/Trojan295/discord-bingo-bot/pkg/bingo/game"
	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
)

func boardToString(board *game.Board) string {
	tableString := strings.Builder{}
	table := tablewriter.NewWriter(&tableString)
	table.SetRowLine(true)
	table.SetColWidth(20)

	for y, v := range board.FieldTexts {
		row := make([]string, 5)
		copy(row, v)

		for x := 0; x < 5; x++ {
			row[x] = fmt.Sprintf("%d. ", y*5+x+1) + row[x]
			if board.Ticked[y][x] == 1 {
				row[x] = "██████" + row[x] + "██████"
			}
		}

		table.Append(row)
	}

	table.Render()
	return tableString.String()
}

func getBoardMessages(board *game.Board) []string {
	messages := make([]string, 0)

	boardString := boardToString(board)
	lines := strings.Split(boardString, "\n")
	buffer := strings.Builder{}

	for _, line := range lines {
		if buffer.Len()+len(line) > 2000 {
			messages = append(messages, "`"+buffer.String()+"`")
			buffer.Reset()
		}

		buffer.WriteString(line)
		buffer.WriteRune('\n')
	}

	messages = append(messages, "`"+buffer.String()+"`")
	return messages
}

func PrintMessage(v interface{}, m *discordgo.MessageCreate) []string {
	switch msg := v.(type) {
	case game.TextMessage:
		return []string{msg.Text}

	case game.ListTextsResponse:
		builder := strings.Builder{}
		builder.WriteString("**Texts:**\n")

		for idx, text := range msg.Texts {
			builder.WriteString(fmt.Sprintf("%d. %s\n", idx+1, text))
		}
		return []string{builder.String()}

	case game.ShowBoardResponse:
		return getBoardMessages(msg.Board)

	case game.BingoWonResponse:
		boardMsgs := getBoardMessages(msg.Board)
		return append([]string{fmt.Sprintf("**Yeah! %s, you won!**", m.Author.Mention())}, boardMsgs...)

	case game.HelpResponse:
		builder := strings.Builder{}
		builder.WriteString("**Help:**\n")
		for _, command := range msg.Commands {
			builder.WriteString(fmt.Sprintf("`%s` - %s\n", command.Cmd, command.Description))
		}
		return []string{builder.String()}
	}

	return []string{}
}
