package model

type Farm struct {
	KEY       string `json:"farm_key,omitempty"`
	ID        string `json:"farm_id"`
	Name      string `json:"farm_name"`
	IsDeleted bool   `json:"is_deleted,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type FarmShow struct {
	ID        string `json:"farm_id"`
	Name      string `json:"farm_name"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type FarmStore struct {
	ID   string `json:"farm_id"`
	Name string `json:"farm_name"`
}
