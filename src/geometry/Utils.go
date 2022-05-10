package geometry

import (
	"math"
)

// Calculate the distance from the given plane to the given point
func DistFromPlane(planeNormalVector, point Vector) float64 {
	return math.Abs(planeNormalVector.X*point.X+planeNormalVector.Y*point.Y+planeNormalVector.Z*point.Z) / planeNormalVector.GetLength()
}

// This function takes the parameters a and b of a plane (ax + by = z) that goes
// through the origin in coordinate representation and outputs the unified normal vector of the plane
func CoordinateReprToNormalVec(a, b float64) Vector {
	x := 1.0
	y := b / a
	z := -1 / a

	return Vector{X: x, Y: y, Z: z}
}
