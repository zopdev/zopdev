package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/zopdev/zopdev/api/integration/models"
	"gofr.dev/pkg/gofr"
)

var (
	errFailedToSaveIntegration = errors.New("failed to save integration")
	errIntegrationNotFound     = errors.New("integration not found")
	errFailedToGetIntegration  = errors.New("failed to get integration")
)

type Store struct{}

func New() *Store {
	return &Store{}
}

func (s *Store) SaveIntegration(ctx *gofr.Context, i models.Integration) error {
	query := `INSERT INTO integrations (id, external_id, role_name, template_url) VALUES (?, ?, ?, ?)`

	_, err := ctx.SQL.ExecContext(ctx, query, i.IntegrationID, i.ExternalID, i.RoleName, i.TemplateURL)
	if err != nil {
		ctx.Logger.Errorf("%v: %v", errFailedToSaveIntegration, err)
		return errFailedToSaveIntegration
	}

	return nil
}

func (s *Store) GetIntegration(ctx *gofr.Context, id string) (models.Integration, error) {
	var i models.Integration

	query := fmt.Sprintf(`SELECT id, external_id, role_name, template_url FROM integrations WHERE id = '%s'`, id)
	err := ctx.SQL.QueryRowContext(ctx, query).Scan(&i.IntegrationID, &i.ExternalID, &i.RoleName, &i.TemplateURL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return i, fmt.Errorf("%w: %s", errIntegrationNotFound, id)
		}
		ctx.Logger.Errorf("%v: %v", errFailedToGetIntegration, err)

		return i, errFailedToGetIntegration
	}

	return i, nil
}
