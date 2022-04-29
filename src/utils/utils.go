package utils

import "math/rand"

type NormalDist struct {
	Mean, Stddev float64
}

func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func RandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func FloatFromNormalDist(noise NormalDist) float64 {
	return rand.NormFloat64()*noise.Stddev + noise.Mean
}
