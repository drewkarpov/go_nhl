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
	db d.DbWrapper
}

func (a Application) New(database d.DbWrapper) Application {
	a.db = database
	return a
}

func (a Application) CreatePlayer(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var person Player
	json.NewDecoder(request.Body).Decode(&person)
	if person.Name == "" {
		response.WriteHeader(http.StatusUnprocessableEntity)
		response.Write([]byte(`{"message":"name is empty"}`))
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := a.db.Collection.InsertOne(ctx, person)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(response).Encode(result)
}

func (a Application) GetPlayers(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var persons []Player
	cursor, err := a.db.Collection.Find(a.db.Ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(a.db.Ctx)
	for cursor.Next(a.db.Ctx) {
		var person Player
		cursor.Decode(&person)
		persons = append(persons, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(persons)
}

func (a Application) GetPlayerById(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person Player
	err := a.db.Collection.FindOne(a.db.Ctx, bson.M{"_id": id}).Decode(&person)

	if err != nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(&person)

}
