package updater

import (
	"fmt"
	"strconv"
	"time"

	"github.com/CCPupp/states.osutools/internal/api"
	"github.com/CCPupp/states.osutools/internal/player"
)

var IsUpdating = false

func StartUpdate() {
	go updateLoop()
	go backupLoop()
}

func updateLoop() {
	for {
		doUpdate()
		time.Sleep(8 * time.Hour)
	}
}

func doUpdate() {
	IsUpdating = true
	fmt.Println("Starting Update")
	clientToken := api.GetClientToken()
	for i := 0; i < len(player.UserList); i++ {
		updateUser(player.UserList, i, clientToken)
	}
	IsUpdating = false
	fmt.Println("Update Complete!")
}

func updateUser(users []player.User, i int, clientToken string) {
	updatedUser := api.GetUser(strconv.Itoa(users[i].ID), clientToken)
	player.NewOverwriteExistingUser(player.NewGetUserById(users[i].ID), updatedUser)
}

func backupLoop() {
	for {
		backup()
		time.Sleep(30 * time.Minute)
	}
}

func backup() {
	player.NewBackupUserJSON()
}
