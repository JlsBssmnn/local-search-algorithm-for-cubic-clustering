package algorithm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConstraints(t *testing.T) {
	allConstraints := parseConstraints("../../temp/constraint_files/constraints.json")
	assert.Equal(t, Edge{1, 2}, allConstraints.SamePartition[0])
	assert.Equal(t, Edge{4, 7}, allConstraints.SamePartition[1])
	assert.Equal(t, Edge{7, 9}, allConstraints.SamePartition[2])
	assert.Equal(t, Edge{9, 10}, allConstraints.SamePartition[3])

	assert.Equal(t, Edge{1, 3}, allConstraints.DifferentPartition[0])
	assert.Equal(t, Edge{1, 4}, allConstraints.DifferentPartition[1])
	assert.Equal(t, Edge{4, 8}, allConstraints.DifferentPartition[2])
	assert.Equal(t, Edge{8, 10}, allConstraints.DifferentPartition[3])

	constraints, partitions := translateConstraints(allConstraints, 11)

	assert.False(t, constraints.Get(1, 2))
	assert.False(t, constraints.Get(0, 1))
	assert.False(t, constraints.Get(3, 4))
	assert.False(t, constraints.Get(10, 1))

	assert.True(t, constraints.Get(1, 3))
	assert.True(t, constraints.Get(1, 4))
	assert.True(t, constraints.Get(8, 4))
	assert.True(t, constraints.Get(10, 8))

	assert.Equal(t, 2, len(partitions))
	for key, partition := range partitions {
		if partition.Length() == 1 {
			assert.Equal(t, 1, key)
			iter := partition.Iterator()
			assert.Equal(t, 2, iter.Next())
		} else if partition.Length() == 3 {
			assert.Equal(t, 4, key)
			iter := partition.Iterator()
			assert.Equal(t, 7, iter.Next())
			assert.Equal(t, 9, iter.Next())
			assert.Equal(t, 10, iter.Next())
		} else {
			t.Error("There was a precomputed partition which didn't match the expected size")
		}
	}
}

func TestInitializeSingletonSets(t *testing.T) {
	var part PartitioningArray
	part.InitializeSingletonSets(50)

	for i, v := range part {
		assert.Equal(t, i, v)
	}
}
