package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

func addResourcesTable() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			// write your migrations here
			_, err := d.SQL.Exec(`CREATE TABLE IF NOT EXISTS resources (
										id INTEGER PRIMARY KEY AUTOINCREMENT,
										resource_uid VARCHAR(255) NOT NULL UNIQUE,
										name VARCHAR(255) NOT NULL,
										state CHAR(10) NOT NULL,
										region VARCHAR(50) NOT NULL,
										cloud_account_id BIGINT NOT NULL,
										cloud_provider VARCHAR(50) NOT NULL,
										resource_type VARCHAR(50) NOT NULL,
										settings TEXT NOT NULL,
										created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
										updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
