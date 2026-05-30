package util

import "math/rand"

// jitter can be -1 or 1
func RandJitter(min, max int) int {
	res := rand.Intn(max-min) + min
	if res == 0 {
		res = RandJitter(min, max)
	}
	return res
}
