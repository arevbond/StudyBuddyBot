package game

import (
	"log"
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
	_, monthLastTry, dayLastTry := user.DateChangeDick.Date()
	_, month, today := time.Now().Date()
	return (month == monthLastTry && today > dayLastTry) || month > monthLastTry
}

// Duel return true if dick1 wins.
func Duel(dick1 int, dick2 int) (bool, float64, float64) {
	rand.Seed(time.Now().UnixNano())
	allChance := dick1 + dick2
	chance1 := float64(dick1) / float64(allChance) * 100
	chance2 := float64(dick2) / float64(allChance) * 100

	if chance1 <= 0 {
		chance1 = 2
		chance2 = 98
	} else if chance2 <= 0 {
		chance1 = 98
		chance2 = 2
	}

	result := float64(rand.Intn(100))
	log.Printf("duel between dick1 = %d and dick2 = %d. chance1 = %f and chance2 = %f", dick1, dick2, chance1, chance2)
	return result <= chance1, chance1, chance2
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func PositiveRandomValue() int {
	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(20)
	if value < 10 {
		value = 12
	}
	return value
}
