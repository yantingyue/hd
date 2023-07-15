package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

func DialWithKeyFile() (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User:            os.Getenv("SSH_USER"),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if k, err := ioutil.ReadFile(os.Getenv("SSH_KEYFILE")); err != nil {
		return nil, err
	} else {
		signer, err := ssh.ParsePrivateKey(k)
		if err != nil {
			return nil, err
		}
		config.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
	}

	address := fmt.Sprintf("%s:%s", os.Getenv("SSH_HOST"), os.Getenv("SSH_PORT"))
	log.Printf("%s\n", address)
	log.Printf("%+v\n", config)
	return ssh.Dial("tcp", address, config)
}
