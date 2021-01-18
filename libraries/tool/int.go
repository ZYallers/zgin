package tool

import (
	"math/rand"
	"time"
)

// RandIntn
func RandIntn(max int) int {
	rad := rand.New(rand.NewSource(time.Now().Unix()))
	return rad.Intn(max)
}

