package algorithm

import (
	"testing"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
	"github.com/stretchr/testify/assert"
)

var delta = 0.00000001

var vectors = []struct {
	v1, v2, v3 geometry.Vector
	maxDist    float64
}{
	{geometry.Vector{X: 0, Y: 0, Z: 0}, geometry.Vector{X: 1, Y: 2, Z: 3}, geometry.Vector{X: 7, Y: 0, Z: 0}, 7},
	{geometry.Vector{X: -3, Y: 2, Z: -1}, geometry.Vector{X: 1, Y: -1, Z: 8}, geometry.Vector{X: 3, Y: 4, Z: 0}, 10.295630140987},
	{geometry.Vector{X: 0, Y: 0, Z: 0}, geometry.Vector{X: 0, Y: 0, Z: 0}, geometry.Vector{X: 0, Y: 0, Z: 0}, 0},
	{geometry.Vector{X: 0, Y: 1, Z: 0}, geometry.Vector{X: 1, Y: 0, Z: 0}, geometry.Vector{X: 0, Y: 0, Z: 1}, 1.4142135623730951},
	{geometry.Vector{X: -100, Y: 0, Z: 2}, geometry.Vector{X: -2, Y: 2, Z: 6}, geometry.Vector{X: 2, Y: 1, Z: 2}, 102.00490184299969},
	{geometry.Vector{X: 0, Y: 0, Z: 0}, geometry.Vector{X: 1, Y: 1, Z: 1}, geometry.Vector{X: 2, Y: 2, Z: 2}, 3.4641016151377544},
}

func TestGetMaxDist(t *testing.T) {
	for _, element := range vectors {
		assert.InDelta(t, element.maxDist, GetMaxDist(element.v1, element.v2, element.v3), delta, "MaxDist should be calculated correctly")
	}
}

func TestTripleCost1Part(t *testing.T) {
	calc := CostCalculator{treshold: 2, amplification: 4}
	assert.InDelta(t, 33.182520563948,
		calc.TripleCost1Part(geometry.Vector{X: -3, Y: 2, Z: -1}, geometry.Vector{X: 1, Y: -1, Z: 8}, geometry.Vector{X: 3, Y: 4, Z: 0}),
		delta,
		"The MaxDist is above the threshold",
	)
	assert.InDelta(t, 400.01960737199875,
		calc.TripleCost1Part(geometry.Vector{X: -100, Y: 0, Z: 2}, geometry.Vector{X: -2, Y: 2, Z: 6}, geometry.Vector{X: 2, Y: 1, Z: 2}),
		delta,
		"The MaxDist is above the threshold",
	)

	assert.Zero(t,
		calc.TripleCost1Part(geometry.Vector{X: 0, Y: 0, Z: 0}, geometry.Vector{X: 0, Y: 0, Z: 0}, geometry.Vector{X: 0, Y: 0, Z: 0}),
		"The MaxDist is below the threshold",
	)
	assert.Zero(t,
		calc.TripleCost1Part(geometry.Vector{X: 0, Y: 1, Z: 0}, geometry.Vector{X: 1, Y: 0, Z: 0}, geometry.Vector{X: 0, Y: 0, Z: 1}),
		"The MaxDist is below the threshold",
	)
}

func TestTripleCost3Part(t *testing.T) {
	calc := CostCalculator{treshold: 5.2, amplification: 1.5}

	assert.Zero(t,
		calc.TripleCost3Part(geometry.Vector{X: -3, Y: 2, Z: -1}, geometry.Vector{X: 1, Y: -1, Z: 8}, geometry.Vector{X: 3, Y: 4, Z: 0}),
		"The MaxDist is above the threshold",
	)
	assert.Zero(t,
		calc.TripleCost3Part(geometry.Vector{X: -100, Y: 0, Z: 2}, geometry.Vector{X: -2, Y: 2, Z: 6}, geometry.Vector{X: 2, Y: 1, Z: 2}),
		"The MaxDist is above the threshold",
	)

	assert.InDelta(t,
		7.8,
		calc.TripleCost3Part(geometry.Vector{X: 0, Y: 0, Z: 0}, geometry.Vector{X: 0, Y: 0, Z: 0}, geometry.Vector{X: 0, Y: 0, Z: 0}),
		delta,
		"The MaxDist is below the threshold",
	)
	assert.InDelta(t,
		5.678679656440358,
		calc.TripleCost3Part(geometry.Vector{X: 0, Y: 1, Z: 0}, geometry.Vector{X: 1, Y: 0, Z: 0}, geometry.Vector{X: 0, Y: 0, Z: 1}),
		delta,
		"The MaxDist is below the threshold",
	)
}

func TestCubicPartitionCost(t *testing.T) {
	calc := CostCalculator{treshold: 1, amplification: 1}

	v1 := geometry.Vector{X: 1, Y: 1, Z: 1}
	v2 := geometry.Vector{X: 0.9, Y: 0.9, Z: 0.9}
	v3 := geometry.Vector{X: 1.1, Y: 1.1, Z: 1.1}

	assert.Zero(t, calc.CubicPartitionCost([][]geometry.Vector{{v1, v2, v3}}), "MaxDist is below the treshold")
	assert.InDelta(t, calc.CubicPartitionCost([][]geometry.Vector{{v1}, {v2}, {v3}}),
		1-0.34641016151377,
		delta,
		"MaxDist is below the treshold so there should be costs")
	assert.Zero(t, calc.CubicPartitionCost([][]geometry.Vector{{v1, v2}, {v3}}), "There should be no costs")

	calc = CostCalculator{treshold: 5, amplification: 3}

	v1 = geometry.Vector{X: -10, Y: 0, Z: 1}
	v2 = geometry.Vector{X: 0, Y: 0, Z: 0}
	v3 = geometry.Vector{X: 20, Y: 1, Z: -5}

	assert.InDelta(t,
		calc.CubicPartitionCost([][]geometry.Vector{{v1, v2, v3}}),
		(30.610455730027933-5)*3,
		delta,
		"MaxDist is above the treshold so there should be costs")
	assert.Zero(t, calc.CubicPartitionCost([][]geometry.Vector{{v1}, {v2}, {v3}}),
		"MaxDist is above the treshold so there should not be costs")
}
