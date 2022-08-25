//go:build !prod
// +build !prod

package main

import (
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func init() {

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path) // for example /home/user

	cmdString := `./start_ctrs.sh`
	cmd := exec.Command(cmdString)

	err = cmd.Run()

	if err != nil {
		log.Fatalf("unable to start database container %v", err)

	}

}
