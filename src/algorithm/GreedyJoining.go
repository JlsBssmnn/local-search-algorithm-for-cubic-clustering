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

	if join.tripleJoinCosts != nil && len(*join.tripleJoinCosts) == 0 {
		return math.Inf(1), true
	}
	return join.joinCost, join.tripleJoinCosts != nil
}

// Extracts the real cost for joining i and j. If they are both one-elementary, 0 is returned.
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

// Deletes the costs for partition k in the third dimension for
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
	} else if twoPartitionsCosts.bestJoin > index3D {
		twoPartitionsCosts.bestJoin = twoPartitionsCosts.bestJoin - 1
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
	algorithm.partitioning.InitializeSingletonSets(len(*algorithm.input))
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
	// Cache that stores join(i, part1) + join(i, part2) + TripleCost3Part(i, part1, part2) with i being the key
	cache := make(map[int]float64)
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

				// Access cache
				costI, ok := cache[i]
				if !ok {
					costI = algorithm.costs.RealJoinCost(i, part1) + algorithm.costs.RealJoinCost(i, part2) + algorithm.TripleCost3Part(i, part1, part2)
					cache[i] = costI
				}
				costJ, ok := cache[secondElement]
				if !ok {
					costJ = algorithm.costs.RealJoinCost(secondElement, part1) + algorithm.costs.RealJoinCost(secondElement, part2) + algorithm.TripleCost3Part(secondElement, part1, part2)
					cache[secondElement] = costJ
				}

				// Triple costs that are already calculated
				tripleCostPart1, err1 := algorithm.costs.TripleJoinCost(i, secondElement, part1)
				tripleCostPart2, err2 := algorithm.costs.TripleJoinCost(i, secondElement, part2)

				if err1 != nil || err2 != nil {
					panic("Tried to access triple costs that don't exist!")
				}

				// update the cost for join with part1
				newTripleCost := costI + costJ + tripleCostPart1 + tripleCostPart2
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
			tripleJoinCosts := (*algorithm.costs)[i].twoPartitionsCosts[j].tripleJoinCosts
			if tripleJoinCosts != nil {
				algorithm.costs.Delete3D(i, j+i+1, part2)
			}
		}
		// delete part2 in the second dimension
		algorithm.costs.Delete2D(i, part2)

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
			twoPartitionsCosts[j].tripleJoinCosts = nil
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
func GreedyJoining[data any](input *[]data, calc CostCalculator[data]) PartitioningArray {
	// algorithm := Algorithm[data]{input: input, calc: calc}
	algorithm := GetAlgorithm(input, calc)
	nextJoin, cost := algorithm.InitializeAlgorithm()

	for cost < 0 && nextJoin[0] != -1 && nextJoin[1] != -1 {
		nextJoin, cost = algorithm.Join(nextJoin[0], nextJoin[1])
	}
	return algorithm.partitioning
}
