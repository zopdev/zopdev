package migrations

import "gofr.dev/pkg/gofr/migration"

func addResourceGroup() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			// write your migrations here
			_, err := d.SQL.Exec(`CREATE TABLE IF NOT EXISTS resource_groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    cloud_account_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL);`)
			if err != nil {
				return err
			}

			_, err = d.SQL.Exec(`CREATE TABLE IF NOT EXISTS resource_group_memberships (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    resource_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (resource_id) REFERENCES resources(id),
    FOREIGN KEY (group_id) REFERENCES resource_groups(id),
    UNIQUE(resource_id, group_id)
);`)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
