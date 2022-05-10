package algorithm

import (
	"testing"

	_ "github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
	"github.com/stretchr/testify/assert"
)

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

func TestCubicPartitionCost(t *testing.T) {
	calc := CharCostCalc{}
	cost := CubicPartitionCost[string](Partitioning[string]{{"hey", "hello", "hohoho", "howdy"}, {"a", "b"}}, calc)
	assert.Equal(t, -4.0, cost, "First partition gives -3 cost, second gives 0 because there are < 3 elements")

	cost = CubicPartitionCost[string](Partitioning[string]{{"hey", "hello", "howdy"}, {"a", "b", "c"}}, calc)
	assert.Equal(t, 0.0, cost, "First partition gives -1 cost, second gives +1")

	cost = CubicPartitionCost[string](Partitioning[string]{{"hey", "hello", "howdy"}, {"a", "b", "c", "b"}}, calc)
	assert.Equal(t, 3.0, cost, "First partition gives -1 cost, second gives +4")

	cost = CubicPartitionCost[string](Partitioning[string]{{"hey", "hello"}, {"a", "b", "c", "b", "b"}}, calc)
	assert.Equal(t, 8.0, cost, "First partition gives 0 cost, second gives 8")
}
