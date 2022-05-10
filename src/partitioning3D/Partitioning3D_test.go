package partitioning3D

import (
	"testing"

	g "github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
	"github.com/stretchr/testify/assert"
)

func TestFitPlane(t *testing.T) {
	plane := FitPlane(g.Vector{X: 5, Y: 5, Z: 5}, g.Vector{X: 6, Y: 6, Z: 3}, g.Vector{X: 4, Y: 4, Z: 7})
	assert.InDelta(t, 1, plane.X, delta, "P1 is on the plane, P2 and P3 are one normalvector away from P1")
	assert.InDelta(t, 1, plane.Y, delta, "P1 is on the plane, P2 and P3 are one normalvector away from P1")
	assert.InDelta(t, -2, plane.Z, delta, "P1 is on the plane, P2 and P3 are one normalvector away from P1")

	plane = FitPlane(g.Vector{X: 0, Y: 0, Z: 5}, g.Vector{X: 0, Y: 0, Z: -3}, g.Vector{X: 0, Y: 0, Z: 7})
	assert.InDelta(t, 1, plane.X, delta, "All points are on the z-axis")
	assert.InDelta(t, 0, plane.Y, delta, "All points are on the z-axis")
	assert.InDelta(t, 0, plane.Z, delta, "All points are on the z-axis")

	plane = FitPlane(g.Vector{X: 4, Y: 1, Z: 2}, g.Vector{X: -13, Y: 2, Z: -3}, g.Vector{X: 0, Y: 3, Z: 2})
	assert.InDelta(t, 1, plane.X, delta, "All points are one a plane that goes through the origin")
	assert.InDelta(t, 2, plane.Y, delta, "All points are one a plane that goes through the origin")
	assert.InDelta(t, -3, plane.Z, delta, "All points are one a plane that goes through the origin")
}
