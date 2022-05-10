package evaluation

import (
	"math"
	"testing"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"
	"github.com/stretchr/testify/assert"
)

var delta = 0.00000001

func TestGenerateDataWithoutNoise(t *testing.T) {
	nPlanes := 5
	pointsPerPlane := 20

	data := GenerateDataWithoutNoise(nPlanes, pointsPerPlane)
	assert.Equal(t, nPlanes, data.numOfPlanes, "number of planes should be stored correctly")
	assert.Equal(t, nPlanes, len(data.planes), "there should be as many planes as specified")

	for i, point := range data.points {
		assert.InDelta(t, 0, geometry.DistFromPlane(
			data.planes[int(math.Floor(float64(i)/float64(pointsPerPlane)))],
			point),
			delta,
			"Every point should be on it's corresponding plane",
		)
	}
}

func TestGenerateDataWithNoise(t *testing.T) {
	nPlanes := 5
	pointsPerPlane := 20
	noise := utils.NormalDist{Mean: 0, Stddev: 5}

	data := GenerateDataWithNoise(nPlanes, pointsPerPlane, noise)

	for i, point := range data.points {
		assert.LessOrEqual(t, 0.001, geometry.DistFromPlane(
			data.planes[int(math.Floor(float64(i)/float64(pointsPerPlane)))],
			point),
			delta,
			`It should basically be impossible for a point to be really close to it's 
			 corresponding plane`,
		)
	}
}
