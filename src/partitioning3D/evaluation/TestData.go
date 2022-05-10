package evaluation

import (
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
)

// This struct holds test data which can be used to evaluate an algorithm.
// numberOfPlanes describes how many planes are in the test data, planes is
// an array of the planes and points is an array of the sampled points.
type TestData struct {
	numOfPlanes int
	planes      []geometry.Vector
	points      []geometry.Vector
}
