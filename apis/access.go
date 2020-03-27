package apis

import (
	"database/sql"
	"net/http"
	"time"
)

type Access struct {
	user_id    string
	server_id  string
	grant_date string
}

func AddAccess(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
	user_id := r.FormValue("user_id")
	server_id := r.FormValue("server_id")
	grant_date := time.Now().UnixNano()

	sqlStatement := `INSERT INTO user_server (user_id, server_id, grant_date) VALUES (?, ?, ?)`
	_, err := db.Exec(sqlStatement, user_id, server_id, grant_date)
	if err != nil {
		return err
	}
	return nil
}

func GetAccess(w http.ResponseWriter, r *http.Request, db *sql.DB) ([]Access, error) {
	rows, err := db.Query("SELECT * FROM user_server")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var accesses []Access
	for rows.Next() {
		var access Access
		err = rows.Scan(&access.user_id, &access.server_id, &access.grant_date)
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return accesses, nil
}

func RevokeAccess(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
	user_id := r.FormValue("user_id")
	server_id := r.FormValue("server_id")

	sqlStatement := `DELETE FROM user_server where user_id = ? and server_id = ?`
	_, err := db.Exec(sqlStatement, user_id, server_id)
	if err != nil {
		return err
	}
	return nil
}
