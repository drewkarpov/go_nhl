package handlers

import (
	"encoding/json"
	"github.com/drewkarpov/go_nhl/internal/app"
	m "github.com/drewkarpov/go_nhl/internal/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"sort"
	"time"
)

type PlayerHandler struct {
	Application *app.Application
}

type Result struct {
	Result string `json:"result" bson:"result"`
}

func (handler PlayerHandler) CreatePlayer(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var player m.PlayerDTO
	json.NewDecoder(request.Body).Decode(&player)
	player.Timestamp = time.Now().Unix()
	result := handler.Application.PlayerService.WritePlayer(player)

	json.NewEncoder(response).Encode(result)
	handler.Application.Logger.Infof("get response %v", request.RequestURI)
}

func (handler PlayerHandler) GetPlayers(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	response.Header().Add("Access-Control-Allow-Origin", "*")
	var players, err = handler.Application.PlayerService.GetAllPlayers()
	if err != nil {
		writeErrorToResponse(response, err)
		return
	}
	sort.Slice(players, func(i, j int) bool { return players[i].Timestamp > players[j].Timestamp })

	json.NewEncoder(response).Encode(players)
	handler.Application.Logger.Infof("get response %v", request.RequestURI)
}

func (handler PlayerHandler) GetPlayerById(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var player, err = handler.Application.PlayerService.GetPlayerById(id)

	if err != nil {
		writeErrorToResponse(response, err)
		return
	}
	json.NewEncoder(response).Encode(&player)
	handler.Application.Logger.Infof("get response %v", request.RequestURI)
}

func (handler PlayerHandler) ChangePlayerById(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var player m.PlayerDTO
	json.NewDecoder(request.Body).Decode(&player)

	chandedPlayer, err := handler.Application.PlayerService.ChangePlayerData(id, player)

	if err != nil {
		writeErrorToResponse(response, err)
		return
	}
	json.NewEncoder(response).Encode(&chandedPlayer)
	handler.Application.Logger.Infof("get response %v", request.RequestURI)
}

func (handler PlayerHandler) DeletePlayer(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	resultStatus, err := handler.Application.PlayerService.DeletePlayer(id)

	if err != nil {
		statusCode := writeErrorToResponse(response, err)
		handler.LoggingRequest(*request, statusCode)
		return
	}
	json.NewEncoder(response).Encode(Result{resultStatus})
	handler.LoggingRequest(*request, 200)
}

func (handler PlayerHandler) AddGameToPlayer(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var game m.Game
	decodeError := json.NewDecoder(request.Body).Decode(&game)
	if decodeError != nil {
		writeErrorToResponse(response, decodeError)
	}
	game.Timestamp = time.Now().Unix()
	game, err := handler.Application.PlayerService.AddGameToPlayer(id, game)

	if err != nil {
		writeErrorToResponse(response, err)
		return
	}
	json.NewEncoder(response).Encode(&game)
	handler.Application.Logger.Infof("get response %v", request.RequestURI)
}

func (handler PlayerHandler) GetPlayerGames(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	games, err := handler.Application.PlayerService.GetPlayerGames(id)

	if err != nil {
		writeErrorToResponse(response, err)
		return
	}
	sort.Slice(games, func(i, j int) bool {
		return games[i].Timestamp > games[j].Timestamp
	})
	json.NewEncoder(response).Encode(&games)
	handler.Application.Logger.Infof("get response %v", request.RequestURI)
}

func writeErrorToResponse(response http.ResponseWriter, err error) int {
	response.WriteHeader(http.StatusUnprocessableEntity)
	response.Write([]byte(`{"message":"` + err.Error() + `"}`))
	return http.StatusUnprocessableEntity
}

func (h *PlayerHandler) LoggingRequest(request http.Request, statusCode int) {
	h.Application.Logger.Infof("method:%v path:%v code:%v", request.Method, request.RequestURI, statusCode)
}
