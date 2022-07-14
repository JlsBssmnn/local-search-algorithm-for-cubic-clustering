package algorithm

import (
	"testing"

	g "github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
	"github.com/stretchr/testify/assert"
)

func GenerateData(numOfPlanes, pointsPerPlane int, sampling func(g.Vector) g.Vector) []g.Vector {
	var planes []g.Vector
	var points []g.Vector

	for i := 0; i < numOfPlanes; i++ {
		plane := g.CreateRandomUnitVec()
		planes = append(planes, plane)
		for j := 0; j < pointsPerPlane; j++ {
			points = append(points, sampling(plane))
		}
	}
	return points
}

type VectorCostCalculator struct {
	Threshold     float64
	Amplification float64
}

func (calc VectorCostCalculator) TripleCost(v1, v2, v3 g.Vector) float64 {
	v1.AddVector(v2)
	v1.AddVector(v3)
	return v1.X
}

var vectorAlgorithm GreedyJoiningAlgorithm[g.Vector]

const NUMBER_OF_PLANES = 5
const POINTS_PER_PLANE = 20

func init() {
	calc := VectorCostCalculator{Threshold: 0.1, Amplification: 1}
	testData := GenerateData(NUMBER_OF_PLANES, POINTS_PER_PLANE, g.SamplePointFromPlane)
	vectorAlgorithm = GetAlgorithm[g.Vector](&testData, &calc)
}

func TestLengthFirstDimVectors(t *testing.T) {
	nextJoin, cost := vectorAlgorithm.InitializeAlgorithm()
	iteration := 0

	for cost < 0 && nextJoin[0] != -1 && nextJoin[1] != -1 {
		assert.Equal(t, (NUMBER_OF_PLANES*POINTS_PER_PLANE)-1-iteration, len(*vectorAlgorithm.costs))
		nextJoin, cost = vectorAlgorithm.Join(nextJoin[0], nextJoin[1])
		iteration++
	}
}

func TestLengthSecondDimVectors(t *testing.T) {
	nextJoin, cost := vectorAlgorithm.InitializeAlgorithm()
	iteration := 0

	for cost < 0 && nextJoin[0] != -1 && nextJoin[1] != -1 {
		for i, onePartitionCosts := range *vectorAlgorithm.costs {
			assert.Equal(t, (NUMBER_OF_PLANES*POINTS_PER_PLANE)-1-iteration-i, len(onePartitionCosts.twoPartitionsCosts))
		}
		nextJoin, cost = vectorAlgorithm.Join(nextJoin[0], nextJoin[1])
		iteration++
	}
}
