package apis

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/util/json"
)

type Server struct {
	Id         string `json:"id"`
	Ip         string `json:"ip"`
	Username   string `json:"username"`
	Created_at string `json:"created_at"`
}

func AddServer(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
	ip := r.FormValue("ip")
	username := r.FormValue("username")

	sqlStatement := `INSERT INTO servers (ip, username, created_at) VALUES ($1, $2, $3)`
	_, err := db.Exec(sqlStatement, ip, username, time.Now())
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		log.Printf("%v", err)
		return nil
	}
	w.WriteHeader(200)
	w.Write([]byte("Server Added"))
	return nil
}

func ListServers(db *sql.DB) ([]Server, error) {
	rows, err := db.Query("SELECT * FROM servers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var servers []Server
	for rows.Next() {
		var server Server
		err = rows.Scan(&server.Id, &server.Ip, &server.Username, &server.Created_at)
		if err != nil {
			return nil, err
		}
		servers = append(servers, server)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return servers, nil
}

func GetServers(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
	rows, err := db.Query("SELECT * FROM servers")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		log.Printf("%v", err)
		return nil
	}
	defer rows.Close()
	var servers []Server
	for rows.Next() {
		var server Server
		err = rows.Scan(&server.Id, &server.Ip, &server.Username, &server.Created_at)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			log.Printf("%v", err)
			return nil
		}
		servers = append(servers, server)
	}
	err = rows.Err()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		log.Printf("%v", err)
		return nil
	}
	serversByte, err := json.Marshal(servers)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return err
	}
	w.WriteHeader(200)
	w.Write(serversByte)
	return nil
}

func DeleteServer(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
	server_id := r.FormValue("server_id")

	sqlStatement := `DELETE FROM servers where id = $1`
	_, err := db.Exec(sqlStatement, server_id)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		log.Printf("%v", err)
		return nil
	}

	w.WriteHeader(200)
	w.Write([]byte("Server Deleted!"))
	return nil
}

func getServer(server_id int, db *sql.DB) (Server, error) {
	var server Server
	sqlStatement := "SELECT * FROM servers where id = $1"
	row := db.QueryRow(sqlStatement, server_id)
	err := row.Scan(&server.Id, &server.Ip, &server.Username, &server.Created_at)
	if err != nil {
		return server, err
	}

	return server, nil
}
