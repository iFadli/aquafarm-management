package repository

import (
	"aquafarm-management/app/model"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// PondRepository mengandung informasi database dan mengandung metode-metode yang dibutuhkan untuk melakukan CRUD pada database
type PondRepository struct {
	DB *sql.DB
}

// NewPondRepository membuat repository baru
func NewPondRepository(db *DbRepository) *PondRepository {
	return &PondRepository{DB: db.DB}
}

// Fetch mengambil semua data pada tabel Pond
func (r *PondRepository) Fetch() ([]model.Pond, error) {
	query := `	SELECT f.farm_id, f.farm_name, p.pond_id, p.pond_name, p.created_at, p.updated_at
				FROM ponds p
				INNER JOIN farms f on p.farm_key = f.farm_key
				WHERE p.is_deleted = ? AND f.is_deleted = ?`

	rows, err := r.DB.Query(query, false, false)
	if err != nil {
		return nil, fmt.Errorf("error querying items: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error closing rows: %v", err)
		}
	}()

	var Ponds []model.Pond
	for rows.Next() {
		var Pond model.Pond
		if err := rows.Scan(&Pond.FarmID, &Pond.FarmName, &Pond.ID, &Pond.Name, &Pond.CreatedAt, &Pond.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		Ponds = append(Ponds, Pond)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error fetching rows: %w", err)
	}

	return Ponds, nil
}

// GetById mengambil data Pond berdasarkan Pond_id
func (r *PondRepository) GetById(farmKey string, pondId string) (*model.Pond, bool, error) {
	query := `	SELECT f.farm_id, f.farm_name, p.pond_id, p.pond_name, p.created_at, p.updated_at
				FROM ponds p
				INNER JOIN farms f on p.farm_key = f.farm_key
				WHERE p.is_deleted = ? AND f.is_deleted = ? AND f.farm_key = ? AND p.pond_id = ?
				LIMIT 1`

	row := r.DB.QueryRow(query, false, false, farmKey, pondId)

	var Pond model.Pond
	if err := row.Scan(&Pond.FarmID, &Pond.FarmName, &Pond.ID, &Pond.Name, &Pond.CreatedAt, &Pond.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, true, nil
		}
		return nil, false, err
	}

	return &Pond, false, nil
}

// Store menyimpan data Pond Baru pada tabel Pond
func (r *PondRepository) Store(pond *model.Pond) (*model.Pond, error) {
	// Buat perintah SQL untuk menyimpan data item baru
	query := `
		INSERT INTO ponds (farm_key, pond_id, pond_name)
		VALUES (?, ?, ?)
	`

	// Jalankan perintah SQL
	result, err := r.DB.Exec(query, pond.FarmKey, pond.ID, pond.Name)
	if err != nil {
		return nil, err
	}

	// Ambil KEY (Primary Key) dari Pond baru yang disimpan
	key, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Update Pond_KEY pada struct Pond
	pond.KEY = strconv.FormatInt(key, 10)

	return pond, nil
}

// UpdateById memperbarui 1 data Pond berdasarkan Pond_key, atau jika data tidak ada maka tambah data baru
func (r *PondRepository) UpdateById(pond *model.Pond) (bool, error) {
	query := `UPDATE ponds SET pond_name=?, updated_at=? WHERE farm_key=? AND pond_id=? AND is_deleted=?`
	stmt, err := r.DB.Prepare(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	timeNow := time.Now()

	res, err := stmt.Exec(pond.Name, timeNow, pond.FarmKey, pond.ID, false)
	if err != nil {
		return false, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if affected != 1 {
		return false, errors.New("failed to update the data. please try again later")
	}
	primaryLastInsert, err := res.LastInsertId()
	if err != nil {
		return true, errors.New("data not found. try to insert new data but failed. please try again later")
	}
	if primaryLastInsert > 0 {
		return true, nil
	}
	return false, nil
}

// SoftDeleteById menghapus 1 data Pond berdasarkan Pond_id
func (r *PondRepository) SoftDeleteById(farmKey *string, pondId string) error {
	query := `UPDATE ponds SET is_deleted=? WHERE farm_key=? AND pond_id=? AND is_deleted=?`
	stmt, err := r.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(true, farmKey, pondId, false)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.New("failed to delete the data. please try again later")
	}
	return nil
}
