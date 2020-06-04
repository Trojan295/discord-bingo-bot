package game

import (
	"fmt"
	"regexp"
	"strconv"
)

type TextMessage struct {
	Text string
}

type ListTextsResponse struct {
	Texts []string
}

type ShowBoardResponse struct {
	Board *Board
}

type BingoWonResponse struct {
	Board *Board
}

type GameRepository interface {
	Persist(*Game) error
	Get(ID string) (*Game, error)
}

type Controller struct {
	gamesRepository GameRepository
}

func NewController(repo GameRepository) *Controller {
	return &Controller{
		gamesRepository: repo,
	}
}

func (ctrl *Controller) ProcessMessage(gameID, message string) (interface{}, error) {
	game, err := ctrl.gamesRepository.Get(gameID)
	if err != nil {
		return nil, err
	}

	if game == nil {
		game = NewGame(gameID)
		if err := ctrl.gamesRepository.Persist(game); err != nil {
			return nil, err
		}
	}

	unknownMessage := TextMessage{
		Text: "Emm.. Could you repeat, please?",
	}

	if message == ".bingo list" {
		return ctrl.listTexts(game)
	}

	if message == ".bingo clear" {
		return ctrl.clearTexts(game)
	}

	addTextRegex := regexp.MustCompile(`^.bingo add (.+)$`)
	if match := addTextRegex.FindStringSubmatch(message); len(match) > 0 {
		return ctrl.addText(game, match[1])
	}

	removeTextRegex := regexp.MustCompile(`^.bingo remove (\d+)$`)
	if match := removeTextRegex.FindStringSubmatch(message); len(match) > 0 {
		pos, err := strconv.Atoi(match[1])
		if err != nil {
			return unknownMessage, err
		}

		return ctrl.removeText(game, pos-1)
	}

	if message == ".bingo new" {
		return ctrl.newBoard(game)
	}

	if message == ".bingo show" {
		if game.Board == nil {
			return TextMessage{
				Text: "You don't have started. Type .bingo new",
			}, nil
		}

		return ShowBoardResponse{
			Board: game.Board,
		}, nil
	}

	markFieldRegex := regexp.MustCompile(`^.bingo mark (.+)$`)
	if match := markFieldRegex.FindStringSubmatch(message); len(match) > 0 {
		pos, err := strconv.Atoi(match[1])
		if err != nil {
			return unknownMessage, err
		}

		return ctrl.markField(game, pos-1)
	}

	return unknownMessage, nil
}

func (ctrl *Controller) listTexts(game *Game) (interface{}, error) {
	return ListTextsResponse{
		Texts: game.Texts,
	}, nil
}

func (ctrl *Controller) clearTexts(game *Game) (interface{}, error) {
	game.Texts = make([]string, 0)

	if err := ctrl.gamesRepository.Persist(game); err != nil {
		return nil, err
	}

	return ListTextsResponse{
		Texts: game.Texts,
	}, nil
}

func (ctrl *Controller) addText(game *Game, text string) (interface{}, error) {
	if err := game.AddText(text); err != nil {
		return TextMessage{
			Text: fmt.Sprintf("Ops.. %s", err.Error()),
		}, nil
	}

	if err := ctrl.gamesRepository.Persist(game); err != nil {
		return nil, err
	}

	return TextMessage{
		Text: fmt.Sprintf("Added: %s", text),
	}, nil
}

func (ctrl *Controller) removeText(game *Game, pos int) (interface{}, error) {
	if err := game.RemoveText(pos); err != nil {
		return nil, err
	}

	if err := ctrl.gamesRepository.Persist(game); err != nil {
		return nil, err
	}

	return TextMessage{
		Text: "Text removed",
	}, nil
}

func (ctrl *Controller) newBoard(game *Game) (interface{}, error) {
	if err := game.NewBoard(); err != nil {
		return nil, err
	}

	if err := ctrl.gamesRepository.Persist(game); err != nil {
		return nil, err
	}

	return ShowBoardResponse{
		Board: game.Board,
	}, nil
}

func (ctrl *Controller) markField(game *Game, pos int) (interface{}, error) {
	if game.Board.IsBingo() {
		return TextMessage{
			Text: "Ops.. this board is already finished...",
		}, nil
	}

	if marked, err := game.Board.Mark(pos); err != nil {
		return TextMessage{
			Text: "Wrong field number",
		}, nil

	} else if !marked {
		return TextMessage{
			Text: "Field already marked",
		}, nil
	}

	if err := ctrl.gamesRepository.Persist(game); err != nil {
		return nil, err
	}

	if game.Board.IsBingo() {
		return BingoWonResponse{
			Board: game.Board,
		}, nil
	}

	return ShowBoardResponse{
		Board: game.Board,
	}, nil
}
