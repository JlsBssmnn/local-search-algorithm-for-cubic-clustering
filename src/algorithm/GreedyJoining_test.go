package algorithm

import (
	"math"
	"testing"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"
	"github.com/stretchr/testify/assert"
)

var algorithm GreedyJoiningAlgorithm[string]

func init() {
	dataPoints := []string{"b", "c", "hello", "but", "howdy", "charly", "big", "delta", "brother", "humor"}

	costs, _, _ := InitializeCosts[string](&dataPoints, CharCostCalc{})
	algorithm = GreedyJoiningAlgorithm[string]{costs: &costs, input: &dataPoints, calc: &CharCostCalc{}, partitioning: []int{}}
}

type CharCostCalc struct{}

// if two strings start with the same char -> cost -1
// otherwise cost 1
func (calc CharCostCalc) TripleCost(c1, c2, c3 string) float64 {
	if c1[0:1] == c2[0:1] && c2[0:1] == c3[0:1] {
		return -1
	} else {
		return 1
	}
}

func TestInitializeCosts(t *testing.T) {
	dataPoints := []string{"b", "c", "hello", "but", "howdy", "charly", "big", "delta", "brother", "humor"}

	costs, bestJoin, bestJoinCost := InitializeCosts[string](&dataPoints, CharCostCalc{})
	t.Run("Test length of costs in first dimension", func(t *testing.T) {
		assert.Equal(t, len(dataPoints)-1, len(costs), "At the start every partition has costs with others except the last one")
		assert.Equal(t, len(dataPoints)-1, len((*costs[0]).twoPartitionsCosts), "The first element has costs with all the others")
		assert.Equal(t, len(dataPoints)-3, len((*costs[2]).twoPartitionsCosts), "The third element has costs with 7 others")
		assert.Equal(t, len(dataPoints)-9, len((*costs[8]).twoPartitionsCosts), "The 9th element has costs with only the last element")
	})

	t.Run("Test lengths of costs in third dimension", func(t *testing.T) {
		getTripleJoinCosts := func(i, j int) []float64 {
			return *((*costs[i]).twoPartitionsCosts[j]).tripleJoinCosts
		}
		assert.Equal(t, len(dataPoints)-2, len(getTripleJoinCosts(0, 0)), "The list of triple join costs for 1 and 2")
		assert.Equal(t, len(dataPoints)-4, len(getTripleJoinCosts(0, 2)), "The list of triple join costs for 1 and 4")
		assert.Equal(t, 0, len(getTripleJoinCosts(0, 8)), "Join of 1 and last element has no triple join costs")
		assert.Equal(t, len(dataPoints)-7, len(getTripleJoinCosts(4, 1)), "The list of triple join costs for 5 and 7")
	})

	t.Run("Test for correct costs in third dimension", func(t *testing.T) {
		getTripleJoinCosts := func(i, j int) []float64 {
			return *((*costs[i]).twoPartitionsCosts[j]).tripleJoinCosts
		}
		assert.Equal(t, []float64{1, 1, 1, 1, 1, 1, 1, 1}, getTripleJoinCosts(0, 0), "Triple join costs for 1 and 2")
		assert.Equal(t, []float64{1, 1, -1, 1, -1, 1}, getTripleJoinCosts(0, 2), "Triple join costs for 1 and 4")
		assert.Equal(t, []float64{1, 1, 1, 1, 1, 1}, getTripleJoinCosts(2, 0), "Triple join costs for 3 and 4")
		assert.Equal(t, []float64{1, 1, 1, 1, -1}, getTripleJoinCosts(2, 1), "Triple join costs for 3 and 5")
	})

	t.Run("Test that the join cost of 2 partitions is the minimum of the triple join costs", func(t *testing.T) {
		getTripleJoinCosts := func(i, j int) []float64 {
			return *((*costs[i]).twoPartitionsCosts[j]).tripleJoinCosts
		}
		getJoinCost := func(i, j int) float64 {
			return ((*costs[i]).twoPartitionsCosts[j]).joinCost
		}
		assert.Equal(t, utils.Min(getTripleJoinCosts(0, 0)), getJoinCost(0, 0))
		assert.Equal(t, utils.Min(getTripleJoinCosts(0, 3)), getJoinCost(0, 3))
		assert.Equal(t, utils.Min(getTripleJoinCosts(2, 4)), getJoinCost(2, 4))
		assert.Equal(t, utils.Min(getTripleJoinCosts(3, 1)), getJoinCost(3, 1))
		assert.Equal(t, 0, len(getTripleJoinCosts(6, 2)))
	})

	assert.LessOrEqual(t, 0, bestJoin[0], "Some partition should be in the best join")
	assert.LessOrEqual(t, 0, bestJoin[1], "Some partition should be in the best join")
	assert.Equal(t, -1.0, bestJoinCost, "Best cost is -1")
	assert.Equal(t, [2]int{0, 3}, bestJoin, "The best join is the first pair of elements that start with b")
}

