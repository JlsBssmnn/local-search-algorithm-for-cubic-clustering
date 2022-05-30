package algorithm

import (
	"errors"
	"fmt"
	"math"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"
)

// This struct holds all the information that is necessary to perform the greedy
// joining algorithm
type GreedyJoiningAlgorithm[data any] struct {
	input        *[]data
	calc         CostCalculator[data]
	partitioning PartitioningArray
	costs        *Costs
}

// A data structure which stores the costs that were calculated for the greedy joining
// algorithm such that they can be reused again without recomputation
type Costs []*OnePartitionCosts

type OnePartitionCosts struct {
	minCost            float64
	bestJoin           int
	twoPartitionsCosts []*TwoPartitionsCosts
}

type TwoPartitionsCosts struct {
	joinCost        float64
	bestJoin        int
	tripleJoinCosts *[]float64
}

// -------------------------- Methods for the cost data structure and the GreedyJoiningAlgorithm struct

// Verify that the given indicies can be used for operations on the
// cost data structure
func (costs *Costs) verifyIndicies(indicies ...*int) {
	if !utils.AllDifferent(utils.Map(indicies, func(element *int) int {
		return *element
	})) {
		panic("You input one partition more than once")
	} else if utils.Any(indicies, func(x *int) bool { return *x < 0 }) {
		panic("Indicies to the partitions must not be negative")
	} else if utils.Any(indicies, func(x *int) bool { return *x > len(*costs) }) {
		panic("Partition indicies are out of bounds")
	}
	utils.SortInts(indicies...)
}

// Extracts the join cost when joining the 3 partitions at the three given indicies.
// If this cost is not stored in the data stucture an error will be returned.
func (costs *Costs) TripleJoinCost(i, j, k int) (float64, error) {
	costs.verifyIndicies(&i, &j, &k)

	if triples := (*((*costs)[i]).twoPartitionsCosts[j-i-1]).tripleJoinCosts; triples != nil {
		return (*triples)[k-j-1], nil
	} else {
		return 0, errors.New("The requested triple cost is not stored in the data structure")
	}
}

// Extracts the join cost when joining the 2 partitions at the two given indicies.
// It also returns whether these costs are future costs (so both partitions contain one
// element and the costs correspond to a join of 3 partitions) where the second return result
// will be true or if they are the real costs for joining the two paritions (return will be false)
func (costs *Costs) JoinCost(i, j int) (float64, bool) {
	costs.verifyIndicies(&i, &j)

	join := (*((*costs)[i]).twoPartitionsCosts[j-i-1])

	return join.joinCost, join.tripleJoinCosts != nil
}

func (costs *Costs) RealJoinCost(i, j int) float64 {
	cost, future := (*costs).JoinCost(i, j)
	if future {
		return 0
	} else {
		return cost
	}
}

// This function computes the index either in the second dimension
// where the cost of joining partition i and j is located or in the
// third dimension where the cost of joining some partition with i
// and j is located
func (costs *Costs) GetIndex(i, j int) int {
	costs.verifyIndicies(&i, &j)
	return j - i - 1
}

// Deletes the costs for the jth partition in the second dimension for
// partition i
func (costs *Costs) Delete2D(i, j int) {
	index := costs.GetIndex(i, j)
	(*(*costs)[i]).twoPartitionsCosts = append((*(*costs)[i]).twoPartitionsCosts[:index], (*(*costs)[i]).twoPartitionsCosts[index+1:]...)
}

// Deletes the costs for the kth partition in the third dimension for
// the join of partition i and j
func (costs *Costs) Delete3D(i, j, k int) {
	index2D := costs.GetIndex(i, j)
	index3D := costs.GetIndex(j, k)
	twoPartitionsCosts := (*(*costs)[i]).twoPartitionsCosts[index2D]
	*twoPartitionsCosts.tripleJoinCosts = append((*twoPartitionsCosts.tripleJoinCosts)[:index3D], (*twoPartitionsCosts.tripleJoinCosts)[index3D+1:]...)

	// If we deleted the triple join cost that were also the best cost, we have to
	// determinate the new minimum and store it
	if twoPartitionsCosts.bestJoin == index3D {
		joinCost, bestJoin := utils.MinAndArgMin(*twoPartitionsCosts.tripleJoinCosts)
		if bestJoin == -1 {
			// In this case the list of triple joins is empty
			twoPartitionsCosts.bestJoin = -1
			twoPartitionsCosts.joinCost = math.Inf(1)
		} else {
			twoPartitionsCosts.joinCost = joinCost
			twoPartitionsCosts.bestJoin = bestJoin
		}
	}
}

