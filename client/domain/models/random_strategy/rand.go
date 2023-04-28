package random_strategy

import (
	"math/rand"
	"time"
)

func SetRandom() {
	rand.Seed(time.Now().UnixNano())
}
