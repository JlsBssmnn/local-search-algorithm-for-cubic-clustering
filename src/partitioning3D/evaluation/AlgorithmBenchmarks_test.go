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
	{numOfPlanes: 3, pointsPerPlane: 3},
	{numOfPlanes: 3, pointsPerPlane: 30},
}

func BenchmarkGreedyJoiningV2WithoutNoise(b *testing.B) {
	calc := partitioning3D.CostCalculator{Threshold: 0.1, Amplification: 1}
	for _, v := range testTable {
		b.StopTimer()
		testData := GenerateDataWithoutNoise(v.numOfPlanes, v.pointsPerPlane)
		b.StartTimer()
		b.Run(fmt.Sprintf("numOfPlanes: %d, pointsPerPlane: %d", v.numOfPlanes, v.pointsPerPlane), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				algorithm.GeedyJoiningV2[geometry.Vector](testData.points, calc)
			}
		})
	}
}

func BenchmarkGreedyJoiningV2WithNoise(b *testing.B) {
	calc := partitioning3D.CostCalculator{Threshold: 0.1, Amplification: 1}
	noise := utils.NormalDist{Mean: 0, Stddev: 1}
	for _, v := range testTable {
		b.StopTimer()
		testData := GenerateDataWithNoise(v.numOfPlanes, v.pointsPerPlane, noise)
		b.StartTimer()
		b.Run(fmt.Sprintf("numOfPlanes: %d, pointsPerPlane: %d", v.numOfPlanes, v.pointsPerPlane), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				algorithm.GeedyJoiningV2[geometry.Vector](testData.points, calc)
			}
		})
	}
}

func BenchmarkCubicPartitionCost(b *testing.B) {
	calc := partitioning3D.CostCalculator{Threshold: 0.1, Amplification: 1}
	for _, v := range testTable {
		numOfPartitions := []int{v.numOfPlanes - 1, v.numOfPlanes, v.numOfPlanes}
		if v.pointsPerPlane == 500 {
			return
		}
		for _, n := range numOfPartitions {
			b.StopTimer()
			if n < 1 {
				continue
			}
			testData := GenerateDataWithoutNoise(v.numOfPlanes, v.pointsPerPlane)
			part := createPartitioning(testData.points, n)
			b.StartTimer()
			b.Run(fmt.Sprintf("numOfPlanes: %d, pointsPerPlane: %d, numOfPartitions: %d",
				v.numOfPlanes, v.pointsPerPlane, n), func(b *testing.B) {
				for j := 0; j < b.N; j++ {
					algorithm.CubicPartitionCost[geometry.Vector](part, calc)
				}
			})
		}
	}
}

// Helper function for creating a partitioning out of a slice of given points and
// the number of desired partitions
func createPartitioning(points []geometry.Vector, numOfPartitions int) [][]geometry.Vector {
	part := [][]geometry.Vector{}
	for i := 0; i < numOfPartitions; i++ {
		part = append(part, []geometry.Vector{})
	}
	for _, point := range points {
		index := utils.RandomInt(0, numOfPartitions)
		part[index] = append(part[index], point)
	}
	return part
}
