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
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	Description    string  `json:"description,omitempty"`
	CloudAccountID int64   `json:"cloud_account_id,omitempty"`
	ResourceIDs    []int64 `json:"resource_ids,omitempty"`
}

type RGUpdate RGCreate
