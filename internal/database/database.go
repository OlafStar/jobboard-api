package database

import (
	"database/sql"
	"log"
)

type Store struct {
	db *sql.DB
}

func InitDatabase() Store {
	db, err := sql.Open("mysql", "root:root@(127.0.0.1:3306)/test-database?parseTime=true")
	if err != nil {
			log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
			log.Fatal(err)
	}

	createJobofferTable(db)
	createSalaryRangeTable(db)
	createCompanyTable(db)

	return Store{
			db: db,
	}
}

func createJobofferTable(db *sql.DB) {
	{
		query := `
		CREATE TABLE IF NOT EXISTS joboffers (
				id VARCHAR(36) NOT NULL,
				company_name VARCHAR(255) NOT NULL,
				company_size INT,
				company_website VARCHAR(255),
				company_logo TEXT,
				title VARCHAR(255) NOT NULL,
				experience TEXT NOT NULL,
				skill TEXT NOT NULL,
				description TEXT NOT NULL,
				location_country VARCHAR(100),
				location_city VARCHAR(100),
				location_address TEXT,
				operating_mode VARCHAR(50),
				type_of_work VARCHAR(50),
				apply_email VARCHAR(255),
				apply_url VARCHAR(255),
				consent BOOLEAN,
				contact_name VARCHAR(255),
				contact_email VARCHAR(255),
				contact_phone VARCHAR(50),
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				PRIMARY KEY (id)
		);`
		if _, err := db.Exec(query); err != nil {
				log.Fatal(err)
		}
	}
}

func createSalaryRangeTable(db *sql.DB) {
	{
		query := `
		CREATE TABLE IF NOT EXISTS salary_ranges (
				id INT AUTO_INCREMENT PRIMARY KEY,
				joboffer_id VARCHAR(36) NOT NULL,
				employment_type VARCHAR(100),
				min_salary DECIMAL(10,2),
				max_salary DECIMAL(10,2),
				currency VARCHAR(6),
				FOREIGN KEY (joboffer_id) REFERENCES joboffers(id)
		);`
		if _, err := db.Exec(query); err != nil {
				log.Fatal(err)
		}
	}
}

func createCompanyTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS companies (
		id VARCHAR(36) NOT NULL,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		country VARCHAR(100),
		city VARCHAR(100),
		address TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	);`
	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
	}
}