package algorithm

import (
	"math"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"
)

type NaiveGreedyMovingAlgorithm[data any] struct {
	input        *[]data
	calc         CostCalculator[data]
	partitioning PartitioningArray
	partitions   map[int]*[]int
}

func (algorithm *NaiveGreedyMovingAlgorithm[data]) FindBestMove() ([3]int, float64) {
	// find the best move of element a (and potentially b) to a set U
	minCostDif := math.Inf(1)
	a := -1
	b := -1
	U := -1

	// iterate over input data
	for i := range *algorithm.input {
		// iterate over partitions
		for j := range algorithm.partitions {
			// if the element is in the partition we don't consider the move
			if algorithm.partitioning[i] == j {
				continue
			}
			// if the partition contains just 1 element, we have to consider an additional move
			// s.t. we can calculate triple costs
			if len(*algorithm.partitions[j]) == 1 {
				// iterate over all input elements after the ith element
				for k := i + 1; k < len(*algorithm.input); k++ {
					if algorithm.partitioning[k] == j {
						continue
					}
					moveCost := algorithm.costDiff2Moves(j, i, k)
					if moveCost < minCostDif {
						minCostDif = moveCost
						a = i
						b = k
						U = j
					}
				}
			} else {
				moveCost := algorithm.costDiff1Move(j, i)
				if moveCost < minCostDif {
					minCostDif = moveCost
					a = i
					b = -1
					U = j
				}
			}
		}
		// if the partition of i contains more than 1 element we consider the case
		// where we remove i from it's partition
		if len(*algorithm.partitions[algorithm.partitioning[i]]) > 1 {
			moveCost := algorithm.costDiffRemoveElement(i)
			if moveCost < minCostDif {
				minCostDif = moveCost
				a = i
				b = -1
				U = -1
			}
		}
	}
	if a != -1 && U == -1 {
		// find maximum partition to create a new one that is one higher
		maxPart := -1
		for part := range algorithm.partitions {
			if part > maxPart {
				maxPart = part
			}
		}
		U = maxPart + 1
	}
	return [3]int{U, a, b}, minCostDif
}

// Initializes the partitioning array and the partitions map
func (algorithm *NaiveGreedyMovingAlgorithm[data]) initialize() {
	algorithm.partitioning = make(PartitioningArray, len(*algorithm.input))
	algorithm.partitions = make(map[int]*[]int)

	for i := 0; i < len(*algorithm.input); i++ {
		algorithm.partitioning[i] = i
		algorithm.partitions[i] = &[]int{i}
	}
}

// Moves element a into partition U
func (algorithm *NaiveGreedyMovingAlgorithm[data]) moveElement(U, a int) {
	partA := algorithm.partitions[algorithm.partitioning[a]]
	indexA := utils.Find(*partA, a)

	if len(*partA) == 1 {
		delete(algorithm.partitions, algorithm.partitioning[a])
	} else {
		*partA = append((*partA)[:indexA], (*partA)[indexA+1:]...)
	}
	partU, ok := algorithm.partitions[U]
	if !ok {
		algorithm.partitions[U] = &[]int{a}
	} else {
		*partU = append(*partU, a)
	}
	algorithm.partitioning[a] = U
}

// Computes the difference in costs for the partitioning if the given element a is moved
// to partition U
func (algorithm *NaiveGreedyMovingAlgorithm[data]) costDiff1Move(U, a int) float64 {
	cost := 0.0
	cost += algorithm.costDiffRemoveElement(a)

	partU := *algorithm.partitions[U]

	for i := 0; i < len(partU); i++ {
		for j := i + 1; j < len(partU); j++ {
			cost += algorithm.calc.TripleCost((*algorithm.input)[partU[i]], (*algorithm.input)[partU[j]], (*algorithm.input)[a])
		}
	}
	return cost
}

// Computes the cost difference when moving the elements a1 and a2 into partition U
func (algorithm *NaiveGreedyMovingAlgorithm[data]) costDiff2Moves(U, a1, a2 int) float64 {
	cost := 0.0
	partU := *algorithm.partitions[U]

	// a1 and a2 are in different partitions
	if algorithm.partitioning[a1] != algorithm.partitioning[a2] {
		cost += algorithm.costDiff1Move(U, a1)
		cost += algorithm.costDiff1Move(U, a2)
		for i := 0; i < len(partU); i++ {
			cost += algorithm.calc.TripleCost((*algorithm.input)[partU[i]], (*algorithm.input)[a1], (*algorithm.input)[a2])
		}
		return cost
	}

	partA := *algorithm.partitions[algorithm.partitioning[a1]]
	for i := 0; i < len(partA); i++ {
		if a1 == partA[i] || a2 == partA[i] {
			continue
		}
		cost -= algorithm.calc.TripleCost((*algorithm.input)[partA[i]], (*algorithm.input)[a1], (*algorithm.input)[a2])
		for j := i + 1; j < len(partA); j++ {
			if a1 == partA[j] || a2 == partA[j] {
				continue
			}
			cost -= algorithm.calc.TripleCost((*algorithm.input)[partA[i]], (*algorithm.input)[partA[j]], (*algorithm.input)[a1])
			cost -= algorithm.calc.TripleCost((*algorithm.input)[partA[i]], (*algorithm.input)[partA[j]], (*algorithm.input)[a2])
		}
	}
	for i := 0; i < len(partU); i++ {
		cost += algorithm.calc.TripleCost((*algorithm.input)[partU[i]], (*algorithm.input)[a1], (*algorithm.input)[a2])
		for j := i + 1; j < len(partU); j++ {
			cost += algorithm.calc.TripleCost((*algorithm.input)[partU[i]], (*algorithm.input)[partU[j]], (*algorithm.input)[a1])
			cost += algorithm.calc.TripleCost((*algorithm.input)[partU[i]], (*algorithm.input)[partU[j]], (*algorithm.input)[a2])
		}
	}
	return cost
}

// Computes the cost difference when moving the element a into a new partition
func (algorithm *NaiveGreedyMovingAlgorithm[data]) costDiffRemoveElement(a int) float64 {
	cost := 0.0
	partA := *algorithm.partitions[algorithm.partitioning[a]]
	for i := 0; i < len(partA)-1; i++ {
		if a == partA[i] {
			continue
		}
		for j := i + 1; j < len(partA); j++ {
			if a == partA[j] {
				continue
			}
			cost -= algorithm.calc.TripleCost((*algorithm.input)[partA[i]], (*algorithm.input)[partA[j]], (*algorithm.input)[a])
		}
	}
	return cost
}

func GreedyMoving[data any](input *[]data, calc CostCalculator[data]) PartitioningArray {
	algorithm := NaiveGreedyMovingAlgorithm[data]{input: input, calc: calc}
	algorithm.initialize()

	nextMove, cost := algorithm.FindBestMove()
	U, a, b := nextMove[0], nextMove[1], nextMove[2]

	for cost < 0 && U != -1 && a != -1 {
		algorithm.moveElement(U, a)
		if b != -1 {
			algorithm.moveElement(U, b)
		}
		nextMove, cost = algorithm.FindBestMove()
		U, a, b = nextMove[0], nextMove[1], nextMove[2]
	}
	return algorithm.partitioning
}
