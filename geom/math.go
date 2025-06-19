package geom

import (
	"cmp"
	"math"
)

const (
	// 0.001
	Epsilon32 = 1e-3
	// 0.000001
	Epsilon64 = 1e-6

	// sqrt(2) / 2 == 1/sqrt(2)
	Sqrt2Over2 = math.Sqrt2 / 2
)

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Abs returns the absolute value of the number.
func Abs[T Number](s T) T {
	if s < 0 {
		return -s
	}
	return s
}

// Clamp clamps the value between min and max.
func Clamp[T cmp.Ordered](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// IsFinite checks if the number is finite.
func IsFinite[T Number](s T) bool {
	return !(math.IsNaN(float64(s)) || math.IsInf(float64(s), 0))
}