func TestJoinPart1is0(t *testing.T) {
	dataPoints := []string{"b", "c", "hello", "but", "howdy", "charly", "big", "delta", "brother", "humor"}

	costs, _, _ := InitializeCosts[string](&dataPoints, CharCostCalc{})
	algorithm := GreedyJoiningAlgorithm[string]{costs: &costs, input: &dataPoints, calc: &CharCostCalc{}, partitioning: []int{}}
	bestJoin, bestCost := algorithm.Join(0, 3)

	assert.Equal(t, 0, bestJoin[0], "Next join with cost -1 is between the partitions 0 and 5")
	assert.Equal(t, 5, bestJoin[1], "Next join with cost -1 is between the partitions 0 and 5")
	assert.Equal(t, -1.0, bestCost, "Best cost should be -1")

	t.Run("Test length of costs in first dimension", func(t *testing.T) {
		assert.Equal(t, len(dataPoints)-2, len(costs), "At the start every partition has costs with others except the last one")
		assert.Equal(t, len(dataPoints)-2, len((*costs[0]).twoPartitionsCosts), "The first element has costs with all the others")
		assert.Equal(t, len(dataPoints)-4, len((*costs[2]).twoPartitionsCosts), "The third element has costs with 7 others")
		assert.Equal(t, len(dataPoints)-6, len((*costs[4]).twoPartitionsCosts), "The 5th element has costs with 4 others")
	})

	t.Run("Test lengths of costs in third dimension", func(t *testing.T) {
		getTripleJoinCosts := func(i, j int) *[]float64 {
			return ((*costs[i]).twoPartitionsCosts[j]).tripleJoinCosts
		}
		assert.Nil(t, getTripleJoinCosts(0, 0), "0 has 2 elements, so no tripple join costs")
		assert.Nil(t, getTripleJoinCosts(0, 2), "0 has 2 elements, so no tripple join costs")
		assert.Nil(t, getTripleJoinCosts(0, 7), "0 has 2 elements, so no tripple join costs")
		assert.Equal(t, len(dataPoints)-4, len(*getTripleJoinCosts(1, 0)), "The list of triple join costs for 2 and 3")
		assert.Equal(t, len(dataPoints)-6, len(*getTripleJoinCosts(1, 2)), "The list of triple join costs for 2 and 5")
		assert.Equal(t, len(dataPoints)-6, len(*getTripleJoinCosts(3, 0)), "Element 4 should be same as 5 in previous step")
		assert.Equal(t, len(dataPoints)-9, len(*getTripleJoinCosts(3, 3)), "Element 4 should be same as 5 in previous step")
	})

	t.Run("Test for correct costs in third dimension", func(t *testing.T) {
		getTripleJoinCosts := func(i, j int) []float64 {
			return *((*costs[i]).twoPartitionsCosts[j]).tripleJoinCosts
		}
		assert.Equal(t, []float64{1, 1, 1, 1, -1}, getTripleJoinCosts(2, 0), "Triple join costs for 3 and 4")
		assert.Equal(t, []float64{1, 1, 1, 1}, getTripleJoinCosts(2, 1), "Triple join costs for 3 and 5")
		assert.Equal(t, []float64{1, 1, 1}, getTripleJoinCosts(2, 2), "Triple join costs for 3 and 5")
		assert.Equal(t, []float64{1, 1, 1}, getTripleJoinCosts(4, 0), "Triple join costs for 5 and 6")
		assert.Equal(t, []float64{1, 1}, getTripleJoinCosts(4, 1), "Triple join costs for 5 and 7")
		assert.Equal(t, []float64{1}, getTripleJoinCosts(4, 2), "Triple join costs for 5 and 8")
	})

	t.Run("Test join cost in second dimension", func(t *testing.T) {
		getCostOfSecondDim := func(i int) []float64 {
			costs2D := make([]float64, len(dataPoints)-i-2)
			for j, join := range (*costs[i]).twoPartitionsCosts {
				costs2D[j] = join.joinCost
			}
			return costs2D
		}
		assert.Equal(t, []float64{1, 1, 1, 1, -1, 1, -1, 1}, getCostOfSecondDim(0))
		assert.Equal(t, []float64{1, 1, 1, 1, 1, 1, math.Inf(1)}, getCostOfSecondDim(1))
		assert.Equal(t, []float64{1, 1, 1, 1, math.Inf(1)}, getCostOfSecondDim(3))
		assert.Equal(t, utils.Min(getCostOfSecondDim(0)), (*costs[0]).minCost)
		assert.Equal(t, utils.Min(getCostOfSecondDim(1)), (*costs[1]).minCost)
		assert.Equal(t, utils.Min(getCostOfSecondDim(3)), (*costs[3]).minCost)
	})
}

