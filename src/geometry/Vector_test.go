package geometry

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

var lengthData = []struct {
	vector Vector
	length float64
}{
	{Vector{0, 0, 0}, 0},
	{Vector{1, 1, 1}, math.Sqrt(3)},
	{Vector{0.5, -0.7, 0.2}, math.Sqrt(0.25 + 0.49 + 0.04)},
	{Vector{2, 2, 1}, 3},
	{Vector{4, 3, 0}, 5},
}

func TestLength(t *testing.T) {
	for _, val := range lengthData {
		assert.Equal(t, val.length, val.vector.GetLength(), "The vectors should have same length")
	}
}

func TestEqual(t *testing.T) {
	for i, val1 := range lengthData {
		for j, val2 := range lengthData {
			if i == j {
				assert.True(t, EqualVector(val1.vector, val2.vector), "The vectors should be equal")
			} else {
				assert.False(t, EqualVector(val1.vector, val2.vector), "The vectors should not be equal")
			}
		}
	}
}
