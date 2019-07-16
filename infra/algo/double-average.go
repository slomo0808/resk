package algo

import (
	"math/rand"
	"time"
)

const min int64 = 1

func DoubleAverage(count, amount int64) int64 {
	if count == 1 {
		return amount
	}

	max := amount - min*count

	avg := max / count

	doubleAvg := 2*avg + min

	rand.Seed(time.Now().UnixNano())
	x := rand.Int63n(doubleAvg) + min
	return x
}
