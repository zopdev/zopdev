package store

import "gofr.dev/pkg/gofr"

type ApplicationStore interface {
	InsertApplication(ctx *gofr.Context, application *Application) (*Application, error)
	GetALLApplications(ctx *gofr.Context) ([]Application, error)
	GetApplicationByName(ctx *gofr.Context, name string) (*Application, error)
	GetApplicationByID(ctx *gofr.Context, id int) (*Application, error)
	InsertEnvironment(ctx *gofr.Context, environment *Environment) (*Environment, error)
}
