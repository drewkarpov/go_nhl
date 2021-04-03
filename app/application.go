package app

import (
	d "github.com/drewkarpov/go_nhl/mongo"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Application struct {
	Server    http.Server
	Db        d.DbWrapper
	AppConfig AppConfig
	Logger    *logrus.Logger
}

type AppConfig struct {
	Router mux.Router
}

func (a *Application) Setup(logger *logrus.Logger, database d.DbWrapper, router mux.Router) Application {
	a.AppConfig = AppConfig{Router: router}
	return Application{Db: database, AppConfig: a.AppConfig, Logger: logger}
}

func (app *Application) Run() {
	shutdown := make(chan error, 1)

	go func() {
		app.Logger.Info("Application is started and listen port 2222")
		err := http.ListenAndServe(":2222", &app.AppConfig.Router)
		shutdown <- err

	}()

	app.Logger.Infof("%v", <-shutdown)

}
