package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"sshmanager/libraries"

	"golang.org/x/crypto/ssh"
)

func executeCmd(cmd, hostname string, config *ssh.ClientConfig) string {

	conn, err := ssh.Dial("tcp", hostname+":22", config)
	if err != nil {
		log.Fatalf("dial failed:%v", err)
	}
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("session failed:%v", err)
	}
	defer session.Close()
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	err = session.Run(cmd)
	if err != nil {
		log.Fatalf("Run failed:%v", err)
	}
	return stdoutBuf.String()
}

func main() {
	psHost := flag.String("pshost", "localhost", "Hostname of Postgres database")
	psPort := flag.Int("psport", 5432, "Port to connect to postgres host")
	psUser := flag.String("psuser", "postgres", "Username to connect to postgres")
	psPassword := flag.String("pspassword", "postgres", "Password of the provided user")
	psDatabase := flag.String("psdatabase", "database", "Database to use")
	keyFile := flag.String("keypath", "/.ssh/id_rsa", "Relative Path of private key to use from HOME")

	flag.Parse()

	_, err := libraries.GetPostgresClient(*psHost, *psPort, *psUser, *psPassword, *psDatabase)
	if err != nil {
		panic(err)
	}

	cmd := "ls"
	hosts := []string{}
	results := make(chan string, 10)
	timeout := time.After(30 * time.Second)

	pemBytes, err := ioutil.ReadFile(os.Getenv("HOME") + *keyFile)
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

	for _, hostname := range hosts {
		go func(hostname string) {
			results <- executeCmd(cmd, hostname, config)
		}(hostname)
	}

	for i := 0; i < len(hosts); i++ {
		select {
		case res := <-results:
			fmt.Print(res)
		case <-timeout:
			fmt.Println("Timed out!")
			return
		}
	}
}
