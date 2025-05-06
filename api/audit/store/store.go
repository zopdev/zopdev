package store

import (
	"database/sql"
	"errors"

	"gofr.dev/pkg/gofr"
)

type Store struct{}

func New() *Store { return &Store{} }

func (*Store) GetAll(ctx *gofr.Context, cloudAccID int64, status string) ([]*Result, error) {
	rows, err := ctx.SQL.QueryContext(ctx, "SELECT * FROM results WHERE cloud_account_id = ? AND status = ?",
		cloudAccID, status)
	if err != nil || rows.Err() != nil {
		return nil, err
	}

	var results []*Result

	for rows.Next() {
		var res Result

		err = rows.Scan(&res.ID, &res.CloudAccountID, &res.RuleID, &res.Status, &res.Result, &res.EvaluatedAt)
		if err != nil {
			return results, err
		}

		results = append(results, &res)
	}

	return results, nil
}

func (*Store) GetLastRun(ctx *gofr.Context, cloudAccountID int64, rule string) (*Result, error) {
	var res Result

	row := ctx.SQL.QueryRowContext(ctx,
		"SELECT * FROM results WHERE cloud_account_id = ? AND rule_id = ? ORDER BY evaluated_at DESC LIMIT 1",
		cloudAccountID, rule)

	err := row.Scan(&res.ID, &res.CloudAccountID, &res.RuleID, &res.Status, &res.Result, &res.EvaluatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &res, nil
}

func (*Store) CreatePending(ctx *gofr.Context, result *Result) (*Result, error) {
	res, err := ctx.SQL.ExecContext(ctx,
		"INSERT INTO results (cloud_account_id, rule_id, status, evaluated_at) VALUES (?, ?, ?, ?)",
		result.CloudAccountID, result.RuleID, "pending", result.EvaluatedAt)
	if err != nil {
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
	_, err := ctx.SQL.ExecContext(ctx, "UPDATE results SET status = ?, result = ? WHERE id = ?",
		result.Status, result.Result, result.ID)
	if err != nil {
		return err
	}

	return nil
}
