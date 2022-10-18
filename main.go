package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/fuadop/altschool-go-task-3/internal/atm"
)

const fileName = "kjb0sg47ms40000gn.json"

var users map[string]atm.User = map[string]atm.User{}
var currentUser *atm.User

func main() {
	sig := make(chan os.Signal)
	go listenForShutdown(sig)

	loadUsers()
	authPrompt(sig)
	opsPrompt(sig)
}

func listenForShutdown(sig chan os.Signal) {
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-sig
	gracefulShutdown(sig)
}

func gracefulShutdown(sig chan os.Signal) {
	fmt.Println("\nThanks for using our service. Have a great time.")
	close(sig)

	// save the users to the file cache
	tmpDir := os.TempDir()
	if file, err := os.Create(path.Join(tmpDir, fileName)); err == nil {
		defer file.Close()

		if currentUser != nil {
			users[currentUser.Username] = *currentUser
		}
		json.NewEncoder(file).Encode(&users)
	}

	os.Exit(0)
}

func loadUsers() {
	tmpDir := os.TempDir()
	buf, err := os.ReadFile(path.Join(tmpDir, fileName))
	if err != nil {
		buf = []byte("{}")
		os.Create(path.Join(tmpDir, fileName))
	}

	json.Unmarshal(buf, &users) // ignore the error
}
