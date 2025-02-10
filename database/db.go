package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	connStr := "user=alimchik dbname=toysdb password=lego1234 host=localhost sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Delete all existing data in the table before creating a new one (optional)
	_, err = DB.Exec("DELETE FROM products")
	if err != nil {
		log.Fatal("Failed to delete data:", err)
	}

	// Create the table if it does not exist
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS products (
		id VARCHAR(50) PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT NOT NULL
	);`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	fmt.Println("Database initialized successfully!")
}
