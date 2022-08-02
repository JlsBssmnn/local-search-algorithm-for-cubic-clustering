package partitioning3D

import (
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
)

// A struct for calculating costs.
type CostCalculator struct {
	Threshold     float64 // The distance of a point to a plane where costs a zero
	Amplification float64 // A factor that is used in the cost calculation
}

// This function computes the cost for three points that are in the same partition
// It fits a plane through the points and calculates costs depending on the maximum
// distance of one point to the plane
func (calc CostCalculator) TripleCost(v1, v2, v3 *geometry.Vector) float64 {
	plane := FitPlane(v1, v2, v3)
	maxDistance := GetMaxDist(&plane, v1, v2, v3)

	return calc.Amplification * (maxDistance - calc.Threshold)
}

// This function accepts a plane normal vector and an arbitrary number of vectors
// and computes the maximal distance from a point to the plane
func GetMaxDist(plane *geometry.Vector, vectors ...*geometry.Vector) float64 {
	maxDist := 0.0
	for _, vector := range vectors {
		if d := geometry.DistFromPlane(plane, vector); d > maxDist {
			maxDist = d
		}
	}
	return maxDist
}
