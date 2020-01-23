package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sshmanager/apis"
	"sshmanager/libraries"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"

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

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load environment variables from .env
	psHost := os.Getenv("POSTGRES_HOST")
	psPort := os.Getenv("POSTGRES_PORT")
	psUser := os.Getenv("POSTGRES_USER")
	servicePort := os.Getenv("SERVICE_PORT")
	psPassword := os.Getenv("POSTGRES_PASSWORD")
	psDatabase := os.Getenv("POSTGRES_DATABASE")
	privateKeyFile := os.Getenv("PRIVATE_KEY_PATH")
	publicKeyFile := os.Getenv("PUBLIC_KEY_PATH")
	// _ = os.Getenv("PUBLIC_KEY_PATH")

	// Get postgres DB connection
	psPortInt, err := strconv.Atoi(psPort)
	_, err = libraries.GetPostgresClient(psHost, psPortInt, psUser, psPassword, psDatabase)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apis.HomePage(w, r, publicKeyFile)
	})))
	log.Println("Server started at :" + servicePort)
	log.Fatal(http.ListenAndServe(":"+servicePort, handlers.RecoveryHandler()(mux)))

	cmd := "ls"
	hosts := []string{}
	results := make(chan string, 10)
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
