package client

type CloudAccount struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Provider    string `json:"provider"`
	Credentials any    `json:"credentials"`
}
