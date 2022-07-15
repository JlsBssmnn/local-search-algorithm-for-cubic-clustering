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
	tripleCosts  *TripleCosts
}

func NaiveGreedyMoving[data any](input *[]data, calc CostCalculator[data]) PartitioningArray {
	algorithm := NaiveGreedyMovingAlgorithm[data]{input: input, calc: calc}
	algorithm.initialize()

	nextMove, costDiff := algorithm.findBestMove()
	U, a, b := nextMove[0], nextMove[1], nextMove[2]

	for costDiff < 0 {
		algorithm.moveElement(U, a)
		if b != -1 {
			algorithm.moveElement(U, b)
		}
		nextMove, costDiff = algorithm.findBestMove()
		U, a, b = nextMove[0], nextMove[1], nextMove[2]
	}
	return algorithm.partitioning
}

func (algorithm *NaiveGreedyMovingAlgorithm[data]) findBestMove() (bestMove [3]int, minCostDiff float64) {
	// find the best move of element a (and potentially b) to a set U
	minCostDiff = math.Inf(1)
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
					currentDiff := algorithm.costDiff2Moves(j, i, k)
					if currentDiff < minCostDiff {
						minCostDiff = currentDiff
						a = i
						b = k
						U = j
					}
				}
			} else {
				currentDiff := algorithm.costDiff1Move(j, i)
				if currentDiff < minCostDiff {
					minCostDiff = currentDiff
					a = i
					b = -1
					U = j
				}
			}
		}
		// if the partition of i contains more than 1 element we consider the case
		// where we remove i from it's partition
		if len(*algorithm.partitions[algorithm.partitioning[i]]) > 1 {
			currentDiff := algorithm.costDiffRemoveElement(i)
			if currentDiff < minCostDiff {
				minCostDiff = currentDiff
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

	bestMove = [3]int{U, a, b}
	return bestMove, minCostDiff
}

// Initializes the partitioning array and the partitions map
func (algorithm *NaiveGreedyMovingAlgorithm[data]) initialize() {
	algorithm.partitioning = make(PartitioningArray, len(*algorithm.input))
	algorithm.partitions = make(map[int]*[]int)

	for i := 0; i < len(*algorithm.input); i++ {
		algorithm.partitioning[i] = i
		algorithm.partitions[i] = &[]int{i}
	}
	algorithm.InitializeTripleCosts()
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
			cost += algorithm.tripleCosts.GetTripleCost(partU[i], partU[j], a)
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
			cost += algorithm.tripleCosts.GetTripleCost(partU[i], a1, a2)
		}
		return cost
	}

	partA := *algorithm.partitions[algorithm.partitioning[a1]]
	for i := 0; i < len(partA); i++ {
		if a1 == partA[i] || a2 == partA[i] {
			continue
		}
		cost -= algorithm.tripleCosts.GetTripleCost(partA[i], a1, a2)
		for j := i + 1; j < len(partA); j++ {
			if a1 == partA[j] || a2 == partA[j] {
				continue
			}
			cost -= algorithm.tripleCosts.GetTripleCost(partA[i], partA[j], a1)
			cost -= algorithm.tripleCosts.GetTripleCost(partA[i], partA[j], a2)
		}
	}
	for i := 0; i < len(partU); i++ {
		cost += algorithm.tripleCosts.GetTripleCost(partU[i], a1, a2)
		for j := i + 1; j < len(partU); j++ {
			cost += algorithm.tripleCosts.GetTripleCost(partU[i], partU[j], a1)
			cost += algorithm.tripleCosts.GetTripleCost(partU[i], partU[j], a2)
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
			cost -= algorithm.tripleCosts.GetTripleCost(partA[i], partA[j], a)
		}
	}
	return cost
}

func (algorithm *NaiveGreedyMovingAlgorithm[data]) InitializeTripleCosts() {
	n := len(*algorithm.input)
	firstDim := make(TripleCosts, n-2)

	for i := 0; i < n-2; i++ {
		secondDim := make([][]float64, n-i-2)

		for j := i + 1; j < n-1; j++ {
			thirdDim := make([]float64, n-j-1)

			for k := j + 1; k < n; k++ {
				thirdDim[k-j-1] = algorithm.calc.TripleCost((*algorithm.input)[i], (*algorithm.input)[j], (*algorithm.input)[k])
			}
			secondDim[j-i-1] = thirdDim
		}
		firstDim[i] = secondDim
	}
	algorithm.tripleCosts = &firstDim
}
