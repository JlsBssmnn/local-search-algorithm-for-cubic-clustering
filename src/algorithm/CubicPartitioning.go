package algorithm

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"
	"github.com/go-playground/validator/v10"
)

type CostCalculator[data any] interface {
	TripleCost(d1, d2, d3 data) float64
}

// Represents a partitioning of data points. The indicies of the array
// correspond to the data points and the values of the array to the
// partition the data point belongs to.
type PartitioningArray []int

// Initializes the PartitioningArray into a partitioning of singleton sets. `n` is the
// number of elements that are partitioned. Singleton sets are a starting point for a lot of algorithms.
func (array *PartitioningArray) InitializeSingletonSets(n int) {
	*array = make(PartitioningArray, n)
	for i := range *array {
		(*array)[i] = i
	}
}

// This is the function signiture which every partitioning algorithm should have
type PartitioningAlgorithm[data any] func(input *[]data, calc CostCalculator[data]) PartitioningArray

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

type AllConstraints struct {
	SamePartition      []Edge `json:"same_partition" validate:"required,dive,required"`
	DifferentPartition []Edge `json:"different_partition" validate:"required,dive,required"`
}

type Edge [2]int

type PrecomputedPartitions map[int]*utils.LinkedList[int]

type DisjointSets []int

// Gets the index for the element i and j and checks if the values are correct
func (array *Constraints) getIndex(i, j int) int {
	if array == nil || len(array.array) == 0 || i == j {
		return -1
	}
	if i < 0 || j < 0 || i >= array.numOfElements || j >= array.numOfElements {
		panic("At least one index is out of bounds")
	}
	if j < i {
		i, j = j, i
	}
	matrixRow := (i * (array.numOfElements - 1)) - (i*(i-1))/2
	columnOffset := j - i - 1
	return matrixRow + columnOffset
}

// Returns the value that would be in the ith row and jth column (or vice verca) in the 2D contraint matrix
func (array *Constraints) Get(i, j int) bool {
	if index := array.getIndex(i, j); index == -1 {
		return false
	} else {
		return array.array[index]
	}
}

// Sets the value for element i and j to true
func (array *Constraints) setTrue(i, j int) {
	if index := array.getIndex(i, j); index == -1 {
		panic("Specified indicies for constraints are invalid for this context")
	} else {
		array.array[array.getIndex(i, j)] = true
	}
}

// Gets the representative (the element which has value -1) for the given element
func (disjointSets *DisjointSets) getRepresentative(i int) int {
	value := (*disjointSets)[i]
	for value != -1 {
		i = value
		value = (*disjointSets)[i]
	}
	return i
}

func parseConstraints(path string) *AllConstraints {
	jsonFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteArray, _ := ioutil.ReadAll(jsonFile)
	var allConstraints AllConstraints
	json.Unmarshal(byteArray, &allConstraints)

	validator := validator.New()
	err = validator.Struct(allConstraints)
	if err != nil {
		panic("Contraint file is not correct")
	}

	return &allConstraints
}

func translatetConstraints(allConstraints *AllConstraints, length int) (constraints Constraints, partitions PrecomputedPartitions) {
	constraints = Constraints{array: make([]bool, (length*(length-1))/2), numOfElements: length}
	for _, edge := range allConstraints.DifferentPartition {
		constraints.setTrue(edge[0], edge[1])
	}

	// use a disjoint set data structure to compute the partitions given by the "same_partition" constraints
	disjointSets := make(DisjointSets, length)
	for i := 0; i < length; i++ {
		disjointSets[i] = -1
	}
	for _, edge := range allConstraints.SamePartition {
		representative0 := disjointSets.getRepresentative(edge[0])
		representative1 := disjointSets.getRepresentative(edge[1])
		disjointSets[representative1] = representative0
	}

	// create the precomputed partitions out of the disjoint set datastructure
	partitions = make(PrecomputedPartitions)
	for element, value := range disjointSets {
		representative := disjointSets.getRepresentative(element)
		if value == -1 {
			continue
		}
		if list, ok := partitions[representative]; ok {
			list.Add(element)
		} else {
			partitions[representative] = utils.CreateLinkedList(element)
		}
	}

	return constraints, partitions
}
