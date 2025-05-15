package store

import (
	"database/sql"
	"fmt"
	"github.com/zopdev/zopdev/api/integration/models"

	"gofr.dev/pkg/gofr"
)

type Store struct{}

func New() *Store {
	return &Store{}
}

func (s *Store) SaveIntegration(ctx *gofr.Context, i models.Integration) error {
	query := `INSERT INTO integrations (id, external_id, role_name, template_url) VALUES (?, ?, ?, ?)`

	_, err := ctx.SQL.ExecContext(ctx, query, i.IntegrationID, i.ExternalID, i.RoleName, i.TemplateURL)
	if err != nil {
		ctx.Logger.Errorf("Failed to save integration: %v", err)
		return err
	}
	return nil
}

func (s *Store) GetIntegration(ctx *gofr.Context, id string) (models.Integration, error) {
	var i models.Integration

	query := fmt.Sprintf(`SELECT id, external_id, role_name, template_url FROM integrations WHERE id = '%s'`, id)
	err := ctx.SQL.QueryRowContext(ctx, query).Scan(&i.IntegrationID, &i.ExternalID, &i.RoleName, &i.TemplateURL)

	if err != nil {
		if err == sql.ErrNoRows {
			return i, fmt.Errorf("integration not found with id: %s", id)
		}
		ctx.Logger.Errorf("Failed to get integration: %v", err)
		return i, err
	}

	return i, nil
}
