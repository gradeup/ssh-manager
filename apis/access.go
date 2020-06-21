package apis

import (
	"database/sql"
	"fmt"
	"net/http"
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

func GetAccess(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
	servers, err := ListServers(db)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		fmt.Printf("%v", err)
		return nil
	}
	users, err := ListUsers(db)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		fmt.Printf("%v", err)
		return nil
	}

	rows, err := db.Query("SELECT us.*, u.username as user, s.username as server FROM user_server us, users u, servers s where us.user_id=u.id and us.server_id=s.id")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		fmt.Printf("%v", err)
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
			fmt.Printf("%v", err)
			return nil
		}
		accesses = append(accesses, access)
	}
	err = rows.Err()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		fmt.Printf("%v", err)
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
