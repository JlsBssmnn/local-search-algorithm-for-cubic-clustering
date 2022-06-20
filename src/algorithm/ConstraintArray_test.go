package algorithm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	array := []bool{true, false, false, true, false, true, true, false, false, true}
	constraints := Constraints{array: array, numOfElements: 5}

	assert.False(t, constraints.Get(0, 2))
	assert.False(t, constraints.Get(0, 3))
	assert.False(t, constraints.Get(1, 2))
	assert.False(t, constraints.Get(2, 3))
	assert.False(t, constraints.Get(2, 4))

	assert.True(t, constraints.Get(0, 1))
	assert.True(t, constraints.Get(0, 4))
	assert.True(t, constraints.Get(1, 3))
	assert.True(t, constraints.Get(1, 4))
	assert.True(t, constraints.Get(3, 4))
}
