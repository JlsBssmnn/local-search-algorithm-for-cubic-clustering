package algorithm

type CostCalculator[data any] interface {
	TripleCost(d1, d2, d3 data) float64
}

// Represents a partitioning of data points. The indicies of the array
// correspond to the data points and the values of the array to the
// partition the data point belongs to.
type PartitioningArray []int
