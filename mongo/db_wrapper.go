package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

type DbWrapper struct {
	Client     *mongo.Client
	Collection *mongo.Collection
	Ctx        context.Context
}

func (d DbWrapper) New() DbWrapper {
	var cred options.Credential = options.Credential{Username: os.Getenv("MONGODB_USER"), Password: os.Getenv("MONGODB_PWD")}
	d.Ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	d.Client, _ = mongo.Connect(d.Ctx, options.Client().ApplyURI("mongodb://mongo:27017").SetAuth(cred))
	d.Collection = d.Client.Database("nhl").Collection("players")
	return DbWrapper{Client: d.Client, Ctx: d.Ctx, Collection: d.Collection}
}