// Calculates the sum over all triples cost where each element in the triple is in
// either partition i, j or k (triples that are exclusive to this combination of partitions)
func (algorithm *GreedyJoiningAlgorithm[data]) TripleCost3Part(i, j, k int) float64 {
	cost := 0.0
	iPartition := utils.LinkedList[int]{}
	jPartition := utils.LinkedList[int]{}
	kPartition := utils.LinkedList[int]{}

	for index, partition := range algorithm.partitioning {
		switch partition {
		case i:
			iPartition.Add(index)
		case j:
			jPartition.Add(index)
		case k:
			kPartition.Add(index)
		}
	}

	iIter := iPartition.Iterator()
	for iIter.HasNext() {
		iElement := iIter.Next()
		jIter := jPartition.Iterator()
		for jIter.HasNext() {
			jElement := jIter.Next()
			kIter := kPartition.Iterator()
			for kIter.HasNext() {
				kElement := kIter.Next()
				cost += algorithm.calc.TripleCost(
					(*algorithm.input)[iElement],
					(*algorithm.input)[jElement],
					(*algorithm.input)[kElement])
			}
		}
	}
	return cost
}

// Recompute the minimum over the second dimension of partition i.
// It also checks if this cost is better than the cost stored in
func (costs *Costs) Min2D(i int, bestJoinOverall *[2]int, bestJoinCostOverall *float64) {
	// recompute the minimum over the second dimension of i
	onePartitionCosts := (*costs)[i]
	min, argmin := utils.MinAndArgMin(utils.Map(onePartitionCosts.twoPartitionsCosts, func(element *TwoPartitionsCosts) float64 {
		return element.joinCost
	}))
	onePartitionCosts.bestJoin = argmin
	onePartitionCosts.minCost = min

	// if these costs are better than the previous best cost over the first dimension we store them
	if min < *bestJoinCostOverall {
		*bestJoinCostOverall = min
		(*bestJoinOverall)[0] = i
		(*bestJoinOverall)[1] = i + argmin + 1
	}
}

// Updates the partitioning array to the join of the two partitions part1 and part2
func (algorithm *GreedyJoiningAlgorithm[data]) updatePartitioningArray(part1, part2 int) {
	for i, part := range algorithm.partitioning {
		if part == part2 {
			algorithm.partitioning[i] = part1
		} else if part > part2 {
			algorithm.partitioning[i]--
		}
	}
}

// --------------------------

// This function can be used to get a GreedyJoiningAlgorithm struct outside of this package,
// this should mainly be used for testing
func GetAlgorithm[T any](input *[]T, calc CostCalculator[T]) GreedyJoiningAlgorithm[T] {
	return GreedyJoiningAlgorithm[T]{input: input, calc: calc}
}

// Sets up the costs and the partitioning into singleton sets of the algorithm
func (algorithm *GreedyJoiningAlgorithm[data]) InitializeAlgorithm() ([2]int, float64) {
	algorithm.partitioning = make(PartitioningArray, len(*algorithm.input))
	for i := range algorithm.partitioning {
		algorithm.partitioning[i] = i
	}
	costs, bestJoinOverall, bestJoinCostOverall := InitializeCosts(algorithm.input, algorithm.calc)
	algorithm.costs = &costs
	return bestJoinOverall, bestJoinCostOverall
}

