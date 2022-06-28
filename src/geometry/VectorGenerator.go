package geometry

import (
	m "math"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"
)

// Create a random vector where each coordinate has a value between -1 and 1
func CreateRandomVec() Vector {
	x := utils.RandomFloat(-1, 1)
	y := utils.RandomFloat(-1, 1)
	z := utils.RandomFloat(-1, 1)

	return Vector{x, y, z}
}

// Same as CreateRandomVec but the vector has length 1
func CreateRandomUnitVec() Vector {
	x := utils.RandomFloat(-1, 1)
	y := utils.RandomFloat(-m.Sqrt(1-m.Pow(x, 2)), m.Sqrt(1-m.Pow(x, 2)))
	z := m.Sqrt(1 - m.Pow(x, 2) - m.Pow(y, 2)) // will always be positive
	if utils.RandomInt(0, 2) == 0 {
		z *= -1 // negate z with 50% probability
	}

	return Vector{x, y, z}
}

// Given a plane normal vector this function samples one point
// that is on the plane
func SamplePointFromPlane(planeNormalVector Vector) Vector {
	point := CreateRandomVec()
	dist := (planeNormalVector.X*point.X + planeNormalVector.Y*point.Y + planeNormalVector.Z*point.Z) / planeNormalVector.GetLength()
	planeNormalVector.ScalarMultiplication(-dist)

	point.AddVector(planeNormalVector)

	if DistFromPlane(planeNormalVector, point) > 0.00000001 {
		panic("The result of SamplePointFromPlane() was a point that is not on the provided plane")
	}
	return point
}

// Generate a point near the given plane with a distance to the plane which
// is sampled from a normal distribution
func SamplePointFromPlaneWithNoise(planeNormalVector Vector, noise utils.NormalDist) Vector {
	distance := utils.FloatFromNormalDist(noise)

	// add this vector to a sampled point from the plane
	point := SamplePointFromPlane(planeNormalVector)

	// scale the planeNormalVector s.t. it has the length of the distance we want to move
	planeNormalVector.ScalarMultiplication(distance / planeNormalVector.GetLength())

	point.AddVector(planeNormalVector)

	return point
}
