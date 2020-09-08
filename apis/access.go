package apis

import (
	"bytes"
	"database/sql"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"
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

type Result struct {
	Error    error
	Response string
}

func ToggleAccess(w http.ResponseWriter, r *http.Request, db *sql.DB, privateKeyFile string) error {
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

	// Get user's public key
	user, err := getUser(user_id_int, db)

	// Get server's IP address
	log.Printf("%v", server_id_int)
	server, err := getServer(server_id_int, db)

	// Update user access on the Server
	if access == "true" {
		err = updateAccess(user, server, true, privateKeyFile)
	} else {
		err = updateAccess(user, server, false, privateKeyFile)
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		log.Printf("%v", err)
		return nil
	}

	// Update user_server table for updated access
	if access == "true" {
		sqlStatement := `INSERT INTO user_server(user_id, server_id, grant_date) VALUES ($1, $2, $3)`
		_, err = db.Exec(sqlStatement, user_id_int, server_id_int, grant_date)
	} else {
		sqlStatement := `DELETE FROM user_server where user_id = $1 and server_id = $2`
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

func executeCmd(cmd, hostname string, config *ssh.ClientConfig) Result {
	conn, err := ssh.Dial("tcp", hostname+":22", config)
	if err != nil {
		log.Printf("dial failed:%v", err)
		return Result{
			Error:    err,
			Response: "",
		}
	}

	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil {
		log.Printf("session failed:%v", err)
		return Result{
			Error:    err,
			Response: "",
		}
	}

	defer session.Close()
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	err = session.Run(cmd)
	if err != nil {
		log.Printf("Run failed:%v", err)
		return Result{
			Error:    err,
			Response: "",
		}
	}

	return Result{
		Error:    nil,
		Response: stdoutBuf.String(),
	}
}

func updateAccess(user User, server Server, access bool, privateKeyFile string) error {

	var cmd string
	if access == true {
		cmd = "echo '" + user.Public_key + " " + user.Email + "' >> /home/ubuntu/.ssh/authorized_keys"
	} else {
		cmd = "sed -i -n '/" + user.Email + "/!p' /home/ubuntu/.ssh/authorized_keys"
	}
	results := make(chan Result, 5)
	timeout := time.After(30 * time.Second)

	pemBytes, err := ioutil.ReadFile(os.Getenv("HOME") + privateKeyFile)
	if err != nil {
		log.Fatal(err)
	}

	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		log.Fatalf("parse key failed:%v", err)
	}

	config := &ssh.ClientConfig{
		User:            "ubuntu",
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	log.Printf("%v", server.Ip)

	go func() {
		results <- executeCmd(cmd, server.Ip, config)
	}()

	select {
	case res := <-results:
		log.Printf("%v", res)
		return res.Error
	case <-timeout:
		log.Printf("SSH Timed out!")
		return errors.New("SSH Timed out!")
	}
}
