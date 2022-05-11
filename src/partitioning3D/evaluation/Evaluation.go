package evaluation

import (
	"math"

	alg "github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/algorithm"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"
)

// This is the function signiture of a partitioning algorithm for 3D Vectors
type PartAlgorithm func([]geometry.Vector, alg.CostCalculator[geometry.Vector]) alg.Partitioning[geometry.Vector]

// The output after an evaluation of an algorithm
type Evaluation struct {
	NumOfPlanesError float64
	Accuracy         float64
}

// Evaluate an algorithm by specifying the algorithm as string (the function name of the algorithm)
// the test data is generated with the help of the provided parameters. This function panics
// if the algorithm isn't defined
func Evaluate(algorithm string, costCalc alg.CostCalculator[geometry.Vector], numOfPlanes, pointsPerPlane int, noise utils.NormalDist) Evaluation {
	testData := GenerateDataWithNoise(numOfPlanes, pointsPerPlane, noise)

	switch algorithm {
	case "GreedyJoiningV1":
		return EvaluateAlgorithm(alg.GeedyJoiningV1[geometry.Vector], costCalc, testData)
	case "GreedyJoiningV2":
		return EvaluateAlgorithm(alg.GeedyJoiningV2[geometry.Vector], costCalc, testData)
	default:
		panic("This algorithm is not supported")
	}
}

// Evaluates a given algorithm with the given test data
func EvaluateAlgorithm(algorithm PartAlgorithm, costCalc alg.CostCalculator[geometry.Vector], testData TestData) Evaluation {
	part := algorithm(testData.points, costCalc)
	numOfPlanesError := math.Abs(float64(len(part)-testData.numOfPlanes)) / float64(testData.numOfPlanes)
	n := len(testData.points)

	pointsPerPlane := n / len(testData.planes)
	correctPartitioned := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			samePartition := i/pointsPerPlane == j/pointsPerPlane
			for _, partition := range part {
				iIn := utils.Contains(partition, testData.points[i])
				if !iIn {
					continue
				} else if utils.Contains(partition, testData.points[j]) {
					if samePartition {
						correctPartitioned++
					}
					break
				} else {
					if !samePartition {
						correctPartitioned++
					}
					break
				}
			}
		}
	}
	accuracy := float64(correctPartitioned) / (float64(n*(n-1)) / 2.0)

	return Evaluation{NumOfPlanesError: numOfPlanesError, Accuracy: accuracy}
}
