package repository

import (
	"aquafarm-management/app/config"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

// DbRepository menangani akses data item
type DbRepository struct {
	DB *sql.DB
}

func NewDB(cfg *config.Config) *DbRepository {

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	))
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	//Menunggu koneksi DB berhasil; terdapat retry sebanyak 30.
	for i := 0; i < 30; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		log.Println("Waiting for database connection...")
		time.Sleep(time.Second)
	}

	//Set DB agar selalu Connect dengan maksimal koneksi.
	db.SetMaxOpenConns(3)
	db.SetMaxIdleConns(1)

	return &DbRepository{
		DB: db,
	}
}

//SetupDB untuk Create Database pertama kali Projek dijalankan namun hanya jika table belum ada.
//TODO:Gunakan Library SetupDB yang sudah ada; seperti gorm.
func SetupDB(db *DbRepository) {
	tablesConfig := []string{"farm", "pond", "access", "logs"}
	tables, _ := FetchTables(db)
	for _, table := range tables {
		tablesConfig = strings.Split(strings.Join(tablesConfig, "|"), table+"|")
	}

	for _, table := range tablesConfig {
		var err error
		var query string
		switch table {
		case "farm":
			query = `
					CREATE TABLE farms (
					    farm_key INT NOT NULL AUTO_INCREMENT COMMENT 'Primary Key of Farms' ,
					    farm_id VARCHAR(255) NOT NULL COMMENT 'ID of Farm' ,
					    farm_name VARCHAR(255) NOT NULL COMMENT 'Name of Farm' ,
					    is_deleted BOOLEAN NOT NULL DEFAULT FALSE COMMENT 'Soft Delete Flagger.' ,
					    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created Timestamp of Data' ,
					    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Updated Timestamp of Data' ,
					    PRIMARY KEY (farm_key)
                   ) ENGINE = InnoDB COMMENT = 'Main Table of Farm Data';`
			fmt.Println("Table Farm Not Found")
		case "pond":
			query = `
					CREATE TABLE ponds (
					    farm_key INT NOT NULL COMMENT 'Secondary Key of Ponds' ,
					    pond_key INT NOT NULL AUTO_INCREMENT COMMENT 'Primary Key of Ponds' ,
					    pond_id VARCHAR(255) NOT NULL COMMENT 'ID of Pond' ,
					    pond_name VARCHAR(255) NOT NULL COMMENT 'Name of Pond' ,
					    is_deleted BOOLEAN NOT NULL DEFAULT FALSE COMMENT 'Soft Delete Flagger.' ,
					    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created Timestamp of Data' ,
					    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Updated Timestamp of Data' ,
					    PRIMARY KEY (pond_key),
					    INDEX index_farm_key_on_ponds (farm_key)
					) ENGINE = InnoDB COMMENT = 'Main Table of Pond Data';`
			fmt.Println("Table Pond Not Found")
		case "access":
			query = `
					CREATE TABLE access (
    					access_id INT NOT NULL AUTO_INCREMENT COMMENT 'Primary Key of Access' ,
    					access_key VARCHAR(255) NOT NULL COMMENT 'Access Key for Header Request' ,
    					access_name VARCHAR(255) NOT NULL COMMENT 'Name of Access this Key for' ,
    					is_disabled BOOLEAN NOT NULL DEFAULT FALSE COMMENT 'Disable Flag' ,
    					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created Timestamp of Data' ,
    					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Updated Timestamp of Data' ,
    					PRIMARY KEY (access_id),
    					UNIQUE access_key_access_unique (access_key)
                    ) ENGINE = InnoDB COMMENT = 'Main Table of Access Data';`
			fmt.Println("Table Access Not Found")
		case "logs":
			query = `
					CREATE TABLE logs (
					    sequence INT NOT NULL AUTO_INCREMENT COMMENT 'PK of Logs' ,
					    access INT NULL COMMENT 'SK of Access Token' ,
					    details TEXT NULL DEFAULT NULL COMMENT 'Details Log (internal)' ,
					    request VARCHAR(255) NOT NULL COMMENT 'Request of This Request' ,
					    response VARCHAR(255) NOT NULL COMMENT 'Response of This Request' ,
					    ip_address VARCHAR(255) NOT NULL COMMENT 'IP Address of Client' ,
					    user_agent VARCHAR(255) NOT NULL COMMENT 'User Agent of Client' ,
					    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Time of Access' ,
					    PRIMARY KEY (sequence), INDEX ip_address_logs (ip_address),
					    INDEX user_agent_logs (user_agent),
					    INDEX access_logs_secondary_access (access),
					    INDEX request_logs (request),
					    INDEX response_logs (response)
                  	) ENGINE = InnoDB COMMENT = 'Main Table of Logs';`
			fmt.Println("Table Logs Not Found")
		default:
			fmt.Println("Unknown Table!")
		}

		if query != "" {
			_, err = db.DB.Exec(query)
			if err != nil {
				fmt.Println("Create Table " + table + " Failed.")
				fmt.Println(err.Error())
				return
			}
			fmt.Println("Create Table " + table + " Successfully.")

			if table == "access" {
				if InsertFirstAccess(db) {
					fmt.Println("Insert Demo Access Successfully.")
				} else {
					fmt.Println("Insert Demo Access Failed.")
				}
			}
		}
	}
}

func InsertFirstAccess(r *DbRepository) bool {
	query := `INSERT INTO access (access_key, access_name, is_disabled) VALUES (?, ?, ?)`

	result, err := r.DB.Exec(query, "ini_dia_si_jali_jali", "Demo Access API.", false)
	if err != nil {
		log.Fatalf("Error inserting data: %v", err)
		return false
	}
	// Check the number of rows affected
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0
}

// FetchTables mengambil semua nama table
func FetchTables(r *DbRepository) ([]string, error) {
	rows, err := r.DB.Query("SHOW TABLES")
	if err != nil {
		return nil, fmt.Errorf("error querying items: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error closing rows: %v", err)
		}
	}()

	var tables []string
	for rows.Next() {
		var tablesName []string
		var table string
		var err error
		if tablesName, err = rows.Columns(); err != nil {
			return nil, fmt.Errorf("error get table name: #{err}")
		}
		if len(tablesName) < 1 {
			return nil, fmt.Errorf("error get table name: tablesName length 0")
		}
		if err = rows.Scan(&table); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error fetching rows: %w", err)
	}

	return tables, nil
}
