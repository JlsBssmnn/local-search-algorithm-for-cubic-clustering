package evaluation

import (
	"math"

	alg "github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/algorithm"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"
)

// The output after an evaluation of an algorithm
type Evaluation struct {
	NumOfPlanesError float64
	Accuracy         float64
	FalsePositives   float64
	FalseNegatives   float64
}

// Evaluate an algorithm by specifying the algorithm as string (the function name of the algorithm)
// the test data is generated with the help of the provided parameters. This function panics
// if the algorithm isn't defined
func Evaluate(algorithm string, costCalc alg.CostCalculator[geometry.Vector], numOfPlanes, pointsPerPlane int, noise utils.NormalDist) Evaluation {
	testData := GenerateDataWithNoise(numOfPlanes, pointsPerPlane, noise)

	return EvaluateAlgorithm(alg.AlgorithmStringToFunc[geometry.Vector](algorithm), costCalc, &testData)
}

// Evaluates a given algorithm with the given test data
func EvaluateAlgorithm(algorithm alg.PartitioningAlgorithm[geometry.Vector], costCalc alg.CostCalculator[geometry.Vector], testData *TestData) Evaluation {
	part := algorithm(&testData.points, costCalc)
	numOfPlanesError := math.Abs(float64(len(utils.ToSet(part))-testData.numOfPlanes)) / float64(testData.numOfPlanes)
	n := len(testData.points)

	pointsPerPlane := n / len(testData.planes)
	correctPartitioned := 0
	fN := 0
	fP := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			samePartition := i/pointsPerPlane == j/pointsPerPlane
			if samePartition && part[i] == part[j] {
				correctPartitioned++
			} else if !samePartition && part[i] != part[j] {
				correctPartitioned++
			} else if samePartition {
				fP++
			} else {
				fN++
			}
		}
	}
	totalEdges := (float64(n*(n-1)) / 2.0)
	accuracy := float64(correctPartitioned) / totalEdges
	falsePositives := float64(fP) / (totalEdges - float64(correctPartitioned))
	falseNegatives := float64(fN) / (totalEdges - float64(correctPartitioned))

	return Evaluation{NumOfPlanesError: numOfPlanesError, Accuracy: accuracy, FalsePositives: falsePositives, FalseNegatives: falseNegatives}
}
