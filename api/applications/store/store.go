package store

import (
	"time"

	"gofr.dev/pkg/gofr"
)

type Store struct{}

func New() ApplicationStore {
	return &Store{}
}
func (*Store) InsertApplication(ctx *gofr.Context, application *Application) (*Application, error) {
	res, err := ctx.SQL.ExecContext(ctx, INSERTQUERY, application.Name)
	if err != nil {
		return nil, err
	}

	application.ID, err = res.LastInsertId()
	application.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	return application, err
}

func (*Store) GetALLApplications(ctx *gofr.Context) ([]Application, error) {
	rows, err := ctx.SQL.QueryContext(ctx, GETALLQUERY)
	if err != nil {
		return nil, err
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	applications := make([]Application, 0)

	for rows.Next() {
		application := Application{}

		err = rows.Scan(&application.ID, &application.Name, &application.CreatedAt, &application.UpdatedAt)
		if err != nil {
			return nil, err
		}

		applications = append(applications, application)
	}

	return applications, nil
}

func (*Store) GetApplicationByName(ctx *gofr.Context, name string) (*Application, error) {
	row := ctx.SQL.QueryRowContext(ctx, GETBYNAMEQUERY, name)
	if row.Err() != nil {
		return nil, row.Err()
	}

	application := Application{}

	err := row.Scan(&application.ID, &application.Name, &application.CreatedAt, &application.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &application, nil
}

func (*Store) GetApplicationByID(ctx *gofr.Context, id int) (*Application, error) {
	row := ctx.SQL.QueryRowContext(ctx, GETBYIDQUERY, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	application := Application{}

	err := row.Scan(&application.ID, &application.Name, &application.CreatedAt, &application.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &application, nil
}

func (*Store) InsertEnvironment(ctx *gofr.Context, environment *Environment) (*Environment, error) {
	res, err := ctx.SQL.ExecContext(ctx, INSERTENVIRONMENTQUERY, environment.Name, environment.Level, environment.ApplicationID)
	if err != nil {
		return nil, err
	}

	environment.ID, _ = res.LastInsertId()

	return environment, nil
}
