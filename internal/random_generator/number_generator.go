package random_generator

import (
	"github.com/yigitsadic/minigame/internal"
	"math/rand"
	"time"
)

func GenerateRandomNumber() int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	return r1.Intn(internal.PlayerLimit)
}
