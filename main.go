package main

import (
	"fmt"
	a "github.com/drewkarpov/go_nhl/app"
	h "github.com/drewkarpov/go_nhl/handler"
	d "github.com/drewkarpov/go_nhl/mongo"
	"github.com/gorilla/mux"

	"net/http"
)

func main() {
	var application a.Application
	var db d.DbWrapper
	var statHandler h.StatisticHandler
	var playerHandler h.PlayerHandler
	db = db.Init()
	application = application.New(db)
	statHandler = statHandler.New(application)
	playerHandler = playerHandler.New(application)
	mux.NewRouter()
	router := mux.NewRouter()
	router.HandleFunc("/add/player", playerHandler.CreatePlayer).Methods("POST")
	router.HandleFunc("/players", playerHandler.GetPlayers).Methods("GET")
	router.HandleFunc("/player/{id}", playerHandler.GetPlayerById).Methods("GET")
	router.HandleFunc("/player/{id}/change", playerHandler.ChangePlayerById).Methods("PUT")
	router.HandleFunc("/player/{id}/delete", playerHandler.DeletePlayer).Methods("DELETE")
	router.HandleFunc("/players/statistic", statHandler.GetStatistic).Methods("GET")
	fmt.Println("Application is started and listen port 2222")
	http.ListenAndServe(":2222", router)
}
