package game

import (
	"fmt"
	"math"
	"math/rand"
)

const (
	smidge = 0.0001
)

// Generates a random number in a given range
func rnd(min, max int) int {
	return rand.Intn(max-min) + min
}

// Generates a random float number in a given range
func srnd(min, max float32) float32 {
	return rand.Float32()*(max-min) + min
}

// Generates a random number in [0,0xffffffff]
func urnd() uint32 {
	return uint32(rand.Intn(math.MaxInt32))
}

// Generates a random number in [0, 1]
func frnd() float64 {
	n := 1.0 + math.MaxInt32
	return float64(urnd()) / n
}

// Generates a random number in [-0.5, 0.5]
func crnd() float64 {
	m := int32(urnd()) - math.MinInt32
	n := 1.0 + math.MaxInt32
	return float64(m) / float64(n)
}

// Seconds to ticks (1/20th of a second)
func toTicks(seconds int, gamespeed float64) float64 {
	return float64(seconds) * 20 * gamespeed
}

// Weighted random range
func weightedRndRange(min, max float64) float64 {
	return math.Sqrt(min*min + frnd()*(max*max-min*min))
}

// Wraps f so it's within the range [smidge..(max-smidge)]
// Assumes f is not outside this range by more than (max - (2 * smidge))
func fwrap(f, max float64) float64 {
	upp := max - smidge
	rng := upp - smidge

	if f > upp {
		f -= rng
	}

	if f < smidge {
		f += rng
	}

	return f
}

// Checks if f is within the range [smidge..(max-smidge)]
func fclip(f, max float64) bool {
	return f < smidge || f >= (max-smidge)
}

// Contrains float
func fconstrain(f, max float64) float64 {
	max -= smidge

	if f > max {
		return max
	}

	if f < smidge {
		return smidge
	}

	return f
}

// Wraps f so it's within the range [min, max]
func fconstrain2(f, min, max float64) float64 {
	min += smidge
	max -= smidge

	if f > max {
		return max
	}

	if f < min {
		return min
	}

	return f
}

// Finds first null byte and returns the length
func clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}

// Reads uint32 from byte array
func readUint32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

// Formats time
func formatTime(s int, full bool) string {
	min := s / 60000
	sec := s / 1000 % 60
	tenths := s % 1000 / 100

	score := fmt.Sprintf("%2d:%.2d.%d", min, sec, tenths)

	if full {
		return score
	}

	if min == 0 {
		score = fmt.Sprintf("%2d.%d", sec, tenths)
	}

	return score
}
