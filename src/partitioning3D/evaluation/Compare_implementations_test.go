package evaluation

import (
	"flag"
	"math/rand"
	"testing"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/algorithm"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/partitioning3D"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"
)

var iterations, seed, verbose *int
var randomizeParameters *bool
var algorithm1, algorithm2 *string

type AlgorithmFunction func(input *[]geometry.Vector, calc algorithm.CostCalculator[geometry.Vector]) algorithm.PartitioningArray

func init() {
	iterations = flag.Int("iterations", 5, "How many iterations should be executed to test algorithms for equality")
	seed = flag.Int("seed", 5, "The seed for the random number generation")
	randomizeParameters = flag.Bool("randomizeParameters", true, `If true the parameters
		threshold, numOfPlanes, pointsPerPlane, stddev, mean will be randomly choosen in each iteration, otherwise
		the parameters will be according to the command-line arguments`)
	algorithm1 = flag.String("algorithm1", "", "The first algorithm in the equality test")
	algorithm2 = flag.String("algorithm2", "", "The second algorithm in the equality test")
	verbose = flag.Int("verbose", 1, `Controls how much output is generated: 0 - no output, 1 - success for each iteration
		2 - if partitions are not equal print the elements that are partitioned differently`)
}

func TestCompareAlgorithms(t *testing.T) {
	flag.Parse()
	if *algorithm1 == "" || *algorithm2 == "" {
		t.Error("Both algorithms have to be specified in order to compare them")
		t.FailNow()
	}

	var firstAlgorithm, secondAlgorithm AlgorithmFunction

	switch *algorithm1 {
	case "GreedyJoining":
		firstAlgorithm = algorithm.GreedyJoining[geometry.Vector]
	case "GreedyMoving":
		firstAlgorithm = algorithm.GreedyMoving[geometry.Vector]
	case "ImprovedGreedyMoving":
		firstAlgorithm = algorithm.ImprovedGreedyMoving[geometry.Vector]
	default:
		t.Error("The specified algorithm1 is not supported")
		t.FailNow()
	}

	switch *algorithm2 {
	case "GreedyJoining":
		secondAlgorithm = algorithm.GreedyJoining[geometry.Vector]
	case "GreedyMoving":
		secondAlgorithm = algorithm.GreedyMoving[geometry.Vector]
	case "ImprovedGreedyMoving":
		secondAlgorithm = algorithm.ImprovedGreedyMoving[geometry.Vector]
	default:
		t.Error("The specified algorithm2 is not supported")
		t.FailNow()
	}
	rand.Seed(int64(*seed))

	for i := 0; i < *iterations; i++ {
		if *randomizeParameters {
			*threshold = utils.RandomFloat(0.00001, 0.001)
			*numOfPlanes = utils.RandomInt(2, 5)
			*pointsPerPlane = utils.RandomInt(3, 17)
			*stddev = utils.RandomFloat(0.5, 1.5)
			*mean = utils.RandomFloat(-0.5, 0.5)
		}

		success := testForEquality(t, firstAlgorithm, secondAlgorithm)
		if !success {
			t.Errorf("Partitions were not equal at interation %d", i)
		} else {
			t.Log("Partitions were equal!")
		}
	}
}

func testForEquality(t *testing.T, firstAlgorithm, secondAlgorithm AlgorithmFunction) bool {
	success := true
	testData := GenerateDataWithNoise(*numOfPlanes, *pointsPerPlane, utils.NormalDist{Mean: *mean, Stddev: *stddev})
	calc := partitioning3D.CostCalculator{Threshold: *threshold, Amplification: *amplification}

	partAlg1 := firstAlgorithm(&testData.points, calc)
	partAlg2 := secondAlgorithm(&testData.points, calc)

	partMapping1 := make(map[int]int)
	partMapping2 := make(map[int]int)

	for i := range partAlg1 {
		part1 := partAlg1[i]
		part2 := partAlg2[i]

		correspondingPart2, ok := partMapping1[part1]
		if !ok {
			_, ok := partMapping2[part2]
			if ok {
				success = false
				if *verbose >= 2 {
					t.Logf("Element %d is differently partitioned", i)
				}
			} else {
				partMapping1[part1] = part2
				partMapping2[part2] = part1
			}
		} else {
			correspondingPart1, ok := partMapping2[part2]
			if !ok {
				success = false
				if *verbose >= 2 {
					t.Logf("Element %d is differently partitioned", i)
				}
			} else if correspondingPart1 != part1 || correspondingPart2 != part2 {
				success = false
				if *verbose >= 2 {
					t.Logf("Element %d is differently partitioned", i)
				}
			}
		}
	}
	return success
}
