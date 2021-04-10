package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Player struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Status   string             `json:"status" bson:"status"`
	Priority int                `json:"priority" priority:"status"`
	Comment  string             `json:"comment" bson:"comment"`
	Date     int64              `json:"date" bson:"date"`
	Games    []Game             `json:"games" bson:"games"`
}

type Game struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Score      string             `json:"score" bson:"score"`
	IsWin      bool               `json:"is_win" bson:"is_win"`
	IsOvertime bool               `json:"is_overtime" bson:"is_overtime"`
}
