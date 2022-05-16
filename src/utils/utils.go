package utils

import "math/rand"

type NormalDist struct {
	Mean, Stddev float64
}

type Numeric interface {
	int | int64 | uint | int32 | int16 | int8 | float32 | float64
}

// Return a random double in the given range, min is
// inclusive, max is exclusive
func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// Return a random integer in the given range, min is
// inclusive, max is exclusive
func RandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func FloatFromNormalDist(noise NormalDist) float64 {
	return rand.NormFloat64()*noise.Stddev + noise.Mean
}

// Check if an element is in the given slice
func Contains[T comparable](slice []T, element T) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}

// Returns the index of the smallest element in the slice.
// If multiple elements are the smallest then the index of the
// first element will be returned.
func ArgMin[T Numeric](slice []T) int {
	smallest := slice[0]
	index := 0
	for i, element := range slice {
		if element < smallest {
			smallest = element
			index = i
		}
	}
	return index
}
