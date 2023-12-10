package duel

import (
	"log"
	"math/rand"
)

// Duel return true if dick1 wins.
func Duel(dick1 int, dick2 int) (bool, float64, float64) {
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
