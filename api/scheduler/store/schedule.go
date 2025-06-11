package store

import (
	"time"

	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/scheduler/models"
)

const (
	getAllQuery  = `SELECT id, name, timezone, week_schedules, created_at, updated_at FROM schedule WHERE deleted_at IS NULL`
	getByIDQuery = `SELECT id, name, timezone, week_schedules, created_at, updated_at FROM schedule WHERE id = ? AND deleted_at IS NULL`
	insertQuery  = `INSERT INTO schedule (name, timezone, week_schedules) VALUES (?, ?, ?)`
	updateQuery  = `UPDATE schedule SET week_schedules = ?, updated_at = ? WHERE id = ?`
	deleteQuery  = `UPDATE schedule SET deleted_at = ? WHERE id = ?`

	assumedMax = 100
)

type Store struct {
}

func New() *Store {
	return &Store{}
}

func (*Store) GetAllSchedule(ctx *gofr.Context) ([]*models.Schedule, error) {
	rows, err := ctx.SQL.QueryContext(ctx, getAllQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	resp := make([]*models.Schedule, 0, assumedMax)

	for rows.Next() {
		var sch models.Schedule

		err = rows.Scan(&sch.ID, &sch.Name, &sch.TimeZone, &sch.ScheduleString, &sch.CreatedAt, &sch.UpdatedAt)
		if err != nil {
			return resp, err
		}

		resp = append(resp, &sch)
	}

	return resp, nil
}

func (*Store) GetScheduleByID(ctx *gofr.Context, id int64) (*models.Schedule, error) {
	row := ctx.SQL.QueryRowContext(ctx, getByIDQuery, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var sch models.Schedule

	err := row.Scan(&sch.ID, &sch.Name, &sch.TimeZone, &sch.ScheduleString, &sch.CreatedAt, &sch.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &sch, nil
}

func (*Store) CreateSchedule(ctx *gofr.Context, sch *models.Schedule) (*models.Schedule, error) {
	res, err := ctx.SQL.ExecContext(ctx, insertQuery, sch.Name, sch.TimeZone, sch.ScheduleString)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	sch.ID = id

	return sch, nil
}

func (*Store) UpdateSchedule(ctx *gofr.Context, sch *models.Schedule) error {
	_, err := ctx.SQL.ExecContext(ctx, updateQuery, sch.ScheduleString, time.Now(), sch.ID)
	if err != nil {
		return err
	}

	return nil
}

func (*Store) DeleteSchedule(ctx *gofr.Context, id int64) error {
	_, err := ctx.SQL.ExecContext(ctx, deleteQuery, time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}
