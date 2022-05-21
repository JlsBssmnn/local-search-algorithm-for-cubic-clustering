package utils

import (
	"math/rand"
	"sort"
)

type NormalDist struct {
	Mean, Stddev float64
}

type Numeric interface {
	int | int64 | uint | int32 | int16 | int8 | float32 | float64
}

// Return a random double in the given range, min is
// inclusive, max is exclusive
func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// Return a random integer in the given range, min is
// inclusive, max is exclusive
func RandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func FloatFromNormalDist(noise NormalDist) float64 {
	return rand.NormFloat64()*noise.Stddev + noise.Mean
}

// Check if an element is in the given slice
func Contains[T comparable](slice []T, element T) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}

// Returns the smallest element of the given input slice
func Min[T Numeric](slice []T) T {
	smallest := slice[0]
	for _, element := range slice {
		if element < smallest {
			smallest = element
		}
	}
	return smallest
}

// Returns the index of the smallest element in the slice.
// If multiple elements are the smallest then the index of the
// first element will be returned.
func ArgMin[T Numeric](slice []T) int {
	smallest := slice[0]
	index := 0
	for i, element := range slice {
		if element < smallest {
			smallest = element
			index = i
		}
	}
	return index
}

// Returns minimum as well as arg min
func MinAndArgMin[T Numeric](slice []T) (T, int) {
	if len(slice) == 0 {
		return 0, -1
	}
	smallest := slice[0]
	index := 0
	for i, element := range slice {
		if element < smallest {
			smallest = element
			index = i
		}
	}
	return smallest, index
}

// Checks if all elements in the given slice are different
func AllDifferent[T comparable](slice []T) bool {
	for i := 0; i < len(slice)-1; i++ {
		for j := i + 1; j < len(slice); j++ {
			if slice[i] == slice[j] {
				return false
			}
		}
	}

	return true
}

// Checks if the given function returns true for every application
// on a slice element
func All[T any](slice []T, function func(elemenet T) bool) bool {
	for _, element := range slice {
		if !function(element) {
			return false
		}
	}
	return true
}

// Checks if the given function returns true for any application
// on a slice element
func Any[T any](slice []T, function func(elemenet T) bool) bool {
	for _, element := range slice {
		if function(element) {
			return true
		}
	}
	return false
}

// Transforms one slice to another by applying the map function
// to every element
func Map[T any, D any](slice []T, function func(element T) D) []D {
	newSlice := make([]D, len(slice))
	for i, e := range slice {
		newSlice[i] = function(e)
	}
	return newSlice
}

// Sum up all elements in the slice
func Sum[T Numeric](slice []T) T {
	var sum T = 0
	for _, element := range slice {
		sum += element
	}
	return sum
}

// Sums up elements of the slice after mapping them to some numeric value
// with the given function.
func MapSum[T any, D Numeric](slice []T, mapping func(element T) D) D {
	var sum D = 0
	for _, element := range slice {
		sum += mapping(element)
	}
	return sum
}

// Takes an arbitrary number of int pointers and re-assigns their values
// such that the values are ascending in the order of the input parameters
func SortInts(p ...*int) {
	elements := make([]int, len(p))
	for i, e := range p {
		elements[i] = *e
	}
	sort.Slice(elements, func(i, j int) bool {
		return elements[i] < elements[j]
	})
	for i, pointer := range p {
		*pointer = elements[i]
	}
}

// Convert a slice to a set, aka remove any duplicates from a slice
func ToSet[T comparable](slice []T) []T {
	store := make(map[T]bool)
	set := []T{}
	for _, element := range slice {
		_, ok := store[element]
		if !ok {
			store[element] = true
			set = append(set, element)
		}
	}
	return set
}
