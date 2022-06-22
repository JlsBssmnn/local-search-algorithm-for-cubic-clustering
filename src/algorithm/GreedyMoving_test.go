package algorithm

import (
	"testing"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"
	"github.com/stretchr/testify/assert"
)

var movingAlgorithm GreedyMovingAlgorithm[string]

func init() {
	dataPoints := []string{"b", "c", "hello", "but", "howdy", "charly", "big", "delta", "brother", "humor"}

	movingAlgorithm = GreedyMovingAlgorithm[string]{input: &dataPoints, calc: &CharCostCalc{}}
	movingAlgorithm.Initialize()
}

func TestPartitioningDataStruct(t *testing.T) {
	for i := 0; i < len(*movingAlgorithm.input); i++ {
		assert.Equal(t, i, movingAlgorithm.partitioning[i])
		assert.Equal(t, 1, len(*movingAlgorithm.partitions[i]))
		assert.Equal(t, i, (*movingAlgorithm.partitions[i])[0])
	}

	movingAlgorithm.updatePartitioning(-1, 5, 1)

	assert.Equal(t, 1, movingAlgorithm.partitioning[1])
	assert.Equal(t, 1, movingAlgorithm.partitioning[5])
	assert.Equal(t, 2, len(*movingAlgorithm.partitions[1]))
	assert.Equal(t, 2, len(*movingAlgorithm.partitions[5]))
	assert.Equal(t, 1, (*movingAlgorithm.partitions[1])[0])
	assert.Equal(t, 5, (*movingAlgorithm.partitions[1])[1])

	movingAlgorithm.updatePartitioning(-1, 1, 2)

	assert.Equal(t, 1, movingAlgorithm.partitioning[1])
	assert.Equal(t, 1, movingAlgorithm.partitioning[2])
	assert.Equal(t, 3, len(*movingAlgorithm.partitions[1]))
	assert.Equal(t, 3, len(*movingAlgorithm.partitions[2]))
	assert.Equal(t, 1, (*movingAlgorithm.partitions[1])[0])
	assert.Equal(t, 2, (*movingAlgorithm.partitions[1])[1])
	assert.Equal(t, 5, (*movingAlgorithm.partitions[1])[2])

	movingAlgorithm.updatePartitioning(2, 0, 1)

	assert.Equal(t, 0, movingAlgorithm.partitioning[1])
	assert.Equal(t, 2, movingAlgorithm.partitioning[2])
	assert.Equal(t, 2, movingAlgorithm.partitioning[5])
	assert.Equal(t, 2, len(*movingAlgorithm.partitions[1]))
	assert.Equal(t, 2, len(*movingAlgorithm.partitions[2]))
	assert.Equal(t, 0, (*movingAlgorithm.partitions[1])[0])
	assert.Equal(t, 1, (*movingAlgorithm.partitions[1])[1])
	assert.Equal(t, 2, (*movingAlgorithm.partitions[2])[0])
	assert.Equal(t, 5, (*movingAlgorithm.partitions[2])[1])

	movingAlgorithm.updatePartitioning(1, 2, 0)

	assert.Equal(t, 1, movingAlgorithm.partitioning[1])
	assert.Equal(t, 0, movingAlgorithm.partitioning[0])
	assert.Equal(t, 0, movingAlgorithm.partitioning[2])
	assert.Equal(t, 1, len(*movingAlgorithm.partitions[1]))
	assert.Equal(t, 3, len(*movingAlgorithm.partitions[0]))
	assert.Equal(t, 3, len(*movingAlgorithm.partitions[2]))
	assert.Equal(t, 1, (*movingAlgorithm.partitions[1])[0])
	assert.Equal(t, 0, (*movingAlgorithm.partitions[2])[0])
	assert.Equal(t, 2, (*movingAlgorithm.partitions[2])[1])
	assert.Equal(t, 5, (*movingAlgorithm.partitions[2])[2])
}

func TestInitializeCost(t *testing.T) {
	dataPoints := []string{"b", "c", "hello", "but", "howdy", "charly", "big", "delta", "brother", "humor"}
	n := len(dataPoints)

	algorithm := GreedyMovingAlgorithm[string]{input: &dataPoints, calc: &CharCostCalc{}}
	nextMove, cost := algorithm.Initialize()

	assert.Equal(t, 0, nextMove[1])
	assert.Equal(t, 3, nextMove[0])
	assert.Equal(t, 6, nextMove[2])
	assert.Equal(t, -1.0, cost)

	for i := 0; i < n; i++ {
		assert.Equal(t, n, len((*algorithm.costs)[i].moves))
		for j := 0; j < n; j++ {
			if i == j {
				assert.False(t, (*algorithm.costs)[i].moves[j].valid)
			} else {
				x := 0
				if j > i {
					x++
				}
				doubleMoveLen := n - 1 - i - x

				if doubleMoveLen < 1 {
					assert.False(t, (*algorithm.costs)[i].moves[j].valid)
				} else {
					assert.True(t, (*algorithm.costs)[i].moves[j].valid)
					assert.Equal(t, n-1-i-x, len(*(*algorithm.costs)[i].moves[j].doubleMoves))
					assert.Equal(t, utils.Min(utils.Map(*(*algorithm.costs)[i].moves[j].doubleMoves, func(e DoubleMove) float64 {
						return e.cost
					})), (*algorithm.costs)[i].moves[j].cost)
				}
			}
		}
	}
}

