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

func TestCompareAlgorithms(t *testing.T) {
	flag.Parse()
	rand.Seed(130)

	for i := 0; i < 25; i++ {
		*threshold = utils.RandomFloat(0.00001, 0.001)
		*numOfPlanes = utils.RandomInt(2, 5)
		*pointsPerPlane = utils.RandomInt(3, 17)
		*stddev = utils.RandomFloat(0.5, 1.5)
		*mean = utils.RandomFloat(-0.5, 0.5)

		success := testForEquality()
		if !success {
			t.Errorf("Partitions were not equal at interation %d", i)
		} else {
			t.Log("Partitions were equal!")
		}
	}
}

func testForEquality() bool {
	success := true
	testData := GenerateDataWithNoise(*numOfPlanes, *pointsPerPlane, utils.NormalDist{Mean: *mean, Stddev: *stddev})
	calc := partitioning3D.CostCalculator{Threshold: *threshold, Amplification: *amplification}

	partNew := algorithm.GreedyJoining[geometry.Vector](&testData.points, calc)

	partOld := algorithm.GreedyJoiningOld[geometry.Vector](testData.points, calc)

	partArrayOld := make([]int, len(partNew))
	for i, partition := range partOld {
		for _, point := range partition {
			for j, inputPoint := range testData.points {
				if point == inputPoint {
					partArrayOld[j] = i
				}
			}
		}
	}

	partMapping1 := make(map[int]int)
	partMapping2 := make(map[int]int)

	for i := range partNew {
		part1 := partNew[i]
		part2 := partArrayOld[i]

		correspondingPart2, ok := partMapping1[part1]
		if !ok {
			_, ok := partMapping2[part2]
			if ok {
				success = false
			} else {
				partMapping1[part1] = part2
				partMapping2[part2] = part1
			}
		} else {
			correspondingPart1, ok := partMapping2[part2]
			if !ok {
				success = false
			} else if correspondingPart1 != part1 || correspondingPart2 != part2 {
				success = false
			}
		}
	}
	return success
}
