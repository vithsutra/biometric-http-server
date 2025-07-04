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

func (q *Query) InitilizeDatabase() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS admin (
			user_id VARCHAR PRIMARY KEY, 
			user_name VARCHAR NOT NULL UNIQUE, 
			password VARCHAR NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			user_id VARCHAR PRIMARY KEY, 
			user_name VARCHAR NOT NULL UNIQUE, 
			email VARCHAR NOT NULL UNIQUE,
			password VARCHAR NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS biometric (
			user_id VARCHAR NOT NULL, 
			unit_id VARCHAR PRIMARY KEY, 
			online BOOLEAN NOT NULL, 
			label VARCHAR NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS student(
			student_id VARCHAR PRIMARY KEY,
			unit_id VARCHAR NOT NULL,
			student_name VARCHAR NOT NULL,
			student_usn VARCHAR NOT NULL,
			department VARCHAR NOT NULL,
			FOREIGN KEY (unit_id) references biometric(unit_id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS fingerprintdata (
			student_id VARCHAR NOT NULL, 
			student_unit_id VARCHAR NOT NULL, 
			unit_id VARCHAR NOT NULL, 
			fingerprint VARCHAR NOT NULL, 
			FOREIGN KEY (unit_id) REFERENCES biometric(unit_id) ON DELETE CASCADE,
			PRIMARY KEY (student_id, student_unit_id)
		)`,
		`CREATE TABLE IF NOT EXISTS attendance (
			student_id VARCHAR NOT NULL, 
			date VARCHAR NOT NULL, 
			login VARCHAR NOT NULL, 
			logout VARCHAR NOT NULL, 
			FOREIGN KEY (student_id) REFERENCES student(student_id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS times (
			user_id VARCHAR NOT NULL, 
			morning_start VARCHAR NOT NULL, 
			morning_end VARCHAR NOT NULL, 
			afternoon_start VARCHAR NOT NULL, 
			afternoon_end VARCHAR NOT NULL, 
			evening_start VARCHAR NOT NULL, 
			evening_end VARCHAR NOT NULL, 
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS inserts(
			unit_id VARCHAR NOT NULL,
			student_unit_id VARCHAR NOT NULL,
			fingerprint_data VARCHAR NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS deletes(
			unit_id VARCHAR NOT NULL,
			student_unit_id VARCHAR NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS otps (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			otp VARCHAR(10) NOT NULL
 		)`,
		`CREATE TABLE IF NOT EXISTS student_unit_numbers (
			unit_id VARCHAR NOT NULL,
			student_unit_id VARCHAR NOT NULL,
			availability BOOL DEFAULT TRUE,
			PRIMARY KEY (unit_id, student_unit_id), 
			FOREIGN KEY (unit_id) REFERENCES biometric(unit_id) ON DELETE CASCADE
		)`,
	}

	tx, err := q.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			log.Println("Database initilized successfully")
		}
	}()

	for _, query := range queries {
		_, err = tx.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}
