package algorithm

import (
	"math"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"
)

// The greedy joining algorithm with following properties:
// 	- it will evaluation joins of 3 partitions
//	- it joins only 2 partitions, based on which join yields the best costs
//	- if there are less than 3 partitions left, the algorithm will terminate
func GeedyJoiningV1[data any](input []data, calc CostCalculator[data]) Partitioning[data] {
	// the initial partitioning starts with singleton sets
	var part Partitioning[data] = [][]data{}
	for _, point := range input {
		part = append(part, []data{point})
	}
	cost := CubicPartitionCost(part, calc)

	// find {B,C,D} that minimize the cost when joining them
	for {
		nextCost := math.Inf(1) // costs of the best rated join of 2 partitions
		bestCost := math.Inf(1) // best observed costs when joining 2 or 3 partitions
		B, C := -1, -1

		for i := 0; i < len(part); i++ {
			for j := i + 1; j < len(part); j++ {
				for k := j + 1; k < len(part); k++ {
					newPart := joinPartitions(part, i, j, k)
					joinCost := CubicPartitionCost(newPart, calc)

					// the join is better than what we've previously found
					if joinCost < bestCost {
						// if all partitions have 1 element, join 2 of them randomly
						if len(part[i]) == 1 && len(part[j]) == 1 && len(part[k]) == 1 {
							excludePart := utils.RandomInt(0, 3)
							indeces := []int{i, j, k}
							indeces = append(indeces[:excludePart], indeces[excludePart+1:]...)
							B = indeces[0]
							C = indeces[1]
							nextCost = cost
						} else { // otherwise find the best join between the partitions
							part1 := joinPartitions(part, i, j)
							part2 := joinPartitions(part, i, k)
							part3 := joinPartitions(part, j, k)

							cost1 := CubicPartitionCost(part1, calc)
							cost2 := CubicPartitionCost(part2, calc)
							cost3 := CubicPartitionCost(part3, calc)

							best := utils.ArgMin([]float64{cost1, cost2, cost3})
							switch best {
							case 0:
								B = i
								C = j
								nextCost = cost1
							case 1:
								B = i
								C = k
								nextCost = cost2
							case 2:
								B = j
								C = k
								nextCost = cost3
							default:
								panic("Something went wrong")
							}
						}
						bestCost = joinCost
					}
				}
			}
		}

		// There aren't enough partitions to find a join so just
		// return the best found partitioning
		if B == -1 || C == -1 {
			return part
		}

		if bestCost-cost < 0 {
			cost = nextCost
			part = joinPartitions(part, B, C)
			continue
		} else {
			break
		}
	}

	return part
}

// The greedy joining algorithm with following properties:
// 	- it will only evaluate joins of 3 partitions if the first 2 contain one element
// 	- it only joins 2 partitions
// 	- if there is only one partition left the algorithm terminates
func GeedyJoiningV2[data any](input []data, calc CostCalculator[data]) Partitioning[data] {
	// the initial partitioning starts with singleton sets
	var part Partitioning[data] = [][]data{}
	for _, point := range input {
		part = append(part, []data{point})
	}
	cost := CubicPartitionCost(part, calc)

	// find {B,C} or {B,C,D} that minimize the cost when joining them
	for {
		nextCost := math.Inf(1) // costs for the next partition which arises from joining 2 partitions
		bestCost := math.Inf(1) // the best found cost for joining either 2 or 3 partitions
		B, C := -1, -1

		for i := 0; i < len(part); i++ {
			for j := i + 1; j < len(part); j++ {
				// both partitions have one element, so we must consider 3 partitions
				// for the join
				if len(part[i]) == 1 && len(part[j]) == 1 {
					for k := j + 1; k < len(part); k++ {
						newPart := joinPartitions(part, i, j, k)
						joinCost := CubicPartitionCost(newPart, calc)

						if joinCost < bestCost {
							B = i
							C = j
							bestCost = joinCost
							nextCost = cost
						}
					}
				} else { // else both partitions have in sum at least 3 elements
					newPart := joinPartitions(part, i, j)
					joinCost := CubicPartitionCost(newPart, calc)

					if joinCost < bestCost {
						B = i
						C = j
						bestCost = joinCost
						nextCost = joinCost
					}
				}
			}
		}

		// There aren't enough partitions to find a join so just
		// return the best found partitioning
		if B == -1 || C == -1 {
			return part
		}

		if bestCost-cost < 0 {
			cost = nextCost
			part = joinPartitions(part, B, C)
			continue
		} else {
			break
		}
	}

	return part
}

// This function takes a partitioning and joins all the partitions at the given indicies
func joinPartitions[data any](partitioning Partitioning[data], partitions ...int) Partitioning[data] {
	join := []data{}
	for _, index := range partitions {
		join = append(join, partitioning[index]...)
	}

	newPartitioning := [][]data{}
	for i, partition := range partitioning {
		if !utils.Contains(partitions, i) {
			newPartitioning = append(newPartitioning, partition)
		}
	}
	newPartitioning = append(newPartitioning, join)

	return newPartitioning
}
