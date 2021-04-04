package main

import (
	a "github.com/drewkarpov/go_nhl/app"
	h "github.com/drewkarpov/go_nhl/handler"
	d "github.com/drewkarpov/go_nhl/mongo"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	var application a.Application
	var db d.DbWrapper
	db = db.Init()

	statHandler := h.StatisticHandler{Application: &application}
	playerHandler := h.PlayerHandler{Application: &application}

	router := mux.NewRouter()
	router.HandleFunc("/add/player", playerHandler.CreatePlayer).Methods("POST")
	router.HandleFunc("/players", playerHandler.GetPlayers).Methods("GET")
	router.HandleFunc("/player/{id}", playerHandler.GetPlayerById).Methods("GET")
	router.HandleFunc("/player/{id}/change", playerHandler.ChangePlayerById).Methods("PUT")
	router.HandleFunc("/player/{id}/delete", playerHandler.DeletePlayer).Methods("DELETE")
	router.HandleFunc("/players/statistic", statHandler.GetStatistic).Methods("GET")

	application = application.Setup(logger, db, *router)
	application.Run()

}
