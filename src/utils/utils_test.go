package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSet(t *testing.T) {
	s1 := []int{1, 2, 3, 1, 6, 2, 3, 4}
	assert.Equal(t, 5, len(ToSet(s1)))

	s1 = []int{1, 2, 3, 6, 4}
	assert.Equal(t, 5, len(ToSet(s1)))

	s2 := []string{"this", "is", "a", "string", "slice", "without", "duplicates"}
	assert.Equal(t, 7, len(ToSet(s2)))

	s2 = []string{"this", "this", "a", "string", "slice", "slice", "string"}
	assert.Equal(t, 4, len(ToSet(s2)))
}
