package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Player struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Status    string             `json:"status" bson:"status"`
	Priority  int                `json:"priority" priority:"status"`
	Comment   string             `json:"comment" bson:"comment"`
	Timestamp int64              `json:"timestamp" bson:"timestamp"`
	Games     []Game             `json:"games" bson:"games"`
}

type PlayerDTO struct {
	Name      string
	Status    string
	Priority  int
	Comment   string
	Timestamp int64
	Games     []Game
}

func (p PlayerDTO) MapToPlayer() Player {
	return Player{Name: p.Name, Status: p.Status, Priority: p.Priority, Comment: p.Comment, Timestamp: p.Timestamp, Games: p.Games}
}

type Game struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Score      string             `json:"score" bson:"score"`
	IsWin      bool               `json:"is_win" bson:"is_win"`
	IsOvertime bool               `json:"is_overtime" bson:"is_overtime"`
	Timestamp  int64              `json:"timestamp" bson:"timestamp"`
}
