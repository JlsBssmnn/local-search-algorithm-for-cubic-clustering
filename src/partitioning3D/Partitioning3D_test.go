package partitioning3D

import (
	"testing"

	g "github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
	"github.com/stretchr/testify/assert"
)

func TestFitPlane(t *testing.T) {
	p1 := g.Vector{X: 1, Y: 0, Z: 0}
	p2 := g.Vector{X: 0, Y: 1, Z: 0}
	p3 := g.Vector{X: 0, Y: 0, Z: 1}
	plane := FitPlane(&p1, &p2, &p3)
	assert.Less(t, 0.1, g.DistFromPlane(&plane, &p1), delta, "There is no plane which contains all the points")
	assert.Less(t, 0.1, g.DistFromPlane(&plane, &p1), delta, "There is no plane which contains all the points")
	assert.Less(t, 0.1, g.DistFromPlane(&plane, &p1), delta, "There is no plane which contains all the points")

	plane = FitPlane(&g.Vector{X: 0, Y: 0, Z: 5}, &g.Vector{X: 0, Y: 0, Z: -3}, &g.Vector{X: 0, Y: 0, Z: 7})
	assert.InDelta(t, 1, plane.X, delta, "All points are on the z-axis")
	assert.InDelta(t, 0, plane.Y, delta, "All points are on the z-axis")
	assert.InDelta(t, 0, plane.Z, delta, "All points are on the z-axis")

	p1 = g.Vector{X: 4, Y: 1, Z: 2}
	p2 = g.Vector{X: -13, Y: 2, Z: -3}
	p3 = g.Vector{X: 0, Y: 3, Z: 2}

	plane = FitPlane(&p1, &p2, &p3)
	assert.InDelta(t, 0, g.DistFromPlane(&plane, &p1), delta, "All points are one a plane that goes through the origin")
	assert.InDelta(t, 0, g.DistFromPlane(&plane, &p2), delta, "All points are one a plane that goes through the origin")
	assert.InDelta(t, 0, g.DistFromPlane(&plane, &p3), delta, "All points are one a plane that goes through the origin")
}
