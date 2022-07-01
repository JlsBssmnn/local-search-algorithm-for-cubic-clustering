package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	list := LinkedList[int]{}
	assert.Equal(t, 0, list.Length())

	assert.Nil(t, list.firstElement)
	assert.Nil(t, list.lastElement)

	list.Add(1)

	assert.Equal(t, 1, list.Length())
	assert.NotNil(t, list.firstElement)
	assert.NotNil(t, list.lastElement)
	assert.Equal(t, 1, list.firstElement.content)

	list.Add(6)

	assert.Equal(t, 2, list.Length())
	assert.NotNil(t, list.firstElement)
	assert.NotNil(t, list.lastElement)
	assert.Equal(t, 1, list.firstElement.content)
	assert.Equal(t, 6, list.lastElement.content)

	list.Add(4)

	assert.Equal(t, 3, list.Length())
	assert.NotNil(t, list.firstElement)
	assert.NotNil(t, list.lastElement)
	assert.Equal(t, 1, list.firstElement.content)
	assert.Equal(t, 6, list.firstElement.next.content)
	assert.Equal(t, 4, list.lastElement.content)
}

func TestIterator(t *testing.T) {
	list := LinkedList[int]{}
	list.Add(5)
	list.Add(3)
	list.Add(21)
	list.Add(32)
	list.Add(-2)

	iter := list.Iterator()
	assert.True(t, iter.HasNext())

	assert.Equal(t, 5, iter.Next())
	assert.Equal(t, 3, iter.Next())
	assert.Equal(t, 21, iter.Next())
	assert.Equal(t, 32, iter.Next())
	assert.Equal(t, -2, iter.Next())

	assert.False(t, iter.HasNext())
	assert.Panics(t, func() {
		iter.Next()
	})
}

func TestWithVectors(t *testing.T) {
	partitioningArray := []int{0, 0, 0, 0, 0, 0, 0, 1, 2}

	partitioning := make(map[int]*LinkedList[int])
	for i, partition := range partitioningArray {
		val, ok := partitioning[partition]
		if ok {
			length := val.Length()
			val.Add(i)
			assert.Equal(t, length+1, partitioning[partition].length)
		} else {
			list := LinkedList[int]{}
			list.Add(i)
			partitioning[partition] = &list
			assert.Equal(t, 1, partitioning[partition].length)
		}
	}
}

func TestAddList(t *testing.T) {
	l1 := CreateLinkedList(1, 2, 3, 4, 5)
	l2 := CreateLinkedList(6, 7, 8, 9, 10)

	l1.AddList(*l2)

	iter := l1.Iterator()
	i := 0
	for iter.HasNext() {
		i++
		assert.Equal(t, i, iter.Next())
	}
	assert.Equal(t, 10, i)
	assert.Equal(t, 10, l1.Length())

	iter2 := l2.Iterator()
	i = 5
	for iter2.HasNext() {
		i++
		assert.Equal(t, i, iter2.Next())
	}
	assert.Equal(t, 10, i)
	assert.Equal(t, 5, l2.Length())

	l3 := CreateLinkedList(11, 12, 13)
	l1.AddList(*l3)

	iter = l1.Iterator()
	i = 0
	for iter.HasNext() {
		i++
		assert.Equal(t, i, iter.Next())
	}
	assert.Equal(t, 13, i)
	assert.Equal(t, 13, l1.Length())
}
