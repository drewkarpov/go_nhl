package handler

import (
	"context"
	"encoding/json"
	app "github.com/drewkarpov/go_nhl/app"
	m "github.com/drewkarpov/go_nhl/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

type StatisticHandler struct {
	Application *app.Application
}

func (handler *StatisticHandler) GetStatistic(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	var statistic = m.Statistic{0, 0, 0, 0, 0}
	cursor, err := handler.Application.Db.Collection.Find(ctx, bson.M{})
	if err != nil {
		writeErrorToResponse(response, err)
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var player m.Player
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
	h.Application.Logger.Infof("method:%v path:%v code:%v", request.Method, request.RequestURI, statusCode)
}
