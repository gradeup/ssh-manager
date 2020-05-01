package apis

import (
	"database/sql"
	"net/http"
	"time"
)

type User struct {
	id         string
	username   string
	email      string
	public_key string
	created_at string
}

func AddUser(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
	username := r.FormValue("username")
	email := r.FormValue("email")
	public_key := r.FormValue("public_key")

	sqlStatement := `INSERT INTO users (username, email, public_key, created_at) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(sqlStatement, username, email, public_key, time.Now())
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return nil
	}
	w.WriteHeader(200)
	w.Write([]byte("User created"))
	return nil
}

func GetUsers(w http.ResponseWriter, r *http.Request, db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.id, &user.username, &user.email, &user.public_key, &user.created_at)
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func DeleteUser(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
	username := r.FormValue("username")

	sqlStatement := `DELETE FROM users where username = ?`
	_, err := db.Exec(sqlStatement, username)
	if err != nil {
		return err
	}
	return nil
}
