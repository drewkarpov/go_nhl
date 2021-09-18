package main

import (
	a "github.com/drewkarpov/go_nhl/internal/app"
	h "github.com/drewkarpov/go_nhl/internal/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	var application a.Application

	statHandler := h.StatisticHandler{Application: &application}
	playerHandler := h.PlayerHandler{Application: &application}

	router := mux.NewRouter()
	router.HandleFunc("/add/player", playerHandler.CreatePlayer).Methods("POST")
	router.HandleFunc("/add/players", playerHandler.CreatePlayers).Methods("POST")
	router.HandleFunc("/player/{id}/game/add", playerHandler.AddGameToPlayer).Methods("POST")
	router.HandleFunc("/player/{id}/games", playerHandler.GetPlayerGames).Methods("GET")
	router.HandleFunc("/players", playerHandler.GetPlayers).Methods("GET")
	router.HandleFunc("/player/{id}", playerHandler.GetPlayerById).Methods("GET")
	router.HandleFunc("/player/{id}/change", playerHandler.ChangePlayerById).Methods("PUT")
	router.HandleFunc("/player/{id}/delete", playerHandler.DeletePlayer).Methods("DELETE")
	router.HandleFunc("/players/statistic", statHandler.GetStatistic).Methods("GET")

	application = application.Setup(logger, *router)
	application.Run()

}
