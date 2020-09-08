package apis

import (
	"database/sql"
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/util/json"
)

type User struct {
	Id         string `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Public_key string `json:"public_key"`
	Created_at string `json:"created_at"`
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

func ListUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.Username, &user.Email, &user.Public_key, &user.Created_at)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetUsers(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.Username, &user.Email, &user.Public_key, &user.Created_at)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return err
		}
		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return err
	}
	usersByte, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return err
	}
	w.WriteHeader(200)
	w.Write(usersByte)
	return nil
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

func getUser(user_id int, db *sql.DB) (User, error) {
	var user User
	sqlStatement := "SELECT * FROM users where id = $1"
	row := db.QueryRow(sqlStatement, user_id)
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Public_key, &user.Created_at)
	if err != nil {
		return user, err
	}

	return user, nil
}
