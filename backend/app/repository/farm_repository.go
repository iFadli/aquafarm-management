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

// FarmRepository mengandung informasi database dan mengandung metode-metode yang dibutuhkan untuk melakukan CRUD pada database
type FarmRepository struct {
	DB *sql.DB
}

// NewFarmRepository membuat repository baru
func NewFarmRepository(db *DbRepository) *FarmRepository {
	return &FarmRepository{DB: db.DB}
}

// Fetch mengambil semua data pada tabel Farm
func (r *FarmRepository) Fetch() ([]model.Farm, error) {
	query := `	SELECT farm_id, farm_name, created_at, updated_at
				FROM farms
				WHERE is_deleted = ?`

	rows, err := r.DB.Query(query, false)
	if err != nil {
		return nil, fmt.Errorf("error querying items: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error closing rows: %v", err)
		}
	}()

	var farms []model.Farm
	for rows.Next() {
		var farm model.Farm
		if err := rows.Scan(&farm.ID, &farm.Name, &farm.CreatedAt, &farm.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		farms = append(farms, farm)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error fetching rows: %w", err)
	}

	return farms, nil
}

// GetById mengambil data Farm berdasarkan farm_id
func (r *FarmRepository) GetById(id string) (*model.Farm, bool, error) {
	query := `	SELECT farm_id, farm_name, created_at, updated_at
				FROM farms
				WHERE farm_id = ? AND is_deleted = ?`
	row := r.DB.QueryRow(query, id, 0)

	var farm model.Farm
	if err := row.Scan(&farm.ID, &farm.Name, &farm.CreatedAt, &farm.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, true, nil
		}
		return nil, false, err
	}

	return &farm, false, nil
}

// Store menyimpan data Farm Baru pada tabel Farm
func (r *FarmRepository) Store(farm *model.Farm) (*model.Farm, error) {
	// Buat perintah SQL untuk menyimpan data item baru
	query := `
		INSERT INTO farms (farm_id, farm_name)
		VALUES (?, ?)
	`

	// Jalankan perintah SQL
	result, err := r.DB.Exec(query, farm.ID, farm.Name)
	if err != nil {
		return nil, err
	}

	// Ambil KEY (Primary Key) dari Farm baru yang disimpan
	key, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Update FARM_KEY pada struct Farm
	farm.KEY = strconv.FormatInt(key, 10)

	return farm, nil
}

// UpdateById memperbarui 1 data Farm berdasarkan farm_key, atau jika data tidak ada maka tambah data baru
func (r *FarmRepository) UpdateById(farm *model.Farm) (bool, error) {
	query := `UPDATE farms SET farm_name=?, updated_at=? WHERE farm_id=? AND is_deleted=?`
	stmt, err := r.DB.Prepare(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	timeNow := time.Now()

	res, err := stmt.Exec(farm.Name, timeNow, farm.ID, false)
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

// SoftDeleteById menghapus 1 data Farm berdasarkan farm_id
func (r *FarmRepository) SoftDeleteById(id string) error {
	query := `UPDATE farms SET is_deleted=? WHERE farm_id=? AND is_deleted=?`
	stmt, err := r.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(true, id, false)
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
