package handlers

import (
	"encoding/json"
	app "github.com/drewkarpov/go_nhl/internal/app"
	m "github.com/drewkarpov/go_nhl/internal/models"
	"net/http"
)

type StatisticHandler struct {
	Application *app.Application
}

func (handler *StatisticHandler) GetStatistic(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	players, err := handler.Application.PlayerService.GetAllPlayers()
	if err != nil {
		writeErrorToResponse(response, err)
		return
	}
	var statistic = m.Statistic{}

	for index := range players {
		statistic.Total++
		switch players[index].Status {
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
	json.NewEncoder(response).Encode(statistic)
	handler.LoggingRequest(*request, 200)
}

func (h *StatisticHandler) LoggingRequest(request http.Request, statusCode int) {
	h.Application.Logger.Infof("method:%v path:%v code:%v", request.Method, request.RequestURI, statusCode)
}
