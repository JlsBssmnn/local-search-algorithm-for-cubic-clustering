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
	return TestData{NumOfPlanes: numOfPlanes, Planes: planes, Points: points}
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

// Samples the specified number of points from each of the given planes with noise and
// returns everything as a test data struct
func GenerateDataFromPlanesWithNoise(planes []g.Vector, pointsPerPlane int, noise utils.NormalDist) TestData {
	points := make([]g.Vector, 0, len(planes)*pointsPerPlane)
	for _, plane := range planes {
		for j := 0; j < pointsPerPlane; j++ {
			points = append(points, g.SamplePointFromPlaneWithNoise(plane, noise))
		}
	}
	return TestData{NumOfPlanes: len(planes), Planes: planes, Points: points}
}
