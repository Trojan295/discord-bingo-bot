package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Trojan295/discord-bingo-bot/pkg/bingo/game"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/sirupsen/logrus"
)

var (
	ttl = 30 * 24 * time.Hour
)

type DynamoDBGameRepository struct {
	dynamoDB  *dynamodb.DynamoDB
	tableName string
}

func NewDynamoDBGameRepository(sess *session.Session, tableName string) (*DynamoDBGameRepository, error) {
	return &DynamoDBGameRepository{
		dynamoDB:  dynamodb.New(sess),
		tableName: tableName,
	}, nil
}

func (repo *DynamoDBGameRepository) Get(ID string) (*game.Game, error) {
	resp, err := repo.dynamoDB.GetItem(&dynamodb.GetItemInput{
		TableName: &repo.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"GameId": {
				S: &ID,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if resp.Item == nil {
		return nil, nil
	}

	encodedData := resp.Item["Data"].B

	logrus.Debugf("read game %s data: %v", ID, string(encodedData))

	var g game.Game
	if err := json.Unmarshal(encodedData, &g); err != nil {
		return nil, err
	}

	return &g, nil
}

func (repo *DynamoDBGameRepository) Persist(g *game.Game) error {
	encodedGame, err := json.Marshal(g)
	if err != nil {
		return err
	}

	ttlAttr := time.Now().Add(ttl)

	now := fmt.Sprintf("%d", ttlAttr.Unix())

	item := map[string]*dynamodb.AttributeValue{
		"GameId": {
			S: &g.ID,
		},
		"Data": {
			B: encodedGame,
		},
		"UpdatedAt": {
			N: &now,
		},
	}

	_, err = repo.dynamoDB.PutItem(&dynamodb.PutItemInput{
		TableName: &repo.tableName,
		Item:      item,
	})
	return err
}
