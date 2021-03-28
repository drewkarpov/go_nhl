package app

import d "github.com/drewkarpov/go_nhl/mongo"

type Application struct {
	Db d.DbWrapper
}

func (a Application) New(database d.DbWrapper) Application {
	return Application{database}
}
