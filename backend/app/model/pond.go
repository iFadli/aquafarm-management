package model

type Pond struct {
	KEY       string `json:"pond_key,omitempty"`
	ID        string `json:"pond_id"`
	Name      string `json:"pond_name"`
	FarmKey   string `json:"farm_key,omitempty"`
	FarmID    string `json:"farm_id"`
	FarmName  string `json:"farm_name"`
	IsDeleted bool   `json:"is_deleted,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type PondShow struct {
	ID        string `json:"pond_id"`
	Name      string `json:"pond_name"`
	FarmID    string `json:"farm_id"`
	FarmName  string `json:"farm_name"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type PondStore struct {
	FarmID string `json:"farm_id"`
	ID     string `json:"pond_id"`
	Name   string `json:"pond_name"`
}
