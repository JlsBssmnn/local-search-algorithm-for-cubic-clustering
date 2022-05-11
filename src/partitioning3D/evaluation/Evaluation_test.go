package evaluation

import (
	"testing"

	alg "github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/algorithm"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
	"github.com/stretchr/testify/assert"
)

type ZeroCostCalc struct{}

func (calc ZeroCostCalc) TripleCost(v1, v2, v3 geometry.Vector) float64 {
	return 0.0
}

func TestEvaluate(t *testing.T) {
	algorithm1 := func(data []geometry.Vector, costCalc alg.CostCalculator[geometry.Vector]) alg.Partitioning[geometry.Vector] {
		part := [][]geometry.Vector{}
		for i := 0; i < len(data); i += 10 {
			part = append(part, data[i:i+10])
		}
		return part
	}
	algorithm2 := func(data []geometry.Vector, costCalc alg.CostCalculator[geometry.Vector]) alg.Partitioning[geometry.Vector] {
		part := [][]geometry.Vector{}
		for i := 0; i < len(data)-12; i += 10 {
			if i == 80 {
				part = append(part, data[i:i+8])
				break
			}
			part = append(part, data[i:i+10])
		}
		for i := len(data) - 12; i < len(data); i += 2 {
			part = append(part, data[i:i+2])
		}
		return part
	}

	testData := GenerateDataWithoutNoise(10, 10)
	calc := ZeroCostCalc{}

	evaluation1 := EvaluateAlgorithm(algorithm1, calc, testData)
	evaluation2 := EvaluateAlgorithm(algorithm2, calc, testData)

	assert.Equal(t, 0.0, evaluation1.NumOfPlanesError, "The first algorithm should have an error of 0 for the number of planes")
	assert.Equal(t, 0.5, evaluation2.NumOfPlanesError, "The second algorithm predicts 5 planes too much and should have to respective error")

	assert.Equal(t, 1.0, evaluation1.Accuracy, "The first algorithm partitions everything correctly")
	assert.InDelta(t, 0.9886868686868687, evaluation2.Accuracy, delta, "The second algorithm doesn't partition everything correct")
}
