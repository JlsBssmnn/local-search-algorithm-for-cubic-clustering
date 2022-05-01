package algorithm

import (
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
)

// A struct for calculating costs.
type CostCalculator struct {
	treshold      float64 // The minimum or maximum value where costs a non-zero
	amplification float64 // A factor that is used in the cost calculation
}

// This function computes the cost for three points that are in the same partition
func (calc *CostCalculator) TripleCost1Part(v1, v2, v3 geometry.Vector) float64 {
	if maxDistance := GetMaxDist(v1, v2, v3); maxDistance <= calc.treshold {
		return 0
	} else {
		return calc.amplification * (maxDistance - calc.treshold)
	}
}

// This function computes the cost for three points that are all in different partitions
func (calc *CostCalculator) TripleCost3Part(v1, v2, v3 geometry.Vector) float64 {
	if maxDistance := GetMaxDist(v1, v2, v3); maxDistance <= calc.treshold {
		return calc.amplification * (calc.treshold - maxDistance)
	} else {
		return 0
	}
}

// This function computes the cost for a given partition with the help of the functions
// TripleCost1Part and TripleCost3Part
func (calc *CostCalculator) CubicPartitionCost(partitioning Partitioning) float64 {
	totalCost := 0.0

	// Compute cost for the case that all elements are in one partition
	for _, partition := range partitioning {
		l := len(partition)
		if l < 3 {
			continue
		}
		for x := 0; x < l; x++ {
			for y := x + 1; y < l; y++ {
				for z := y + 1; z < l; z++ {
					totalCost += calc.TripleCost1Part(partition[x], partition[y], partition[z])
				}
			}
		}
	}

	// Compute cost for the case that all elements are in different partitions
	n := len(partitioning)
	if n < 3 {
		return totalCost
	}
	// iterate over all possible combinations of 3 partitions
	for part1 := 0; part1 < n; part1++ {
		for part2 := part1 + 1; part2 < n; part2++ {
			for part3 := part2 + 1; part3 < n; part3++ {
				// compute the cost for every combination of 3 elements where each element is in one
				// of the 3 partitions
				for x := 0; x < len(partitioning[part1]); x++ {
					for y := 0; y < len(partitioning[part2]); y++ {
						for z := 0; z < len(partitioning[part3]); z++ {
							totalCost += calc.TripleCost3Part(partitioning[part1][x], partitioning[part2][y], partitioning[part3][z])
						}
					}
				}
			}
		}
	}
	return totalCost
}

// This function accepts an arbitrary number of vectors and computes the maximal distance
// when comparing them pairwise
func GetMaxDist(vectors ...geometry.Vector) float64 {
	maxDist := 0.0
	for i, vector := range vectors {
		for j := i + 1; j < len(vectors); j++ {
			if d := geometry.VectorDist(vector, vectors[j]); d > maxDist {
				maxDist = d
			}
		}
	}
	return maxDist
}
