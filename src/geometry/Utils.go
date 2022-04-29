package geometry

import (
	"math"
)

// Calculate the distance from the given plane to the given point
func DistFromPlane(planeNormalVector, point Vector) float64 {
	return math.Abs(planeNormalVector.X*point.X + planeNormalVector.Y*point.Y + planeNormalVector.Z*point.Z)
}
