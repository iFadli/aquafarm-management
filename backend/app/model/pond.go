package model

type Pond struct {
	FarmKey   string `json:"farm_key,omitempty"`
	FarmID    string `json:"farm_id"`
	FarmName  string `json:"farm_name"`
	KEY       string `json:"pond_key,omitempty"`
	ID        string `json:"pond_id"`
	Name      string `json:"pond_name"`
	IsDeleted bool   `json:"is_deleted,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}
