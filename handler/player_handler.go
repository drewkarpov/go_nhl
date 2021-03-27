package handler

import (
	"context"
	"encoding/json"
	d "github.com/drewkarpov/go_nhl/mongo"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

type Player struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Status   string             `json:"status,omitempty" bson:"status,omitempty"`
	Priority int                `json:"priority,omitempty" priority:"status,omitempty"`
	Comment  string             `json:"comment,omitempty" bson:"comment,omitempty"`
}

type Application struct {
	Db d.DbWrapper
}

func (a Application) New(database d.DbWrapper) Application {
	return Application{database}
}

func (a Application) CreatePlayer(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	var person Player
	json.NewDecoder(request.Body).Decode(&person)
	if person.Name == "" {
		response.WriteHeader(http.StatusUnprocessableEntity)
		response.Write([]byte(`{"message":"name is empty"}`))
		return
	}

	result, err := a.Db.Collection.InsertOne(ctx, person)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(response).Encode(result)
}

func (a Application) GetPlayers(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	var players []Player
	cursor, err := a.Db.Collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Player
		cursor.Decode(&person)
		players = append(players, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(players)
}

func (a Application) GetPlayerById(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person Player
	err := a.Db.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&person)

	if err != nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(&person)

}
