package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Connection struct {
	db *sql.DB
}

func NewDatabase() *Connection {
	conn, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	log.Println("Connected to database")
	return &Connection{
		db: conn,
	}
}

func (c *Connection) CheckStatus() {
	if err := c.db.Ping(); err != nil {
		log.Fatalf("Bad database connection: %v", err)
	}
	log.Println("Database is working correctly")
}

func (c *Connection) Close() {
	if err := c.db.Close(); err != nil {
		log.Printf("Unable to close database connection: %v", err)
		return
	}
	log.Println("Database connection closed successfully")
}
