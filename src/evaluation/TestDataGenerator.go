package evaluation

import (
	g "github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"
)

func GenerateData(numOfPlanes, pointsPerPlane int, sampling func(g.Vector) g.Vector) TestData {
	var planes []g.Vector
	var points []g.Vector

	for i := 0; i < numOfPlanes; i++ {
		plane := g.CreateRandomUnitVec()
		planes = append(planes, plane)
		for j := 0; j < pointsPerPlane; j++ {
			points = append(points, sampling(plane))
		}
	}
	return TestData{numOfPlanes: numOfPlanes, planes: planes, points: points}
}

// Generate test data without noise, thus every point will be on one plane
func GenerateDataWithoutNoise(numOfPlanes, pointsPerPlane int) TestData {
	return GenerateData(numOfPlanes, pointsPerPlane, g.SamplePointFromPlane)
}

// Generate test data with gaussian noise, every point will have a distance to it's
// plane which is sampled out of a normal distribution
func GenerateDataWithNoise(numOfPlanes, pointsPerPlane int, noise utils.NormalDist) TestData {
	sampling := func(plane g.Vector) g.Vector {
		return g.SamplePointFromPlaneWithNoise(plane, noise)
	}
	return GenerateData(numOfPlanes, pointsPerPlane, sampling)
}
