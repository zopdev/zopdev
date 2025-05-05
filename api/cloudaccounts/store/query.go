package store

const (
	INSERTQUERY = "INSERT INTO cloud_account (name, provider,provider_id,provider_details ,credentials) values(? , ?, ?, ? ,?);"
	GETALLQUERY = "SELECT id, name, provider, provider_id, provider_details, created_at, updated_at FROM cloud_account" +
		" WHERE deleted_at IS NULL;"
	GETBYPROVIDERQUERY = "SELECT id, name, provider, provider_id, provider_details, created_at," +
		" updated_at FROM cloud_account WHERE provider = ? " +
		"AND provider_id = ? AND deleted_at IS NULL;"
	GETBYPROVIDERIDQUERY = "SELECT id, name, provider, provider_id, provider_details, created_at," +
		" updated_at FROM cloud_account WHERE " +
		"id = ? AND deleted_at IS NULL;"
	//nolint:gosec //query
	GETCREDENTIALSQUERY = "SELECT credentials from cloud_account WHERE id = ? AND deleted_at IS NULL;"
)
