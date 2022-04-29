package geometry

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

var delta = 0.00000001

var distanceData = []struct {
	planeNormalVector Vector
	pointVector       Vector
	distance          float64
}{
	{planeNormalVector: Vector{1, 0, 0}, pointVector: Vector{0, 0, 0}, distance: 0},
	{planeNormalVector: Vector{1, 0, 0}, pointVector: Vector{-10.2, 0.8, 0}, distance: 10.2},
	{planeNormalVector: Vector{1, 0, 0}, pointVector: Vector{5, 0, -4.2}, distance: 5},
	{planeNormalVector: Vector{1, 0, 0}, pointVector: Vector{1.2, 11.9, -8.1}, distance: 1.2},
	{planeNormalVector: Vector{1, 0, 0}, pointVector: Vector{0, -0.2, -0.01}, distance: 0},

	{planeNormalVector: Vector{math.Sqrt(0.5), math.Sqrt(0.5), 0}, pointVector: Vector{-4.2, 4.2, 0}, distance: 0},
	{planeNormalVector: Vector{math.Sqrt(0.5), math.Sqrt(0.5), 0}, pointVector: Vector{-11.9, 11.9, 4}, distance: 0},
	{planeNormalVector: Vector{math.Sqrt(0.5), math.Sqrt(0.5), 0}, pointVector: Vector{1, 1, 0}, distance: math.Sqrt(2)},
	{planeNormalVector: Vector{math.Sqrt(0.5), math.Sqrt(0.5), 0}, pointVector: Vector{8, -2, 10.23}, distance: math.Sqrt(18)},
}

func TestDistFromPlane(t *testing.T) {
	for _, val := range distanceData {
		assert.InDelta(t, val.distance, DistFromPlane(val.planeNormalVector, val.pointVector), delta, "Distance should be correct")
	}
}
