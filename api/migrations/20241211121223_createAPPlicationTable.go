package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

func createAPPlicationTable() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			const query = `
				CREATE TABLE if not exists application (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					name VARCHAR(255) NOT NULL,
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
