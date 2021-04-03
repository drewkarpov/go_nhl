package handler

import (
	"context"
	"encoding/json"
	app "github.com/drewkarpov/go_nhl/app"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

type StatisticHandler struct {
	application *app.Application
}

func (h *StatisticHandler) New(app *app.Application) StatisticHandler {
	return StatisticHandler{app}
}

type Statistic struct {
	Zero   int64 `json:"zero" bson:"zero"`
	Low    int64 `json:"low" bson:"low"`
	Hard   int64 `json:"hard" bson:"hard"`
	Driver int64 `json:"driver" bson:"driver"`
	Total  int64 `json:"total" bson:"total"`
}

func (handler *StatisticHandler) GetStatistic(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	var statistic = Statistic{0, 0, 0, 0, 0}
	cursor, err := handler.application.Db.Collection.Find(ctx, bson.M{})
	if err != nil {
		writeErrorToResponse(response, err)
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var player Player
		cursor.Decode(&player)
		statistic.Total++
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
	handler.LoggingRequest(*request, 200)
}

func (h *StatisticHandler) LoggingRequest(request http.Request, statusCode int) {
	h.application.Logger.Infof("method:%v path:%v code:%v", request.Method, request.RequestURI, statusCode)
}
