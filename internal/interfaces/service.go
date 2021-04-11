package interfaces

import (
	m "github.com/drewkarpov/go_nhl/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PlayerService interface {
	WritePlayer(playerDTO m.PlayerDTO) *mongo.InsertOneResult
	GetAllPlayers() ([]m.Player, error)
	ChangePlayerData(id primitive.ObjectID, playerDTO m.PlayerDTO) (m.Player, error)
	DeletePlayer(id primitive.ObjectID) (string, error)
	GetPlayerById(id primitive.ObjectID) (m.Player, error)
	AddGameToPlayer(id primitive.ObjectID, game m.Game) (m.Game, error)
	GetPlayerGames(id primitive.ObjectID) ([]m.Game, error)
}
