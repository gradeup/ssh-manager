package libraries

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func getDb(psHost string, psPort int, psUser string, psPassword string, psDatabase string) (*sql.DB, error) {

	// Connect to postgres DB
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", psHost, psPort, psUser, psPassword, "postgres")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// Check if required database exists
	statement := `SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '` + psDatabase + `');`
	row := db.QueryRow(statement)
	var exists bool
	err = row.Scan(&exists)
	if err != nil {
		return nil, err
	}

	// Create database if it does not exist
	if exists == false {
		statement = `CREATE DATABASE ` + psDatabase + `;`
		_, err = db.Exec(statement)
		if err != nil {
			return nil, err
		}
	}

	// Create connection with the required Database
	psqlInfoNew := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", psHost, psPort, psUser, psPassword, psDatabase)
	dbNew, err := sql.Open("postgres", psqlInfoNew)
	if err != nil {
		return nil, err
	}

	// Verify if Connection is established
	err = dbNew.Ping()
	if err != nil {
		return nil, err
	}

	return dbNew, nil

}

func createTables(db *sql.DB) error {
	return nil
}

func seedData(db *sql.DB) error {
	return nil
}

func GetPostgresClient(psHost string, psPort int, psUser string, psPassword string, psDatabase string) (*sql.DB, error) {

	// Get Database instance
	db, err := getDb(psHost, psPort, psUser, psPassword, psDatabase)
	if err != nil {
		log.Println("Error connecting/creating database")
		return nil, err
	}

	// Bootstrap required tables
	err = createTables(db)
	if err != nil {
		log.Println("Error creating required tables")
		return nil, err
	}

	// Seed Admin User & Keys
	err = seedData(db)
	if err != nil {
		log.Println("Error creating admin users & access")
		return nil, err
	}

	return db, nil
}
