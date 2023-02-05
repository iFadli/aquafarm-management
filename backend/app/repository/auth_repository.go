package repository

import (
	"database/sql"
)

// AuthRepository mengandung informasi database dan mengandung metode-metode yang dibutuhkan untuk melakukan CRUD pada database
type AuthRepository struct {
	DB *sql.DB
}

// NewAuthRepository membuat repository baru
func NewAuthRepository(db *DbRepository) *AuthRepository {
	return &AuthRepository{DB: db.DB}
}

// GetApiKey mengambil data FarmKey dari Tabel Farm berdasarkan farm_id
func (r *AuthRepository) GetApiKey(accessKey string) (*string, error) {
	query := `	SELECT access_id
				FROM access
				WHERE access_key = ? AND is_disabled = ?`
	row := r.DB.QueryRow(query, accessKey, 0)

	var accessId *string
	if err := row.Scan(&accessId); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return accessId, nil
}
