package evaluation

import (
	"testing"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/partitioning3D"
	"github.com/stretchr/testify/assert"
)

func TestEvaluate(t *testing.T) {
	algorithm1 := func(testData TestData) partitioning3D.Output {
		return partitioning3D.Output{NumOfPlanes: testData.numOfPlanes}
	}
	algorithm2 := func(testData TestData) partitioning3D.Output {
		return partitioning3D.Output{NumOfPlanes: testData.numOfPlanes + 5}
	}

	testData := GenerateDataWithoutNoise(10, 10)

	evaluation1 := Evaluate(algorithm1, testData)
	evaluation2 := Evaluate(algorithm2, testData)

	assert.Equal(t, 0.0, evaluation1.NumOfPlanesError, "The first algorithm should have an error of 0 for the number of planes")
	assert.Equal(t, 0.5, evaluation2.NumOfPlanesError, "The second algorithm predicts 5 planes too much and should have to respective error")
}
