package store

const (
	INSERTQUERY = "INSERT INTO environment (name,level,application_id) VALUES ( ?,?, ?);"
	GETALLQUERY = "SELECT id, name,level,application_id, created_at, updated_at FROM environment WHERE application_id = ? " +
		"and deleted_at IS NULL;"
	GETBYNAMEQUERY = "SELECT id, name,level,application_id, created_at, updated_at FROM environment WHERE name = ? " +
		"and application_id = ? and deleted_at IS NULL;"
	UPDATEQUERY = "UPDATE environment SET name = ?, level = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?;"
	GETMAXLEVEL = `
    SELECT MAX(level) AS highest_level
    FROM environment
    WHERE application_id = ? AND deleted_at IS NULL;
`
)
