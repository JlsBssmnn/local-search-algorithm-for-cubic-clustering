package algorithm

import (
	"math"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"
)

type NaiveGreedyJoiningAlgorithm[data any] struct {
	input         *[]data
	calc          CostCalculator[data]
	partitioning  PartitioningArray
	partitions    map[int][]int
	partitionList []int
}

func NaiveGreedyJoining[data any](input *[]data, calc CostCalculator[data]) PartitioningArray {
	n := len(*input)

	algorithm := NaiveGreedyJoiningAlgorithm[data]{input: input, calc: calc}
	algorithm.partitioning.InitializeSingletonSets(n)

	algorithm.partitions = make(map[int][]int, n)
	algorithm.partitionList = make([]int, n)

	for i := 0; i < n; i++ {
		algorithm.partitions[i] = []int{i}
		algorithm.partitionList[i] = i
	}

	nextJoin, costDiff := algorithm.FindBestJoin()
	for costDiff < 0 {
		algorithm.join(nextJoin[0], nextJoin[1])
		nextJoin, costDiff = algorithm.FindBestJoin()
	}

	return algorithm.partitioning
}

func (algorithm *NaiveGreedyJoiningAlgorithm[data]) FindBestJoin() (join [2]int, costDiff float64) {
	join = [2]int{-1, -1}
	costDiff = math.Inf(1)
	n := len(algorithm.partitionList)

	for i, partition1 := range algorithm.partitionList {
		for j := i + 1; j < n; j++ {
			partition2 := algorithm.partitionList[j]
			var currentDiff float64

			if len(algorithm.partitions[partition1]) == 1 && len(algorithm.partitions[partition2]) == 1 {
				for k := j + 1; k < n; k++ {
					partition3 := algorithm.partitionList[k]
					currentDiff = algorithm.computeCostDiff3Part(partition1, partition2, partition3)
					if currentDiff < costDiff {
						costDiff = currentDiff
						join[0] = partition1
						join[1] = partition2
					}
				}
			} else {
				currentDiff = algorithm.computeCostDiff(partition1, partition2)
				if currentDiff < costDiff {
					costDiff = currentDiff
					join[0] = partition1
					join[1] = partition2
				}
			}
		}
	}

	return join, costDiff
}

// Computes the difference in cost when joining the 2 given partitions
func (algorithm *NaiveGreedyJoiningAlgorithm[data]) computeCostDiff(part1Idx, part2Idx int) (costDiff float64) {
	part1 := algorithm.partitions[part1Idx]
	part2 := algorithm.partitions[part2Idx]
	n1 := len(part1)
	n2 := len(part2)

	// all cases where 2 elements are in part1 and 1 element is in part2
	for i := 0; i < n1-1; i++ {
		for j := i + 1; j < n1; j++ {
			for k := 0; k < n2; k++ {
				costDiff += algorithm.calc.TripleCost((*algorithm.input)[part1[i]], (*algorithm.input)[part1[j]], (*algorithm.input)[part2[k]])
			}
		}
	}

	// all cases where 1 element is in part1 and 2 elements are in part2
	for i := 0; i < n2-1; i++ {
		for j := i + 1; j < n2; j++ {
			for k := 0; k < n1; k++ {
				costDiff += algorithm.calc.TripleCost((*algorithm.input)[part2[i]], (*algorithm.input)[part2[j]], (*algorithm.input)[part1[k]])
			}
		}
	}

	return costDiff
}

// Computes the difference in cost when joining the two singleton partitions (first 2 parameters)
// with the third one (third parameter)
func (algorithm *NaiveGreedyJoiningAlgorithm[data]) computeCostDiff3Part(part1Idx, part2Idx, part3Idx int) (costDiff float64) {
	if len(algorithm.partitions[part1Idx]) != 1 || len(algorithm.partitions[part1Idx]) != 1 {
		panic("One of the first 2 partitions was not a singleton set")
	}
	elem1 := algorithm.partitions[part1Idx][0]
	elem2 := algorithm.partitions[part2Idx][0]

	for _, elem3 := range algorithm.partitions[part3Idx] {
		costDiff += algorithm.calc.TripleCost((*algorithm.input)[elem1], (*algorithm.input)[elem2], (*algorithm.input)[elem3])
	}

	return costDiff
}

// Adjusts the data structures for the join of the 2 given partitions
func (algorithm *NaiveGreedyJoiningAlgorithm[data]) join(part1, part2 int) {
	utils.SortInts(&part1, &part2)

	// update partitioning array
	for i, part := range algorithm.partitioning {
		if part == part2 {
			algorithm.partitioning[i] = part1
		}
	}

	// update partitions
	algorithm.partitions[part1] = append(algorithm.partitions[part1], algorithm.partitions[part2]...)
	delete(algorithm.partitions, part2)

	// update partitionList
	i := utils.Find(algorithm.partitionList, part2)
	n := len(algorithm.partitionList)
	algorithm.partitionList[i] = algorithm.partitionList[n-1]
	algorithm.partitionList = algorithm.partitionList[:n-1]
}
