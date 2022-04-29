package evaluation

import (
	"math"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/algorithm"
)

type Evaluation struct {
	NumOfPlanesError float64
}

func Evaluate(algorithm func(TestData) algorithm.Output, testData TestData) Evaluation {
	output := algorithm(testData)

	numOfPlanesError := math.Abs(float64(output.NumOfPlanes)-float64(testData.numOfPlanes)) / float64(testData.numOfPlanes)

	return Evaluation{NumOfPlanesError: numOfPlanesError}
}
