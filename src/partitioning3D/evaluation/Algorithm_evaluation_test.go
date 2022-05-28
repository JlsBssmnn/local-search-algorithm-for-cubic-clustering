package evaluation

import (
	"flag"
	"fmt"
	"testing"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/algorithm"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/partitioning3D"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"
)

var threshold, amplification, mean, stddev *float64
var numOfPlanes, pointsPerPlane *int

func init() {
	threshold = flag.Float64("threshold", 1.0, "The threshold for the cost calculation")
	amplification = flag.Float64("amplification", 1.0, "The amplification for the cost calculation")
	mean = flag.Float64("mean", 0, "The mean for the noise")
	stddev = flag.Float64("stddev", 1.0, "The standard deviation for the noise")

	numOfPlanes = flag.Int("numberOfPlanes", 5, "How many planes should be used to sample data points")
	pointsPerPlane = flag.Int("pointsPerPlane", 5, "How many points per plane should be sampled")
}

func TestEvalGreedyJoining(t *testing.T) {
	flag.Parse()
	testData := GenerateDataWithoutNoise(*numOfPlanes, *pointsPerPlane)
	eval := EvaluateAlgorithm(algorithm.GeedyJoining[geometry.Vector], partitioning3D.CostCalculator{Threshold: *threshold, Amplification: *amplification}, &testData)

	fmt.Printf("Greedy Joining on %d planes with %d points per plane gave the following results:\n", *numOfPlanes, *pointsPerPlane)
	fmt.Printf("\tnumber of planes error: %f%%\n\taccuracy: %f%%\n", eval.NumOfPlanesError*100, eval.Accuracy*100)
}

func TestEvalGreedyJoiningWithNoise(t *testing.T) {
	flag.Parse()
	testData := GenerateDataWithNoise(*numOfPlanes, *pointsPerPlane, utils.NormalDist{Mean: *mean, Stddev: *stddev})
	eval := EvaluateAlgorithm(algorithm.GeedyJoining[geometry.Vector], partitioning3D.CostCalculator{Threshold: *threshold, Amplification: *amplification}, &testData)

	fmt.Printf("Greedy Joining on %d planes with %d points per plane gave the following results:\n", *numOfPlanes, *pointsPerPlane)
	fmt.Printf("\tnumber of planes error: %f%%\n\taccuracy: %f%%\n", eval.NumOfPlanesError*100, eval.Accuracy*100)
}
