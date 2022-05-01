package algorithm

import (
	"testing"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
	"github.com/stretchr/testify/assert"
)

func TestValidPartitioning(t *testing.T) {
	v1 := geometry.Vector{X: 1, Y: 1, Z: 1}
	v2 := geometry.Vector{X: 0.9, Y: 0.9, Z: 0.9}
	v3 := geometry.Vector{X: 1.1, Y: 1.1, Z: 1.1}

	assert.True(t, ValidPartitioning([][]geometry.Vector{{v1, v2, v3}}), "This is a valid partitioning")
	assert.False(t, ValidPartitioning([][]geometry.Vector{{v1, v1, v1}}), "This is not a valid partitioning")
}
