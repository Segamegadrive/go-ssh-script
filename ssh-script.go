package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"
)

func main() {

	username := "ubuntu"
	ipaddr := "10.5.31.113"
	port := "22"
	command := "ls"

	key, err := ioutil.ReadFile("/Users/Sagar Gurung/Downloads/5gc-key.pem")
	if err != nil {
		log.Fatalf("Private key cannot be found", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("Unable to parse private key", err)
	}

	hostKeyCallback, err := kh.New(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		log.Fatalf("Could not create hostkeycallback: ", err)
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,
	}

	client, err := ssh.Dial("tcp", ipaddr+":"+port, config)
	if err != nil {
		log.Fatalf("Unable to ssh into the remote server", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(command); err != nil {
		log.Fatalf("Failed to run; " + err.Error())
	}
	fmt.Println(b.String())

}
