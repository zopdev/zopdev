package store

import (
	"database/sql"
	"errors"

	"gofr.dev/pkg/gofr"
)

type Store struct{}

func New() *Store { return &Store{} }

func (*Store) GetLastRun(ctx *gofr.Context, cloudAccountID int64, rule string) (*Result, error) {
	var res Result

	row := ctx.SQL.QueryRowContext(ctx,
		"SELECT id, cloud_account_id, rule_id, result, evaluated_at "+
			"FROM results WHERE cloud_account_id = ? AND rule_id = ? ORDER BY evaluated_at DESC LIMIT 1",
		cloudAccountID, rule)

	err := row.Scan(&res.ID, &res.CloudAccountID, &res.RuleID, &res.Result, &res.EvaluatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		ctx.Metrics().IncrementCounter(ctx,
			"db_error_count", "audit_store", "GetLastRun", "error", err.Error())

		return nil, err
	}

	return &res, nil
}

func (*Store) CreatePending(ctx *gofr.Context, result *Result) (*Result, error) {
	res, err := ctx.SQL.ExecContext(ctx,
		"INSERT INTO results (cloud_account_id, rule_id, evaluated_at) VALUES (?, ?, ?)",
		result.CloudAccountID, result.RuleID, result.EvaluatedAt)
	if err != nil {
		ctx.Metrics().IncrementCounter(ctx,
			"db_error_count", "audit_store", "CreatePending", "error", err.Error())

		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	result.ID = id

	return result, nil
}

func (*Store) UpdateResult(ctx *gofr.Context, result *Result) error {
	_, err := ctx.SQL.ExecContext(ctx, "UPDATE results SET result = ? WHERE id = ?",
		result.Result, result.ID)
	if err != nil {
		ctx.Metrics().IncrementCounter(ctx,
			"db_error_count", "audit_store", "UpdateResult", "error", err.Error())

		return err
	}

	return nil
}
