package apis

import (
	"database/sql"
	"net/http"
	"time"
)

type Server struct {
	id         string
	ip         string
	username   string
	created_at string
}

func AddServer(w http.ResponseWriter, r *http.Request, db *sql.DB) error {

	ip := r.FormValue("ip")
	username := r.FormValue("username")
	created_at := time.Now().UnixNano()

	sqlStatement := `INSERT INTO servers (ip, username, created_at) VALUES (?, ?, ?)`
	_, err := db.Exec(sqlStatement, ip, username, created_at)
	if err != nil {
		return err
	}
	return nil
}

func GetServers(w http.ResponseWriter, r *http.Request, db *sql.DB) ([]Server, error) {
	rows, err := db.Query("SELECT * FROM servers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var servers []Server
	for rows.Next() {
		var server Server
		err = rows.Scan(&server.id, &server.ip, &server.username, &server.created_at)
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return servers, nil
}

func DeleteServer(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
	username := r.FormValue("username")

	sqlStatement := `DELETE FROM servers where username = ?`
	_, err := db.Exec(sqlStatement, username)
	if err != nil {
		return err
	}
	return nil
}
