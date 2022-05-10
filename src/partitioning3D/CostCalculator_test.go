package partitioning3D

import (
	"testing"

	g "github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
	"github.com/stretchr/testify/assert"
)

const delta = 0.00000001

var data = []struct {
	plane, v1, v2, v3 g.Vector
	maxDist           float64
}{
	{g.Vector{X: 0, Y: 0, Z: 1}, g.Vector{X: 0, Y: 0, Z: 0}, g.Vector{X: 1, Y: 2, Z: 3}, g.Vector{X: 7, Y: 0, Z: 0}, 3},
	{g.Vector{X: 0, Y: 0, Z: 1}, g.Vector{X: 0, Y: 0, Z: 0}, g.Vector{X: -7, Y: 50, Z: -12.5}, g.Vector{X: 3, Y: 9, Z: 10}, 12.5},
	{g.Vector{X: 0, Y: 0, Z: 1}, g.Vector{X: -100, Y: -9, Z: 1}, g.Vector{X: 4, Y: 2, Z: 9}, g.Vector{X: 7, Y: 3, Z: -8}, 9},

	{g.Vector{X: 1, Y: 1, Z: -2}, g.Vector{X: 1, Y: 2, Z: 3}, g.Vector{X: 0, Y: 0, Z: 0}, g.Vector{X: 9, Y: 9, Z: 9}, 1.2247448713915892},
	{g.Vector{X: 1, Y: 1, Z: -2}, g.Vector{X: 0, Y: 0, Z: 0}, g.Vector{X: -7, Y: 50, Z: -12.5}, g.Vector{X: 3, Y: 9, Z: 10}, 27.76088375154269},
	{g.Vector{X: 1, Y: 1, Z: -2}, g.Vector{X: -100, Y: -9, Z: 1}, g.Vector{X: 4, Y: 2, Z: 9}, g.Vector{X: 10000, Y: 9999, Z: 10001}, 45.3155602414888},

	{g.Vector{X: -2, Y: 0, Z: 1}, g.Vector{X: -35.5, Y: 0, Z: -71}, g.Vector{X: 0, Y: 0, Z: 0}, g.Vector{X: 9, Y: 192.51, Z: 18}, 0},
	{g.Vector{X: -2, Y: 0, Z: 1}, g.Vector{X: -2, Y: 0, Z: 1}, g.Vector{X: 4, Y: 50, Z: -2}, g.Vector{X: 0.5, Y: 91.3, Z: -0.2}, 4.47213595499958},
	{g.Vector{X: -2, Y: 0, Z: 1}, g.Vector{X: 4, Y: -9, Z: 23}, g.Vector{X: 4.5, Y: 102, Z: 8.7}, g.Vector{X: -1.2, Y: -99, Z: -2}, 6.708203932499369},
}

func TestGetMaxDist(t *testing.T) {
	for _, element := range data {
		assert.InDelta(t, element.maxDist, GetMaxDist(element.plane, element.v1, element.v2, element.v3), delta, "MaxDist should be calculated correctly")
	}
}

func TestTripleCost(t *testing.T) {
	calc := CostCalculator{Threshold: 2, Amplification: 4}

	assert.InDelta(t, -8,
		calc.TripleCost(g.Vector{X: -35.5, Y: 0, Z: -71}, g.Vector{X: 0, Y: 0, Z: 0}, g.Vector{X: 9, Y: 192.51, Z: 18}),
		delta,
		"The points are all in the fitted plane",
	)
	assert.InDelta(t, -8,
		calc.TripleCost(g.Vector{X: 0, Y: 0, Z: 1}, g.Vector{X: 0, Y: 0, Z: 0}, g.Vector{X: 0, Y: 0, Z: 0}),
		delta,
		"The points are all in the origin",
	)
	assert.InDelta(t, -8,
		calc.TripleCost(g.Vector{X: 8, Y: 8, Z: 8}, g.Vector{X: -4.2, Y: -4.2, Z: -4.2}, g.Vector{X: 109.52, Y: 109.52, Z: 109.52}),
		delta,
		"The points are all in the fitted plane",
	)
	assert.InDelta(t, 1.7979589711327115,
		calc.TripleCost(g.Vector{X: 5, Y: 5, Z: 5}, g.Vector{X: 6, Y: 6, Z: 3}, g.Vector{X: 4, Y: 4, Z: 7}),
		delta,
		"The points are all in the fitted plane",
	)
}
