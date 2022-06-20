package algorithm

type CostCalculator[data any] interface {
	TripleCost(d1, d2, d3 data) float64
}

// Represents a partitioning of data points. The indicies of the array
// correspond to the data points and the values of the array to the
// partition the data point belongs to.
type PartitioningArray []int

// This datastructure can store whether two elements should be in different partitions.
// The greedy moving algorithm can use it during it's execution to ensure that the
// resulting partitioning satisfies the constraints.
//
// The contraints can be seen as a 2D symmetric matrix where the values in the matrix are
// booleans. If a bool is true, then the column index and row index seen as indicies for
// the input elements cannot be in the same partition. This data structure stores
// not the entire matrix but only one triangle to be as space efficiant as possible.
type Constraints struct {
	array         []bool
	numOfElements int
}

// Returns the value that would be in the ith row and jth column (or vice verca) in the 2D contraint matrix
func (array *Constraints) Get(i, j int) bool {
	if array == nil || len(array.array) == 0 || i == j {
		return false
	}
	if i < 0 || j < 0 || i >= array.numOfElements || j >= array.numOfElements {
		panic("At least one index is out of bounds")
	}
	if j < i {
		i, j = j, i
	}
	matrixRow := (i * (array.numOfElements - 1)) - (i*(i-1))/2
	columnOffset := j - i - 1
	return array.array[matrixRow+columnOffset]
}
