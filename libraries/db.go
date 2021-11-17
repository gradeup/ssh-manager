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

	createUserTableStatement := "CREATE TABLE IF NOT EXISTS users (id serial PRIMARY KEY, username VARCHAR (50) UNIQUE NOT NULL, email VARCHAR (100) NOT NULL, public_key TEXT NOT NULL, created_at timestamp);"
	_, err := db.Exec(createUserTableStatement)
	if err != nil {
		return err
	}

	createServerTableStatement := "CREATE TABLE IF NOT EXISTS servers (id serial PRIMARY KEY, ip VARCHAR (50) UNIQUE NOT NULL, username VARCHAR (50) NOT NULL, created_at timestamp);"
	_, err = db.Exec(createServerTableStatement)
	if err != nil {
		return err
	}

	createPivotTableStatement := "CREATE TABLE IF NOT EXISTS user_server (user_id integer not null, server_id integer not null, grant_date timestamp, PRIMARY KEY (user_id, server_id), CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES users (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE, CONSTRAINT server_id_fk FOREIGN KEY (server_id) REFERENCES servers (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE);"
	_, err = db.Exec(createPivotTableStatement)
	if err != nil {
		return err
	}

	return nil
}

func seedData(db *sql.DB) error {
	// Nothing to seed as of now
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
