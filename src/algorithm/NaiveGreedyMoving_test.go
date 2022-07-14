package algorithm

import (
	"testing"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"
	"github.com/stretchr/testify/assert"
)

func TestDataStructure(t *testing.T) {
	dataPoints := []string{"b", "c", "hello", "but", "howdy", "charley", "big", "delta", "brother", "humor"}

	algorithm := NaiveGreedyMovingAlgorithm[string]{input: &dataPoints, calc: &CharCostCalc{}}
	algorithm.initialize()

	assert.Equal(t, 10, len(algorithm.partitioning))
	assert.Equal(t, 10, len(algorithm.partitions))
	for _, partition := range algorithm.partitions {
		assert.Equal(t, 1, len(*partition))
	}

	algorithm.moveElement(1, 0)

	assert.Equal(t, 10, len(algorithm.partitioning))
	assert.True(t, algorithm.partitioning[0] == algorithm.partitioning[1])
	assert.Equal(t, 9, len(algorithm.partitions))
	assert.Equal(t, 2, len(*algorithm.partitions[1]))
	_, ok := algorithm.partitions[0]
	assert.False(t, ok)

	algorithm.moveElement(1, 2)

	assert.True(t, algorithm.partitioning[0] == algorithm.partitioning[2])
	assert.Equal(t, 8, len(algorithm.partitions))
	assert.Equal(t, 3, len(*algorithm.partitions[1]))
	_, ok = algorithm.partitions[2]
	assert.False(t, ok)

	algorithm.moveElement(10, 1)

	assert.False(t, algorithm.partitioning[0] == algorithm.partitioning[1])
	assert.False(t, algorithm.partitioning[2] == algorithm.partitioning[1])
	assert.Equal(t, 9, len(algorithm.partitions))
	assert.Equal(t, 2, len(*algorithm.partitions[1]))
	assert.Equal(t, 1, len(*algorithm.partitions[10]))
}

func TestFindBestMove(t *testing.T) {
	dataPoints := []string{"b", "c", "hello", "but", "howdy", "charley", "big", "delta", "brother", "humor"}

	algorithm := NaiveGreedyMovingAlgorithm[string]{input: &dataPoints, calc: &CharCostCalc{}}
	algorithm.initialize()

	nextMove, cost := algorithm.findBestMove()
	U, a, b := nextMove[0], nextMove[1], nextMove[2]

	assert.Equal(t, -1.0, cost)
	assert.Equal(t, 0, a)
	assert.True(t, utils.Contains([]int{3, 6, 8}, U))
	assert.True(t, utils.Contains([]int{3, 6}, b))

	algorithm.moveElement(0, 3)
	algorithm.moveElement(0, 6)

	nextMove, cost = algorithm.findBestMove()
	U, a, b = nextMove[0], nextMove[1], nextMove[2]

	assert.Equal(t, -3.0, cost)
	assert.Equal(t, 0, U)
	assert.Equal(t, a, 8)
	assert.Equal(t, -1, b)

	algorithm.moveElement(0, 8)

	nextMove, cost = algorithm.findBestMove()
	U, a, b = nextMove[0], nextMove[1], nextMove[2]

	assert.Equal(t, -1.0, cost)
	assert.Equal(t, a, 2)
	assert.True(t, utils.Contains([]int{4, 9}, U))
	assert.True(t, utils.Contains([]int{4, 9}, b))

	algorithm.moveElement(2, 4)
	algorithm.moveElement(2, 9)

	nextMove, cost = algorithm.findBestMove()
	U, a, b = nextMove[0], nextMove[1], nextMove[2]

	assert.Equal(t, 1.0, cost)
	assert.Equal(t, a, 1)
	assert.True(t, utils.Contains([]int{5, 7}, U))
	assert.True(t, utils.Contains([]int{5, 7}, b))

	algorithm.moveElement(1, 5)
	algorithm.moveElement(1, 7)

	nextMove, cost = algorithm.findBestMove()
	U, a, b = nextMove[0], nextMove[1], nextMove[2]

	assert.Equal(t, -1.0, cost)
	assert.Equal(t, a, 1)
	assert.Equal(t, 3, U)
	assert.Equal(t, -1, b)
}
