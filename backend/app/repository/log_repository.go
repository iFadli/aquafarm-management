package repository

import (
	"aquafarm-management/app/model"
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

// LogRepository mengandung informasi database dan mengandung metode-metode yang dibutuhkan untuk melakukan CRUD pada database
type LogRepository struct {
	DB *sql.DB
}

// NewLogRepository membuat repository baru
func NewLogRepository(db *DbRepository) *LogRepository {
	return &LogRepository{DB: db.DB}
}

// Fetch mengambil semua data pada tabel Farm
func (r *LogRepository) Fetch() ([]model.Logs, error) {
	query := `	SELECT a.access_name, l.ip_address, l.user_agent, l.request, l.response, l.created_at
				FROM logs l
				LEFT JOIN access a on l.access = a.access_id`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying items: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error closing rows: %v", err)
		}
	}()

	var logs []model.Logs
	for rows.Next() {
		var log model.Logs
		if err := rows.Scan(&log.AccessName, &log.IpAddress, &log.UserAgent, &log.Request, &log.Response, &log.AccessedAt); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		logs = append(logs, log)
	}

	if logs == nil {
		logs = []model.Logs{}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error fetching rows: %w", err)
	}

	return logs, nil
}

// FirstLog menyimpan data Log pada tabel Logs
func (r *LogRepository) FirstLog(log *model.Logs) (*model.Logs, error) {
	// Buat perintah SQL untuk menyimpan data item baru
	query := `
		INSERT INTO logs (access, request, response, ip_address, user_agent)
		VALUES (?, ?, ?, ?, ?)
	`

	// Jalankan perintah SQL
	result, err := r.DB.Exec(query, log.AccessID, log.Request, log.Response, log.IpAddress, log.UserAgent)
	if err != nil {
		return nil, err
	}

	// Ambil KEY (Primary Key) dari Pond baru yang disimpan
	sequence, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Update Pond_KEY pada struct Pond
	log.Sequence = strconv.FormatInt(sequence, 10)

	return log, nil
}

// FetchStatistics mengambil perhitungan data pada tabel Logs
func (r *LogRepository) FetchStatistics() ([]model.StatisticsDB, error) {
	query := `	SELECT 
				  CASE
				    WHEN request LIKE 'GET /v1/pond/%' THEN 'GET /v1/pond/'
				    WHEN request LIKE 'DELETE /v1/pond/%' THEN 'DELETE /v1/pond/'
				    WHEN request LIKE 'GET /v1/farm/%' THEN 'GET /v1/farm/'
				    WHEN request LIKE 'DELETE /v1/farm/%' THEN 'DELETE /v1/farm/'
				    ELSE request
				  END as grouped_request,
				  COUNT(*) as count,
				  COUNT(DISTINCT user_agent) as unique_user_agent,
				  SUM(CASE WHEN response = '200' THEN 1 ELSE 0 END) as response_200,
				  SUM(CASE WHEN response = '404' THEN 1 ELSE 0 END) as response_404,
				  SUM(CASE WHEN response = '500' THEN 1 ELSE 0 END) as response_500,
				  SUM(CASE WHEN response != '200' AND response != '404' AND response != '500' THEN 1 ELSE 0 END) as response_etc
				FROM 
				  logs
				GROUP BY 
				  grouped_request;`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying items: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error closing rows: %v", err)
		}
	}()

	var logs []model.StatisticsDB
	for rows.Next() {
		var log model.StatisticsDB
		if err := rows.Scan(
			&log.Request, &log.Count, &log.UniqueUserAgent, &log.Response200,
			&log.Response404, &log.Response500, &log.ResponseETC); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error fetching rows: %w", err)
	}

	return logs, nil
}
