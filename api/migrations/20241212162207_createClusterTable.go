package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

func createClusterTable() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			const query = `
				CREATE TABLE if not exists cluster (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					deployment_space_id INTEGER NOT NULL,
					cluster_id INTEGER,
					name VARCHAR(255) NOT NULL,
				    region VARCHAR(255) NOT NULL,
				    provider_id VARCHAR(255) NOT NULL,
				    provider VARCHAR(255) NOT NULL,
				    namespace VARCHAR(255),
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
