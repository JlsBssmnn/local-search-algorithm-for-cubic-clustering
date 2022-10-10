package evaluation

import (
	"flag"
	"testing"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/algorithm"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/partitioning3D"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"
)

func BenchmarkAlgorithm(b *testing.B) {
	flag.Parse()
	if *algorithm1 == "" {
		return
	}
	algorithm := algorithm.AlgorithmStringToFunc[geometry.Vector](*algorithm1)
	calc := partitioning3D.CostCalculator{Threshold: *threshold, Amplification: *amplification}
	noise := utils.NormalDist{Mean: *mean, Stddev: *stddev}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		testData := GenerateDataWithNoise(*numOfPlanes, *pointsPerPlane, noise)
		b.StartTimer()

		algorithm(&testData.Points, &calc)
	}
}
