package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Trojan295/discord-bingo-bot/pkg/bingo/game"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoGameRepository struct {
	client *mongo.Client
}

func NewMongoGameRepository(URI string) (*MongoGameRepository, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		return nil, err
	}

	repo := &MongoGameRepository{
		client: client,
	}
	if err = repo.connect(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *MongoGameRepository) connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return r.client.Connect(ctx)
}

func serializeGame(g *game.Game) bson.M {
	var doc bson.M
	err := mapstructure.Decode(g, &doc)
	if err != nil {
		panic(err)
	}

	return doc
}

func (r *MongoGameRepository) Persist(g *game.Game) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := r.client.Database("bingo").Collection("games")
	doc := serializeGame(g)

	result := collection.FindOne(ctx, bson.D{{Key: "ID", Value: g.ID}})

	err := result.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			_, err2 := collection.InsertOne(ctx, doc)
			return err2
		}

		return err
	}

	_, err = collection.UpdateOne(ctx, bson.D{{Key: "ID", Value: g.ID}}, bson.M{
		"$set": doc,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoGameRepository) Get(GameID string) (*game.Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := r.client.Database("bingo").Collection("games")

	result := collection.FindOne(ctx, bson.D{{Key: "ID", Value: GameID}})
	err := result.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	var g game.Game
	result.Decode(&g)
	return &g, nil
}
