package database

import (
	"database/sql"
	"log"
)

type Query struct {
	db *sql.DB
}

func NewQuery(db *sql.DB) *Query {
	return &Query{
		db,
	}
}

func(q *Query) InitilizeDatabase() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS admin (
			user_id VARCHAR(100) PRIMARY KEY, 
			user_name VARCHAR(50) NOT NULL UNIQUE, 
			password VARCHAR(100) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			user_id VARCHAR(100) PRIMARY KEY, 
			user_name VARCHAR(50) NOT NULL UNIQUE, 
			password VARCHAR(100) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS biometric (
			user_id VARCHAR(100), 
			unit_id VARCHAR(50) PRIMARY KEY, 
			online BOOLEAN NOT NULL, 
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS fingerprintdata (
			student_id VARCHAR(100) PRIMARY KEY, 
			student_unit_id VARCHAR(100), 
			unit_id VARCHAR(50), 
			fingerprint VARCHAR(1000), 
			FOREIGN KEY (unit_id) REFERENCES biometric(unit_id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS attendance (
			student_id VARCHAR(100), 
			student_unit_id VARCHAR(100), 
			unit_id VARCHAR(50), 
			date VARCHAR(20), 
			login VARCHAR(20), 
			logout VARCHAR(20), 
			FOREIGN KEY (unit_id) REFERENCES biometric(unit_id) ON DELETE CASCADE, 
			FOREIGN KEY (student_id) REFERENCES fingerprintdata(student_id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS times (
			user_id VARCHAR(200), 
			morning_start VARCHAR(20), 
			morning_end VARCHAR(20), 
			afternoon_start VARCHAR(20), 
			afternoon_end VARCHAR(20), 
			evening_start VARCHAR(20), 
			evening_end VARCHAR(20), 
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS inserts(
			unit_id VARCHAR(200),
			student_unit_id VARCHAR(200),
			fingerprint_data VARCHAR(1000)
		)`,
		`CREATE TABLE IF NOT EXISTS deletes(
			unit_id VARCHAR(200),
			student_unit_id VARCHAR(200)
		)`,
	}

	tx , err := q.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}else {
			tx.Commit()
			log.Println("Database initilized successfully")
		}
	} ()

	for _ , query := range queries {
		_ , err = tx.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

