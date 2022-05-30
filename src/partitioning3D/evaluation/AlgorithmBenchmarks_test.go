package evaluation

import (
	"fmt"
	"testing"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/algorithm"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/partitioning3D"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"
)

var testTable = []struct {
	numOfPlanes, pointsPerPlane int
}{
	{numOfPlanes: 1, pointsPerPlane: 3},
	{numOfPlanes: 2, pointsPerPlane: 3},
	{numOfPlanes: 2, pointsPerPlane: 4},
	{numOfPlanes: 2, pointsPerPlane: 5},
	{numOfPlanes: 2, pointsPerPlane: 6},
	{numOfPlanes: 2, pointsPerPlane: 30},
	{numOfPlanes: 2, pointsPerPlane: 50},
	{numOfPlanes: 2, pointsPerPlane: 100},
	{numOfPlanes: 2, pointsPerPlane: 500},
	{numOfPlanes: 3, pointsPerPlane: 3},
	{numOfPlanes: 3, pointsPerPlane: 30},
	{numOfPlanes: 3, pointsPerPlane: 300},
	{numOfPlanes: 4, pointsPerPlane: 30},
	{numOfPlanes: 4, pointsPerPlane: 250},
}

func BenchmarkGreedyJoiningWithoutNoise(b *testing.B) {
	calc := partitioning3D.CostCalculator{Threshold: 0.1, Amplification: 1}
	for _, v := range testTable {
		b.StopTimer()
		testData := GenerateDataWithoutNoise(v.numOfPlanes, v.pointsPerPlane)
		b.StartTimer()
		b.Run(fmt.Sprintf("numOfPlanes: %d, pointsPerPlane: %d", v.numOfPlanes, v.pointsPerPlane), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				algorithm.GreedyJoining[geometry.Vector](&testData.points, &calc)
			}
		})
	}
}

func BenchmarkGreedyJoiningWithNoise(b *testing.B) {
	calc := partitioning3D.CostCalculator{Threshold: 0.1, Amplification: 1}
	noise := utils.NormalDist{Mean: 0, Stddev: 1}
	for _, v := range testTable {
		b.StopTimer()
		testData := GenerateDataWithNoise(v.numOfPlanes, v.pointsPerPlane, noise)
		b.StartTimer()
		b.Run(fmt.Sprintf("numOfPlanes: %d, pointsPerPlane: %d", v.numOfPlanes, v.pointsPerPlane), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				algorithm.GreedyJoining[geometry.Vector](&testData.points, &calc)
			}
		})
	}
}
