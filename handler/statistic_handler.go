package handler

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

type Statistic struct {
	Zero   int64 `json:"zero,omitempty" bson:"zero,omitempty"`
	Low    int64 `json:"low,omitempty" bson:"low,omitempty"`
	Hard   int64 `json:"hard,omitempty" bson:"hard,omitempty"`
	Driver int64 `json:"driver,omitempty" bson:"driver,omitempty"`
}

func (a Application) GetStatistic(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	var statistic = Statistic{0, 0, 0, 0}
	cursor, err := a.Db.Collection.Find(ctx, bson.M{})
	if err != nil {
		writeErrorToResponse(response, err)
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var player Player
		cursor.Decode(&player)
		switch player.Status {
		case "zero":
			statistic.Zero++
		case "low":
			statistic.Low++
		case "hard":
			statistic.Hard++
		case "driver":
			statistic.Driver++
		}
	}
	if err := cursor.Err(); err != nil {
		writeErrorToResponse(response, err)
		return
	}
	json.NewEncoder(response).Encode(statistic)
}
