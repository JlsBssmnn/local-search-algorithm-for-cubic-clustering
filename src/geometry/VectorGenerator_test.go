package geometry

import (
	"testing"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"

	"github.com/stretchr/testify/assert"
)

func TestCreateRandomNormalVec(t *testing.T) {
	oldVector := Vector{}
	for i := 0; i < 10; i++ {
		v := CreateRandomUnitVec()
		assert.InDelta(t, 1.0, v.GetLength(), 0.00000001, "a random normal vector should have length 1")
		assert.False(t, false, EqualVector(oldVector, v), "random vectors should be different")
		oldVector = v
	}
}

func TestSamplePointFromPlane(t *testing.T) {
	for i := 0; i < 5; i++ {
		assert.NotPanics(t, func() {
			SamplePointFromPlane(CreateRandomUnitVec())
		}, "The function should not panic thus calculating a point with distance 0 to the plane")

		assert.NotPanics(t, func() {
			SamplePointFromPlane(CreateRandomVec())
		}, "The function should not panic thus calculating a point with distance 0 to the plane")
	}
}

func TestSamplePointFromPlaneWithNoise(t *testing.T) {
	plane := CreateRandomUnitVec()
	noise := utils.NormalDist{Mean: 0, Stddev: 1}

	point := SamplePointFromPlaneWithNoise(plane, noise)

	// an element is within 3 stddev with 99.7% probability
	assert.Less(t, DistFromPlane(plane, point), 3.0, "This test should fail with probability of 0.3%")

	// that an element is within 1 stddev in at least 1 of 10 iterations is ~0.001%
	inOneStddev := false
	for i := 0; i < 10; i++ {
		point = SamplePointFromPlaneWithNoise(plane, noise)
		inOneStddev = inOneStddev || DistFromPlane(plane, point) <= 1
	}
	assert.True(t, inOneStddev, "This test should fail with probability of 0.001%")

	// that an element is outside of 2 stddev in at least 1 of 200 iterations is ~0.0035%
	outside2Stddev := false
	for i := 0; i < 200; i++ {
		point = SamplePointFromPlaneWithNoise(plane, noise)
		outside2Stddev = outside2Stddev || DistFromPlane(plane, point) >= 2
	}
	assert.True(t, outside2Stddev, "This test should fail with probability of 0.0035%")
}
