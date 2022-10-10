package evaluation

import (
	"math"

	alg "github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/algorithm"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/partitioning3D"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"
)

// The output after an evaluation of an algorithm
// A `positive` is an edge in the multicut, so an edge which was cut.
// Analogously, an edge which was not cut, is a negative.
type Evaluation struct {
	NumOfPlanesError float64
	Accuracy         float64
	TotalEdges       int
	TruePositives    int
	TrueNegatives    int
	FalsePositives   int
	FalseNegatives   int
	ComputedPlanes   []geometry.Vector
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
	part := algorithm(&testData.Points, costCalc)
	numOfPlanes := len(utils.ToSet(part))
	numOfPlanesError := math.Abs(float64(numOfPlanes-testData.NumOfPlanes)) / float64(testData.NumOfPlanes)
	n := len(testData.Points)

	pointsPerPlane := n / len(testData.Planes)
	correctPartitioned := 0
	tN := 0
	tP := 0
	fN := 0
	fP := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			samePartition := i/pointsPerPlane == j/pointsPerPlane
			if samePartition && part[i] == part[j] {
				correctPartitioned++
				tN++
			} else if !samePartition && part[i] != part[j] {
				correctPartitioned++
				tP++
			} else if samePartition {
				fP++
			} else {
				fN++
			}
		}
	}
	totalEdges := (float64(n*(n-1)) / 2.0)
	accuracy := float64(correctPartitioned) / totalEdges

	// compute planes
	partitioning := make(map[int][]*geometry.Vector, numOfPlanes)
	for i, partition := range part {
		_, ok := partitioning[partition]
		if ok {
			partitioning[partition] = append(partitioning[partition], &testData.Points[i])
		} else {
			partitioning[partition] = []*geometry.Vector{&testData.Points[i]}
		}
	}

	computedPlanes := make([]geometry.Vector, 0, numOfPlanes)
	for _, points := range partitioning {
		computedPlanes = append(computedPlanes, partitioning3D.FitPlane(points...))
	}
	return Evaluation{
		NumOfPlanesError: numOfPlanesError,
		Accuracy:         accuracy,
		TotalEdges:       int(totalEdges),
		TruePositives:    tP,
		TrueNegatives:    tN,
		FalsePositives:   fP,
		FalseNegatives:   fN,
		ComputedPlanes:   computedPlanes,
	}
}