func TestPartsInMiddle(t *testing.T) {
	dataPoints := []string{"b", "c", "hello", "but", "howdy", "charly", "big", "delta", "brother", "humor"}

	costs, _, _ := InitializeCosts[string](&dataPoints, CharCostCalc{})
	algorithm := GreedyJoiningAlgorithm[string]{costs: &costs, input: &dataPoints, calc: &CharCostCalc{}, partitioning: []int{}}
	bestJoin, bestCost := algorithm.Join(5, 8)

	assert.Equal(t, 0, bestJoin[0], "Next join with cost -1 is between the partitions 0 and 3")
	assert.Equal(t, 3, bestJoin[1], "Next join with cost -1 is between the partitions 0 and 3")
	assert.Equal(t, -1.0, bestCost, "Best cost should be -1")

	t.Run("Test length of costs in first dimension", func(t *testing.T) {
		assert.Equal(t, len(dataPoints)-2, len(costs), "At the start every partition has costs with others except the last one")
		assert.Equal(t, len(dataPoints)-2, len((*costs[0]).twoPartitionsCosts), "The first element has costs with all the others")
		assert.Equal(t, len(dataPoints)-4, len((*costs[2]).twoPartitionsCosts), "The third element has costs with 7 others")
		assert.Equal(t, len(dataPoints)-6, len((*costs[4]).twoPartitionsCosts), "The 5th element has costs with 4 others")
	})

	t.Run("Test lengths of costs in third dimension", func(t *testing.T) {
		getTripleJoinCosts := func(i, j int) *[]float64 {
			return ((*costs[i]).twoPartitionsCosts[j]).tripleJoinCosts
		}
		assert.Equal(t, 7, len(*getTripleJoinCosts(0, 0)))
		assert.Equal(t, 5, len(*getTripleJoinCosts(0, 2)))
		assert.Equal(t, 0, len(*getTripleJoinCosts(0, 7)))
		assert.Nil(t, getTripleJoinCosts(5, 0))
		assert.Nil(t, getTripleJoinCosts(5, 1))
		assert.Nil(t, getTripleJoinCosts(0, 4))
		assert.Nil(t, getTripleJoinCosts(1, 3))
		assert.Nil(t, getTripleJoinCosts(2, 2))
		assert.Equal(t, 0, len(*getTripleJoinCosts(7, 0)))
	})

	t.Run("Test for correct costs in third dimension", func(t *testing.T) {
		getTripleJoinCosts := func(i, j int) []float64 {
			return *((*costs[i]).twoPartitionsCosts[j]).tripleJoinCosts
		}
		assert.Equal(t, []float64{1, 1, 1, 4, 1, 1, 1}, getTripleJoinCosts(0, 0))
		assert.Equal(t, []float64{1, 2, -1, 1, 1}, getTripleJoinCosts(0, 2))
		assert.Equal(t, []float64{1, 1, 4, 1, 1, 1}, getTripleJoinCosts(1, 0))
		assert.Equal(t, []float64{4, 1, 1, -1}, getTripleJoinCosts(2, 1))
	})

	t.Run("Test join cost in second dimension", func(t *testing.T) {
		getCostOfSecondDim := func(i int) []float64 {
			costs2D := make([]float64, len(dataPoints)-i-2)
			for j, join := range (*costs[i]).twoPartitionsCosts {
				costs2D[j] = join.joinCost
			}
			return costs2D
		}
		assert.Equal(t, []float64{1, 1, -1, 1, 1, 1, 1, math.Inf(1)}, getCostOfSecondDim(0))
		assert.Equal(t, []float64{1, 1, 1, 1, 1, 1, math.Inf(1)}, getCostOfSecondDim(1))
		assert.Equal(t, []float64{1, 1, 1}, getCostOfSecondDim(5))
		assert.Equal(t, utils.Min(getCostOfSecondDim(0)), (*costs[0]).minCost)
		assert.Equal(t, utils.Min(getCostOfSecondDim(1)), (*costs[1]).minCost)
		assert.Equal(t, utils.Min(getCostOfSecondDim(3)), (*costs[3]).minCost)
	})
}

func TestLengthSecondDim(t *testing.T) {
	nextJoin, cost := algorithm.InitializeAlgorithm()
	iteration := 0

	for cost < 0 && nextJoin[0] != -1 && nextJoin[1] != -1 {
		for i, onePartitionCosts := range *algorithm.costs {
			assert.Equal(t, 9-iteration-i, len(onePartitionCosts.twoPartitionsCosts))
		}
		nextJoin, cost = algorithm.Join(nextJoin[0], nextJoin[1])
		iteration++
	}
}