// Initializes the cost data structure for singleton sets of the input data
// It returns this data structure and two indicies of partitions which have the best
// join cost as well as the cost.
func InitializeCosts[data any](input *[]data, calc CostCalculator[data]) (Costs, [2]int, float64) {
	size := len(*input)
	costs := make(Costs, size-1)
	bestJoinOverall := [2]int{-1, -1}
	bestJoinCostOverall := math.Inf(1)

	for i := 0; i < len(costs); i++ {
		onePartitionCosts := make([]*TwoPartitionsCosts, size-i-1)
		minCost2Dim := math.Inf(1)
		bestJoin := -1
		for j := 0; j < len(onePartitionCosts); j++ {
			tripleJoinCosts := make([]float64, size-(i+j+2))
			minCost3Dim := math.Inf(1)
			bestJoin3Dim := -1
			for k := 0; k < len(tripleJoinCosts); k++ {
				tripleCost := calc.TripleCost((*input)[i], (*input)[i+j+1], (*input)[i+j+k+2])
				if tripleCost < minCost3Dim {
					minCost3Dim = tripleCost
					bestJoin3Dim = k
				}
				tripleJoinCosts[k] = tripleCost
			}
			onePartitionCosts[j] = &TwoPartitionsCosts{
				tripleJoinCosts: &tripleJoinCosts,
				joinCost:        minCost3Dim,
				bestJoin:        bestJoin3Dim,
			}
			if minCost3Dim < minCost2Dim {
				minCost2Dim = minCost3Dim
				bestJoin = j
			}
		}
		costs[i] = &OnePartitionCosts{
			minCost:            minCost2Dim,
			bestJoin:           bestJoin,
			twoPartitionsCosts: onePartitionCosts,
		}
		if minCost2Dim < bestJoinCostOverall {
			bestJoinCostOverall = minCost2Dim
			bestJoinOverall[0] = i
			bestJoinOverall[1] = i + bestJoin + 1
		}
	}

	return costs, bestJoinOverall, bestJoinCostOverall
}

// This function adjusts the cost data structure to the new partitioning
// which arises when joining the two partitions of which the indicies are
// given to this function. It returns the indicies of two partitions
// which have the best join cost in the new data structure and the cost.
func (algorithm *GreedyJoiningAlgorithm[data]) Join(part1, part2 int) ([2]int, float64) {
	algorithm.costs.verifyIndicies(&part1, &part2)

	previousJoinCost := algorithm.costs.RealJoinCost(part1, part2)

	bestJoinOverall := [2]int{-1, -1}
	bestJoinCostOverall := math.Inf(1)

	algorithm.joinStep1(part1, part2, previousJoinCost, &bestJoinOverall, &bestJoinCostOverall)

	// if part1 is the second last partition it can just be removed because after the join it'll be the last
	// partition which is not stored in the data structure
	if part1 == len(*algorithm.costs)-1 {
		*algorithm.costs = (*algorithm.costs)[:part1]
		algorithm.updatePartitioningArray(part1, part2)
		return bestJoinOverall, bestJoinCostOverall
	}
	algorithm.joinStep2(part1, part2, previousJoinCost, &bestJoinOverall, &bestJoinCostOverall)
	algorithm.joinStep3(part1, part2, previousJoinCost, &bestJoinOverall, &bestJoinCostOverall)
	algorithm.joinStep4(part1, part2, previousJoinCost, &bestJoinOverall, &bestJoinCostOverall)

	algorithm.updatePartitioningArray(part1, part2)

	return bestJoinOverall, bestJoinCostOverall
}

