package main

import (
	"fmt"
	a "github.com/drewkarpov/go_nhl/handler"
	d "github.com/drewkarpov/go_nhl/mongo"
	"github.com/gorilla/mux"

	"net/http"
)

func main() {
	var application a.Application
	var db d.DbWrapper

	db = db.Init()
	application = application.New(db)

	mux.NewRouter()
	router := mux.NewRouter()
	router.HandleFunc("/add/player", application.CreatePlayer).Methods("POST")
	router.HandleFunc("/players", application.GetPlayers).Methods("GET")
	router.HandleFunc("/player/{id}", application.GetPlayerById).Methods("GET")
	router.HandleFunc("/player/{id}/change", application.ChangePlayerById).Methods("PUT")
	router.HandleFunc("/player/{id}/delete", application.DeletePlayer).Methods("DELETE")
	fmt.Println("Application is started and listen port 2222")
	http.ListenAndServe(":2222", router)
}