func TestInitialCosts(t *testing.T) {
	getMoveCost := func(i, j int) float64 {
		return (*movingAlgorithm.costs)[i].moves[j].cost
	}
	getTripleMoveCost := func(i, j, k int) float64 {
		return (*(*movingAlgorithm.costs)[i].moves[j].doubleMoves)[k].cost
	}
	t.Run("Test values for element 0", func(t *testing.T) {
		assert.Equal(t, -1.0, (*movingAlgorithm.costs)[0].minCost)
		assert.Equal(t, 1.0, getMoveCost(0, 1))
		assert.Equal(t, 1.0, getTripleMoveCost(0, 1, 0))
		assert.Equal(t, 1.0, getTripleMoveCost(0, 1, 1))
		assert.Equal(t, 1.0, getTripleMoveCost(0, 1, 2))

		assert.Equal(t, 1.0, getMoveCost(0, 2))
		assert.Equal(t, 1.0, getTripleMoveCost(0, 2, 0))
		assert.Equal(t, 1.0, getTripleMoveCost(0, 2, 1))
		assert.Equal(t, 1.0, getTripleMoveCost(0, 2, 2))

		assert.Equal(t, -1.0, getMoveCost(0, 3))
		assert.Equal(t, 1.0, getTripleMoveCost(0, 3, 0))
		assert.Equal(t, 1.0, getTripleMoveCost(0, 3, 1))
		assert.Equal(t, 1.0, getTripleMoveCost(0, 3, 2))
		assert.Equal(t, -1.0, getTripleMoveCost(0, 3, 4))
		assert.Equal(t, -1.0, getTripleMoveCost(0, 3, 6))
	})

	t.Run("Test values for element 1", func(t *testing.T) {
		assert.Equal(t, 1.0, (*movingAlgorithm.costs)[1].minCost)
		assert.Equal(t, 1.0, getMoveCost(1, 0))
		assert.Equal(t, 1.0, getTripleMoveCost(1, 0, 0))
		assert.Equal(t, 1.0, getTripleMoveCost(1, 0, 1))
		assert.Equal(t, 1.0, getTripleMoveCost(1, 0, 2))

		assert.Equal(t, 1.0, getMoveCost(1, 5))
		assert.Equal(t, 1.0, getTripleMoveCost(1, 5, 0))
		assert.Equal(t, 1.0, getTripleMoveCost(1, 5, 1))
		assert.Equal(t, 1.0, getTripleMoveCost(1, 5, 2))
		assert.Equal(t, 1.0, getTripleMoveCost(1, 5, 3))
	})

	t.Run("Test values for element 2", func(t *testing.T) {
		assert.Equal(t, -1.0, (*movingAlgorithm.costs)[2].minCost)
		assert.Equal(t, -1.0, getMoveCost(2, 4))
		assert.Equal(t, 1.0, getTripleMoveCost(2, 4, 0))
		assert.Equal(t, 1.0, getTripleMoveCost(2, 4, 1))
		assert.Equal(t, -1.0, getTripleMoveCost(2, 4, 5))

		assert.Equal(t, 1.0, getMoveCost(2, 5))
	})
}

func TestMove(t *testing.T) {
	dataPoints := []string{"b", "c", "hello", "but", "howdy", "charly", "big", "delta", "brother", "humor"}

	algorithm := GreedyMovingAlgorithm[string]{input: &dataPoints, calc: &CharCostCalc{}}
	nextMove, _ := algorithm.Initialize()

	algorithm.Move(nextMove[0], nextMove[1])

	t.Run("Cost for element 0", func(t *testing.T) {
		assert.True(t, (*algorithm.costs)[0].moves[0].valid)
		assert.Equal(t, 0.0, (*algorithm.costs)[0].moves[0].cost)
		assert.False(t, (*algorithm.costs)[0].moves[3].valid)
	})
	t.Run("Cost for element 1", func(t *testing.T) {
		assert.Equal(t, 1.0, (*algorithm.costs)[1].moves[0].cost)
		assert.Negative(t, (*algorithm.costs)[1].moves[0].bestMove)
		assert.False(t, (*algorithm.costs)[1].moves[3].valid)
	})
	t.Run("Cost for element 3", func(t *testing.T) {
		assert.True(t, (*algorithm.costs)[3].moves[3].valid)
		assert.Equal(t, 0.0, (*algorithm.costs)[3].moves[3].cost)
		assert.False(t, (*algorithm.costs)[3].moves[0].valid)
	})

	algorithm.Move(nextMove[0], nextMove[2])

	t.Run("Cost for element 0", func(t *testing.T) {
		assert.True(t, (*algorithm.costs)[0].moves[0].valid)
		assert.Equal(t, 1.0, (*algorithm.costs)[0].moves[0].cost)
		assert.False(t, (*algorithm.costs)[0].moves[6].valid)
	})
	t.Run("Cost for element 1", func(t *testing.T) {
		assert.Equal(t, 3.0, (*algorithm.costs)[1].moves[0].cost)
		assert.Negative(t, (*algorithm.costs)[1].moves[0].bestMove)
		assert.False(t, (*algorithm.costs)[1].moves[6].valid)
	})
	t.Run("Cost for element 8", func(t *testing.T) {
		assert.True(t, (*algorithm.costs)[8].moves[0].valid)
		assert.Equal(t, -3.0, (*algorithm.costs)[8].moves[0].cost)
		assert.False(t, (*algorithm.costs)[8].moves[3].valid)
		assert.False(t, (*algorithm.costs)[8].moves[6].valid)
	})
}
