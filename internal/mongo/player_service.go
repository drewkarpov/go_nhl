package mongo

import (
	"context"
	m "github.com/drewkarpov/go_nhl/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type MongoPlayerService struct {
	Collection *mongo.Collection
}

func (d MongoPlayerService) Init() MongoPlayerService {
	ctx, _ := context.WithTimeout(context.Background(), 40*time.Second)
	var cred options.Credential = options.Credential{Username: os.Getenv("MONGODB_USER"),
		Password: os.Getenv("MONGODB_PWD")}
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongo:27017").SetAuth(cred))
	d.Collection = client.Database("nhl").Collection("players")
	return MongoPlayerService{Collection: d.Collection}
}

func (d MongoPlayerService) WritePlayer(playerDTO m.PlayerDTO) *mongo.InsertOneResult {
	player := playerDTO.MapToPlayer()
	player.ID = primitive.NewObjectID()
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	result, err := d.Collection.InsertOne(ctx, player)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func (d MongoPlayerService) GetAllPlayers() ([]m.Player, error) {
	var players []m.Player
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	cursor, err := d.Collection.Find(ctx, bson.M{})

	if cursor != nil {
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var person m.Player
			cursor.Decode(&person)
			players = append(players, person)
		}
		return players, err
	}
	return players, err
}

func (d MongoPlayerService) ChangePlayerData(id primitive.ObjectID, playerDTO m.PlayerDTO) (m.Player, error) {
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	player := playerDTO.MapToPlayer()

	_, err := d.Collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{{"name", player.Name}}},
			{"$set", bson.D{{"status", player.Status}}},
			{"$set", bson.D{{"priority", player.Priority}}},
			{"$set", bson.D{{"comment", player.Comment}}},
		},
	)
	return player, err
}

func (d MongoPlayerService) GetPlayerById(id primitive.ObjectID) (m.Player, error) {
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	var player m.Player
	err := d.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&player)
	return player, err
}

func (d MongoPlayerService) DeletePlayer(id primitive.ObjectID) (string, error) {
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	_, err := d.Collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return "failure", err
	}
	return "success", err
}

func (d MongoPlayerService) AddGameToPlayer(id primitive.ObjectID, game m.Game) (m.Game, error) {
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	game.ID = primitive.NewObjectID()

	_, err := d.Collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$push": bson.M{"games": game}},
	)
	return game, err
}
