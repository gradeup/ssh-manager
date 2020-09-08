package apis

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"k8s.io/apimachinery/pkg/util/json"
)

type Access struct {
	User_id    string `json:"user_id"`
	Server_id  string `json:"server_id"`
	Grant_date string `json:"grant_date"`
	User       string `json:"user"`
	Server     string `json:"server"`
}

type AccessList struct {
	Accesses []Access `json:"accesses"`
	Servers  []Server `json:"servers"`
	Users    []User   `json:"users"`
}

func ToggleAccess(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
	user_id := r.FormValue("user_id")
	server_id := r.FormValue("server_id")
	access := r.FormValue("access")
	grant_date := time.Now()

	user_id_int, err := strconv.Atoi(user_id)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		log.Printf("%v", err)
		return nil
	}

	server_id_int, err := strconv.Atoi(server_id)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		log.Printf("%v", err)
		return nil
	}

	if access == "true" {
		sqlStatement := `INSERT INTO user_server(user_id, server_id, grant_date) VALUES ($1, $2, $3)`
		log.Printf(sqlStatement)
		_, err = db.Exec(sqlStatement, user_id_int, server_id_int, grant_date)
	} else {
		sqlStatement := `DELETE FROM user_server where user_id = $1 and server_id = $2`
		log.Printf(sqlStatement)
		_, err = db.Exec(sqlStatement, user_id_int, server_id_int)
	}

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		log.Printf("%v", err)
		return nil
	}
	w.WriteHeader(200)
	w.Write([]byte("Access modified!"))
	return nil
}

func GetAccess(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
	servers, err := ListServers(db)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		log.Printf("%v", err)
		return nil
	}
	users, err := ListUsers(db)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		log.Printf("%v", err)
		return nil
	}

	rows, err := db.Query("SELECT us.*, u.username as user, s.username as server FROM user_server us, users u, servers s where us.user_id=u.id and us.server_id=s.id")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		log.Printf("%v", err)
		return nil
	}
	defer rows.Close()
	var accesses []Access
	for rows.Next() {
		var access Access
		err = rows.Scan(&access.User_id, &access.Server_id, &access.Grant_date, &access.User, &access.Server)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			log.Printf("%v", err)
			return nil
		}
		accesses = append(accesses, access)
	}
	err = rows.Err()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		log.Printf("%v", err)
		return nil
	}

	var accessList = AccessList{
		Accesses: accesses,
		Servers:  servers,
		Users:    users,
	}

	accessListByte, err := json.Marshal(accessList)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return err
	}
	w.WriteHeader(200)
	w.Write(accessListByte)
	return nil
}
