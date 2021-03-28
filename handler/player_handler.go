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
	"time"
)

type PlayerHandler struct {
	application app.Application
}

func (h PlayerHandler) New(app app.Application) PlayerHandler {
	return PlayerHandler{app}
}

type Player struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Status   string             `json:"status,omitempty" bson:"status,omitempty"`
	Priority int                `json:"priority,omitempty" priority:"status,omitempty"`
	Comment  string             `json:"comment,omitempty" bson:"comment,omitempty"`
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

	result, err := handler.application.Db.Collection.InsertOne(ctx, player)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(response).Encode(result)
}

func (handler PlayerHandler) GetPlayers(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	var players []Player
	cursor, err := handler.application.Db.Collection.Find(ctx, bson.M{})
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
	json.NewEncoder(response).Encode(players)
}

func (handler PlayerHandler) GetPlayerById(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person Player

	err := handler.application.Db.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&person)

	if err != nil {
		writeErrorToResponse(response, err)
		return
	}
	json.NewEncoder(response).Encode(&person)

}

func (handler PlayerHandler) ChangePlayerById(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var player Player
	json.NewDecoder(request.Body).Decode(&player)

	_, err := handler.application.Db.Collection.UpdateOne(
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

}

func (handler PlayerHandler) DeletePlayer(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	result, err := handler.application.Db.Collection.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		writeErrorToResponse(response, err)
		return
	}
	json.NewEncoder(response).Encode(result)
}

func writeErrorToResponse(response http.ResponseWriter, err error) {
	log.Println(err)
	response.WriteHeader(http.StatusUnprocessableEntity)
	response.Write([]byte(`{"message":"` + err.Error() + `"}`))
}
