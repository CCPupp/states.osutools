package stats

import (
	"time"

	"github.com/CCPupp/pupper.moe/internal/player"
)

var StatsRunning = false
var TotalUsers = 0

func StartStats() {
	for {
		StatsRunning = true
		setTotalUsers()
		time.Sleep(time.Second * 1)
	}
}

func setTotalUsers() {
	users := player.GetUserJSON()
	TotalUsers = len(users.Users)
}