func (algorithm *GreedyJoiningAlgorithm[data]) joinStep1(part1, part2 int, previousJoinCost float64,
	bestJoinOverall *[2]int, bestJoinCostOverall *float64) {
	// loop through first dimension until part1
	for i := 0; i < part1; i++ {

		// determine the index where part1 is located in the second dimension
		index2DPart1 := algorithm.costs.GetIndex(i, part1)

		// loop through second dimension until the element part1
		for j := 0; j < index2DPart1; j++ {
			twoPartitionsCosts := ((*(algorithm.costs))[i]).twoPartitionsCosts[j]
			if triples := twoPartitionsCosts.tripleJoinCosts; triples != nil {
				// This is the case where the elements in dimension 1 and 2 don't contain at least 3 elements together

				// determine the index where part1 is located in the third dimension
				index3DPart1 := algorithm.costs.GetIndex(i+j+1, part1)

				// the current partition in the second dimension
				secondElement := i + j + 1

				// all triple costs that will be added to get the new cost
				triplesToConsider := [4][3]int{{i, secondElement, part1}, {i, secondElement, part2},
					{i, part1, part2}, {secondElement, part1, part2}}
				// all tuple costs which will be substracted to get the new cost
				tuplesToConsider := [4][2]int{{i, part1}, {i, part2}, {secondElement, part1}, {secondElement, part2}}

				// update the cost for join with part1
				newTripleCost := utils.MapSum(triplesToConsider[:],
					func(element [3]int) float64 {
						cost, _ := algorithm.costs.TripleJoinCost(element[0], element[1], element[2])
						return cost
					},
				) - utils.MapSum(tuplesToConsider[:],
					func(element [2]int) float64 {
						return algorithm.costs.RealJoinCost(element[0], element[1])
					},
				) - (2 * previousJoinCost)
				(*triples)[index3DPart1] = newTripleCost

				// delete cost for join with part2 in third dimension
				algorithm.costs.Delete3D(i, secondElement, part2)

				// check if new triple join is better than previous minimum
				if newTripleCost < twoPartitionsCosts.joinCost {
					twoPartitionsCosts.joinCost = newTripleCost
					twoPartitionsCosts.bestJoin = index3DPart1
				} else if twoPartitionsCosts.bestJoin == index3DPart1 {
					// in this case the previous minimum involed either part1 or part2, but since they have been joined,
					// the new minimum must be found
					joinCost, bestJoin := utils.MinAndArgMin(*triples)
					twoPartitionsCosts.joinCost = joinCost
					twoPartitionsCosts.bestJoin = bestJoin
				}
			}
		}
		{
			// adjust the join of element i with part1
			twoPartitionsCosts := ((*(algorithm.costs))[i]).twoPartitionsCosts[index2DPart1]
			if twoPartitionsCosts.tripleJoinCosts != nil {
				// if i and part1 both have 1 element, then the new join cost has already been calculated
				twoPartitionsCosts.joinCost = (*twoPartitionsCosts.tripleJoinCosts)[algorithm.costs.GetIndex(part1, part2)] - previousJoinCost
				twoPartitionsCosts.tripleJoinCosts = nil
			} else {
				// the new cost require some additional calculations: the cost of triples where each element is in
				// i, part1 and part2 respectively
				additionalCost := algorithm.TripleCost3Part(i, part1, part2)
				twoPartitionsCosts.joinCost = algorithm.costs.RealJoinCost(i, part1) + algorithm.costs.RealJoinCost(i, part2) + additionalCost
			}
		}
		// determine the index where part2 is located in the second dimension
		index2DPart2 := part2 - i - 1

		// delete part2 in every triple cost list
		for j := index2DPart1 + 1; j < index2DPart2; j++ {
			twoPartitionsCost := (*algorithm.costs)[i].twoPartitionsCosts[j]
			tripleJoinCosts := twoPartitionsCost.tripleJoinCosts
			if tripleJoinCosts != nil {
				index3DPart2 := part2 - (i + j + 2)
				*tripleJoinCosts = append((*tripleJoinCosts)[:index3DPart2], (*tripleJoinCosts)[index3DPart2+1:]...)
				// check if the deleted cost were the minimum
				if twoPartitionsCost.bestJoin == index3DPart2 {
					min, argmin := utils.MinAndArgMin(*tripleJoinCosts)
					twoPartitionsCost.bestJoin = argmin
					twoPartitionsCost.joinCost = min
				}
			}
		}
		// delete part2 in the second dimension
		algorithm.costs.Delete2D(i, part2)
		lenItsSupposedToBe := 99 - (100 - len(*algorithm.costs))
		lenItActuallyIs := len((*algorithm.costs)[0].twoPartitionsCosts)
		if i == 0 && lenItsSupposedToBe <= lenItActuallyIs {
			_ = 1
		}

		algorithm.costs.Min2D(i, bestJoinOverall, bestJoinCostOverall)
	}
}

// Adjust the join costs in second dimension of partition part1. If the join procedure can be
// ended after this function it returns true.
func (algorithm *GreedyJoiningAlgorithm[data]) joinStep2(part1, part2 int, previousJoinCost float64,
	bestJoinOverall *[2]int, bestJoinCostOverall *float64) {

	twoPartitionsCosts := ((*(algorithm.costs))[part1]).twoPartitionsCosts
	index2DPart2 := algorithm.costs.GetIndex(part1, part2)

	for j := range twoPartitionsCosts {
		// the cost for part1 and part2 doesn't have to be updated as it will be deleted later on
		if j == index2DPart2 {
			continue
		}
		jPartition := (part1 + j + 1)

		// determine if the new cost have been calculated by looking into triple costs
		var costAlreadyCalculated bool
		if j < index2DPart2 {
			costAlreadyCalculated = twoPartitionsCosts[j].tripleJoinCosts != nil
		} else {
			costAlreadyCalculated = twoPartitionsCosts[algorithm.costs.GetIndex(part1, part2)].tripleJoinCosts != nil
		}

		if costAlreadyCalculated {
			// the new join cost has already been calculated
			tripleCost, err := algorithm.costs.TripleJoinCost(part1, jPartition, part2)
			if err != nil {
				panic(fmt.Errorf("The algorithm tried to access triple costs for partitions %d, %d and %d, but these costs weren't present!", part1, part2, jPartition))
			}
			twoPartitionsCosts[j].joinCost = tripleCost - previousJoinCost
			twoPartitionsCosts[j].tripleJoinCosts = nil
		} else {
			// the new cost require some additional calculations: the cost of triples where each element is in
			// j, part1 and part2 respectively
			additionalCost := algorithm.TripleCost3Part(part1, jPartition, part2)
			twoPartitionsCosts[j].joinCost = algorithm.costs.RealJoinCost(part1, jPartition) + algorithm.costs.RealJoinCost(jPartition, part2) + additionalCost
		}
	}
	// Delete the costs of part1 with part2 in the second dimension
	algorithm.costs.Delete2D(part1, part2)

	// recompute the minimum over all joins of part1
	algorithm.costs.Min2D(part1, bestJoinOverall, bestJoinCostOverall)
}

