package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

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
	cmd := "ls"
	hosts := []string{}
	results := make(chan string, 10)
	timeout := time.After(30 * time.Second)

	pemBytes, err := ioutil.ReadFile(os.Getenv("HOME") + "")
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
