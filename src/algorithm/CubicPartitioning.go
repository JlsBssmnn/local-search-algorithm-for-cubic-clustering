package algorithm

type CostCalculator[data any] interface {
	TripleCost(d1, d2, d3 data) float64
}

// Represents a partitioning of data points. The indicies of the array
// correspond to the data points and the values of the array to the
// partition the data point belongs to.
type PartitioningArray []int

// Represents a partitioning of data points. First dimension are the
// different partitions, second dimension are the elements in one partition.
// Example partitioning: [[a,b,c], [d,e,f,g,h], [i,j]]
type Partitioning[data any] [][]data

// This function computes the cost for a given partition with the help of the functions
// TripleCost1Part and TripleCost3Part
func CubicPartitionCost[data any](partitioning Partitioning[data], calc CostCalculator[data]) float64 {
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
					totalCost += calc.TripleCost(partition[x], partition[y], partition[z])
				}
			}
		}
	}

	return totalCost
}
