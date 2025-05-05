package store

const (
	INSERTQUERY     = "INSERT INTO deployment_space (cloud_account_id, environment_id, type) VALUES ( ?, ?, ?);"
	GETQUERYBYENVID = `SELECT ds.id, ds.cloud_account_id, ds.environment_id, ds.type, ds.created_at,
       ds.updated_at, ca.Name, ev.name as ev_name
FROM deployment_space ds
         JOIN cloud_account ca JOIN environment ev ON ds.cloud_account_id = ca.id AND ds.environment_id = ev.id
WHERE ds.environment_id = ? AND ds.deleted_at IS NULL;`
)
