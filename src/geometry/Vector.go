package geometry

import (
	"math"
)

// A Vector in a 3D real-valued vector space
type Vector struct {
	X float64
	Y float64
	Z float64
}

// Get the length of the vector
func (v *Vector) GetLength() float64 {
	return math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2) + math.Pow(v.Z, 2))
}

// Multiply the vector with a scalar value
func (v *Vector) ScalarMultiplication(scalar float64) {
	v.X *= scalar
	v.Y *= scalar
	v.Z *= scalar
}

// Vector addition
func (v *Vector) AddVector(vector Vector) {
	v.X += vector.X
	v.Y += vector.Y
	v.Z += vector.Z
}

// Compare the 3 coordinates of the vectors and return if all of them are equal
func EqualVector(v1, v2 Vector) bool {
	return v1.X == v2.X && v1.Y == v2.Y && v1.Z == v2.Z
}

// Compare the 3 coordinates of the vectors and return if the difference between
// each of them is smaller than the given delta
func EqualVectorInDelta(v1, v2 Vector, delta float64) bool {
	return (math.Abs(v1.X-v2.X) <= delta) && (math.Abs(v1.Y-v2.Y) <= delta) && (math.Abs(v1.Z-v2.Z) <= delta)
}

// Computes the distance between the two given vectors
func VectorDist(v1, v2 Vector) float64 {
	return math.Sqrt(math.Pow(v1.X-v2.X, 2) + math.Pow(v1.Y-v2.Y, 2) + math.Pow(v1.Z-v2.Z, 2))
}
