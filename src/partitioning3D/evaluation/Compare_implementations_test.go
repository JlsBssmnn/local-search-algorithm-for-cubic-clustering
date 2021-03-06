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

func init() {
	iterations = flag.Int("iterations", 5, "How many iterations should be executed to test algorithms for equality")
	seed = flag.Int("seed", 5, "The seed for the random number generation")
	randomizeParameters = flag.Bool("randomizeParameters", true, `If true the parameters
		threshold, numOfPlanes, pointsPerPlane, stddev, mean will be randomly chosen in each iteration, otherwise
		the parameters will be according to the command-line arguments`)
	algorithm1 = flag.String("algorithm1", "", "The first algorithm in the equality test")
	algorithm2 = flag.String("algorithm2", "", "The second algorithm in the equality test")
	verbose = flag.Int("verbose", 1, `Controls how much output is generated, higher levels include lower ones:
		0 - no output, 1 - print if an iteration was not successful, 2 - print if an iteration was successful
		3 - if partitions are not equal print the elements that are partitioned differently`)
}

func TestCompareAlgorithms(t *testing.T) {
	flag.Parse()
	if *algorithm1 == "" || *algorithm2 == "" {
		t.Log("Both algorithms have to be specified in order to compare them")
		return
	}

	firstAlgorithm := algorithm.AlgorithmStringToFunc[geometry.Vector](*algorithm1)
	secondAlgorithm := algorithm.AlgorithmStringToFunc[geometry.Vector](*algorithm2)

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
			t.Fail()
			if *verbose >= 1 {
				t.Errorf("Partitions were not equal at iteration %d", i)
			}
		} else if *verbose >= 2 {
			t.Log("Partitions were equal!")
		}
	}
}

func testForEquality(t *testing.T, firstAlgorithm, secondAlgorithm algorithm.PartitioningAlgorithm[geometry.Vector]) bool {
	success := true
	testData := GenerateDataWithNoise(*numOfPlanes, *pointsPerPlane, utils.NormalDist{Mean: *mean, Stddev: *stddev})
	calc := partitioning3D.CostCalculator{Threshold: *threshold, Amplification: *amplification}

	partAlg1 := firstAlgorithm(&testData.points, calc)
	partAlg2 := secondAlgorithm(&testData.points, calc)

	if len(partAlg1) != len(partAlg2) {
		t.Logf("Partitioning arrays have different lengths")
		return false
	}

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
				if *verbose >= 3 {
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
				if *verbose >= 3 {
					t.Logf("Element %d is differently partitioned", i)
				}
			} else if correspondingPart1 != part1 || correspondingPart2 != part2 {
				success = false
				if *verbose >= 3 {
					t.Logf("Element %d is differently partitioned", i)
				}
			}
		}
	}
	return success
}
