package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

func createDeploymentSpaceTable() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			const query = `
				CREATE TABLE if not exists deployment_space (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					type VARCHAR(255) NOT NULL,
				    environment_id INTEGER NOT NULL,
				    cloud_account_id INTEGER NOT NULL,
				    details TEXT,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP
				);
				`

			_, err := d.SQL.Exec(query)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