func (algorithm *GreedyJoiningAlgorithm[data]) joinStep3(part1, part2 int, previousJoinCost float64,
	bestJoinOverall *[2]int, bestJoinCostOverall *float64) {
	// Delete part2 in second and third dimension of all elements in first dimension between part1 and part2
	for i := part1 + 1; i < part2; i++ {
		twoPartitionsCosts := (*algorithm.costs)[i].twoPartitionsCosts
		index2DPart2 := part2 - i - 1

		// Delete part2 triple join costs in third dimension
		for j := 0; j < index2DPart2; j++ {
			if tripleCosts := twoPartitionsCosts[j].tripleJoinCosts; tripleCosts != nil {
				algorithm.costs.Delete3D(i, (i + j + 1), part2)
			}
		}

		// Delete part2 join cost in second dimension
		algorithm.costs.Delete2D(i, part2)

		// recompute the minimum over all joins of partition i
		algorithm.costs.Min2D(i, bestJoinOverall, bestJoinCostOverall)
	}
}

// Delete costs for part2 in first dimension
func (algorithm *GreedyJoiningAlgorithm[data]) joinStep4(part1, part2 int, previousJoinCost float64,
	bestJoinOverall *[2]int, bestJoinCostOverall *float64) {
	// check if part2 is the last partition which is not stored in the data structure
	if part2 == len(*algorithm.costs) {
		*algorithm.costs = (*algorithm.costs)[:len(*algorithm.costs)-1]
	} else {
		*algorithm.costs = append((*algorithm.costs)[:part2], (*algorithm.costs)[part2+1:]...)
	}

	// Check if the joins in the partitions after part2 are better than the previously
	// observed join costs
	for i := part2; i < len(*algorithm.costs); i++ {
		if onePartitionCosts := (*algorithm.costs)[i]; onePartitionCosts.minCost < *bestJoinCostOverall {
			*bestJoinCostOverall = onePartitionCosts.minCost
			(*bestJoinOverall)[0] = i
			(*bestJoinOverall)[1] = i + onePartitionCosts.bestJoin + 1
		}
	}
}

// The greedy joining algorithm with following properties:
// 	- it will only evaluate joins of 3 partitions if the first 2 contain one element
// 	- it only joins 2 partitions
// 	- if there is only one partition left the algorithm terminates
func GeedyJoining[data any](input *[]data, calc CostCalculator[data]) PartitioningArray {
	// algorithm := Algorithm[data]{input: input, calc: calc}
	algorithm := GetAlgorithm(input, calc)
	nextJoin, cost := algorithm.InitializeAlgorithm()

	for cost < 0 && nextJoin[0] != -1 && nextJoin[1] != -1 {
		nextJoin, cost = algorithm.Join(nextJoin[0], nextJoin[1])
	}
	return algorithm.partitioning
}

// This is GreedyJoiningV2 from tag v1.0 on master
func GeedyJoiningOld[data any](input []data, calc CostCalculator[data]) Partitioning[data] {
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

	joinAdded := false
	newPartitioning := [][]data{}
	for i, partition := range partitioning {
		if !utils.Contains(partitions, i) {
			newPartitioning = append(newPartitioning, partition)
		} else if !joinAdded {
			newPartitioning = append(newPartitioning, join)
			joinAdded = true
		}
	}

	return newPartitioning
}
