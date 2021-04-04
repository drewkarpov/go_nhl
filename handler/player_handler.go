package handler

import (
	"context"
	"encoding/json"
	"github.com/drewkarpov/go_nhl/app"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"sort"
	"time"
)

type PlayerHandler struct {
	Application *app.Application
}

type Player struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Status   string             `json:"status,omitempty" bson:"status,omitempty"`
	Priority int                `json:"priority,omitempty" priority:"status,omitempty"`
	Comment  string             `json:"comment,omitempty" bson:"comment,omitempty"`
	Date     int64              `json:"date" bson:"date"`
}

func (handler PlayerHandler) CreatePlayer(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	var player Player
	json.NewDecoder(request.Body).Decode(&player)
	if player.Name == "" {
		writeErrorToResponse(response, nil)
		return
	}

	player.Date = time.Now().Unix()

	result, err := handler.Application.Db.Collection.InsertOne(ctx, player)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(response).Encode(result)
	handler.Application.Logger.Infof("get response %v", request.RequestURI)
}

func (handler PlayerHandler) GetPlayers(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	var players []Player
	cursor, err := handler.Application.Db.Collection.Find(ctx, bson.M{})
	if err != nil {
		writeErrorToResponse(response, err)
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Player
		cursor.Decode(&person)
		players = append(players, person)
	}
	if err := cursor.Err(); err != nil {
		writeErrorToResponse(response, err)
		return
	}
	sort.Slice(players, func(i, j int) bool { return players[i].Date > players[j].Date })

	json.NewEncoder(response).Encode(players)
	handler.Application.Logger.Infof("get response %v", request.RequestURI)
}

func (handler PlayerHandler) GetPlayerById(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person Player

	err := handler.Application.Db.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&person)

	if err != nil {
		writeErrorToResponse(response, err)
		return
	}
	json.NewEncoder(response).Encode(&person)
	handler.Application.Logger.Infof("get response %v", request.RequestURI)
}

func (handler PlayerHandler) ChangePlayerById(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var player Player
	json.NewDecoder(request.Body).Decode(&player)

	_, err := handler.Application.Db.Collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{{"name", player.Name}}},
			{"$set", bson.D{{"status", player.Status}}},
			{"$set", bson.D{{"priority", player.Priority}}},
			{"$set", bson.D{{"comment", player.Comment}}},
		},
	)
	if err != nil {
		writeErrorToResponse(response, err)
		return
	}
	json.NewEncoder(response).Encode(&player)
	handler.Application.Logger.Infof("get response %v", request.RequestURI)
}

func (handler PlayerHandler) DeletePlayer(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	result, err := handler.Application.Db.Collection.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		statusCode := writeErrorToResponse(response, err)
		handler.LoggingRequest(*request, statusCode)
		return
	}
	json.NewEncoder(response).Encode(result)
	handler.LoggingRequest(*request, 200)
}

func writeErrorToResponse(response http.ResponseWriter, err error) int {
	response.WriteHeader(http.StatusUnprocessableEntity)
	response.Write([]byte(`{"message":"` + err.Error() + `"}`))
	return http.StatusUnprocessableEntity
}

func (h *PlayerHandler) LoggingRequest(request http.Request, statusCode int) {
	h.Application.Logger.Infof("method:%v path:%v code:%v", request.Method, request.RequestURI, statusCode)
}
