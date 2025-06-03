package models

type ResourceGroup struct {
	ID             int64  `json:"id"`
	CloudAccountID int64  `json:"cloud_account_id"`
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	Status         string `json:"status,omitempty"`
}

type ResourceGroupData struct {
	ResourceGroup
	Resources []Resource `json:"resources,omitempty"`
}

type RGCreate struct {
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	CloudAccountID int64   `json:"cloud_account_id"`
	ResourceIDs    []int64 `json:"resource_ids"`
}

type RGUpdate struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	CloudAccountID int64   `json:"cloud_account_id"`
	ResourceIDs    []int64 `json:"resource_ids"`
}
