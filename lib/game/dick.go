package game

import (
	"math/rand"
	"tg_ics_useful_bot/storage"
	"time"
)

func RandomValue() int {
	rand.Seed(time.Now().UnixNano())
	sign := rand.Intn(5)
	value := rand.Intn(10)
	if sign > 0 {
		return value
	}
	return -1 * value
}

func CanChangeDickSize(user *storage.DBUser) bool {
	_, monthLastTry, dayLastTry := user.LastTryChangeDick.Date()
	_, month, today := time.Now().Date()
	return month >= monthLastTry && today > dayLastTry
}
