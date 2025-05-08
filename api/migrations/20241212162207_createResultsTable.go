package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

const (
	createResultsTableQuery = `CREATE TABLE IF NOT EXISTS results
(
    id               integer 				primary key,
    cloud_account_id int                                not null,
    rule_id          varchar(50)                        not null,
    result           text                               null,
    evaluated_at     datetime default CURRENT_TIMESTAMP null
);`

	addIndexQuery = `CREATE INDEX evaluation_index ON results (evaluated_at desc);`
)

func createTableAuditResults() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(createResultsTableQuery)
			if err != nil {
				return err
			}

			_, err = d.SQL.Exec(addIndexQuery)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
