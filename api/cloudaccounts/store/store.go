package store

import (
	"encoding/json"

	"database/sql"

	"gofr.dev/pkg/gofr"
)

type Store struct {
}

// New creates a new instance of CloudAccountStore.
func New() CloudAccountStore {
	return &Store{}
}

// InsertCloudAccount inserts a new cloud account into the database.
func (*Store) InsertCloudAccount(ctx *gofr.Context, cloudAccount *CloudAccount) (*CloudAccount, error) {
	jsonCredentials, err := json.Marshal(cloudAccount.Credentials)
	if err != nil {
		return nil, err
	}

	res, err := ctx.SQL.ExecContext(ctx, INSERTQUERY, cloudAccount.Name, cloudAccount.Provider, cloudAccount.ProviderID,
		cloudAccount.ProviderDetails, jsonCredentials)
	if err != nil {
		return nil, err
	}

	cloudAccount.ID, err = res.LastInsertId()
	cloudAccount.Credentials = nil

	return cloudAccount, err
}

// GetALLCloudAccounts retrieves all cloud accounts from the database.
func (*Store) GetALLCloudAccounts(ctx *gofr.Context) ([]CloudAccount, error) {
	rows, err := ctx.SQL.QueryContext(ctx, GETALLQUERY)
	if err != nil {
		return nil, err
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	cloudAccounts := make([]CloudAccount, 0)

	for rows.Next() {
		cloudAccount := CloudAccount{}

		var providerDetails sql.NullString

		err := rows.Scan(&cloudAccount.ID, &cloudAccount.Name, &cloudAccount.Provider, &cloudAccount.ProviderID,
			&providerDetails, &cloudAccount.CreatedAt, &cloudAccount.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if providerDetails.Valid {
			cloudAccount.ProviderDetails = providerDetails.String
		}

		cloudAccounts = append(cloudAccounts, cloudAccount)
	}

	return cloudAccounts, nil
}

// GetCloudAccountByProvider retrieves a cloud account by provider type and provider Identifier.
func (*Store) GetCloudAccountByProvider(ctx *gofr.Context, providerType, providerID string) (*CloudAccount, error) {
	row := ctx.SQL.QueryRowContext(ctx, GETBYPROVIDERQUERY, providerType, providerID)

	if row.Err() != nil {
		return nil, row.Err()
	}

	cloudAccount := CloudAccount{}

	var providerDetails sql.NullString

	err := row.Scan(&cloudAccount.ID, &cloudAccount.Name, &cloudAccount.Provider, &cloudAccount.ProviderID,
		&providerDetails, &cloudAccount.CreatedAt, &cloudAccount.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if providerDetails.Valid {
		cloudAccount.ProviderDetails = providerDetails.String
	}

	return &cloudAccount, nil
}

// GetCloudAccountByID retrieves a cloud account by id.
func (*Store) GetCloudAccountByID(ctx *gofr.Context, cloudAccountID int64) (*CloudAccount, error) {
	row := ctx.SQL.QueryRowContext(ctx, GETBYPROVIDERIDQUERY, cloudAccountID)

	if row.Err() != nil {
		return nil, row.Err()
	}

	cloudAccount := CloudAccount{}

	var providerDetails sql.NullString

	err := row.Scan(&cloudAccount.ID, &cloudAccount.Name, &cloudAccount.Provider, &cloudAccount.ProviderID,
		&providerDetails, &cloudAccount.CreatedAt, &cloudAccount.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if providerDetails.Valid {
		cloudAccount.ProviderDetails = providerDetails.String
	}

	return &cloudAccount, nil
}

func (*Store) GetCredentials(ctx *gofr.Context, cloudAccountID int64) (interface{}, error) {
	row := ctx.SQL.QueryRowContext(ctx, GETCREDENTIALSQUERY, cloudAccountID)

	if row.Err() != nil {
		return nil, row.Err()
	}

	var credentials string

	err := row.Scan(&credentials)
	if err != nil {
		return nil, err
	}

	var jsonCred map[string]string

	err = json.Unmarshal([]byte(credentials), &jsonCred)
	if err != nil {
		return nil, err
	}

	return jsonCred, nil
}
