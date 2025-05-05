package migrations

import (
	"fmt"

	"gofr.dev/pkg/gofr/migration"
)

func createEnvironmentTable() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			const query = `
				CREATE TABLE IF NOT EXISTS environment (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					name VARCHAR(255) NOT NULL,
					level INTEGER NOT NULL,
					application_id INTEGER NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP
				);
			`

			_, err := d.SQL.Exec(query)
			if err != nil {
				fmt.Println(err)
				return err
			}

			return nil
		},
	}
}
