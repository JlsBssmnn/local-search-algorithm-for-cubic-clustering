package evaluation

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"testing"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"
	"github.com/stretchr/testify/assert"
)

var delta = 0.00000001

var outputFile *string

func init() {
	outputFile = flag.String("outputFile", "", "The path and file name where the test data should be written to")
}

func TestGenerateDataWithoutNoise(t *testing.T) {
	nPlanes := 5
	pointsPerPlane := 20

	data := GenerateDataWithoutNoise(nPlanes, pointsPerPlane)
	assert.Equal(t, nPlanes, data.numOfPlanes, "number of planes should be stored correctly")
	assert.Equal(t, nPlanes, len(data.planes), "there should be as many planes as specified")

	for i, point := range data.points {
		assert.InDelta(t, 0, geometry.DistFromPlane(
			data.planes[int(math.Floor(float64(i)/float64(pointsPerPlane)))],
			point),
			delta,
			"Every point should be on it's corresponding plane",
		)
	}
}

func TestGenerateDataWithNoise(t *testing.T) {
	nPlanes := 5
	pointsPerPlane := 20
	noise := utils.NormalDist{Mean: 0, Stddev: 5}

	data := GenerateDataWithNoise(nPlanes, pointsPerPlane, noise)

	for i, point := range data.points {
		assert.LessOrEqual(t, 0.001, geometry.DistFromPlane(
			data.planes[int(math.Floor(float64(i)/float64(pointsPerPlane)))],
			point),
			delta,
			`It should basically be impossible for a point to be really close to it's
			 corresponding plane`,
		)
	}

	noise = utils.NormalDist{Mean: 0, Stddev: 0}
	data = GenerateDataWithNoise(nPlanes, pointsPerPlane, noise)

	for i, point := range data.points {
		assert.InDelta(t, 0.0, geometry.DistFromPlane(
			data.planes[int(math.Floor(float64(i)/float64(pointsPerPlane)))],
			point),
			delta,
			`If mu and sigma are 0 in the normal distribution, then every sampled point
			 should be on it's corresponding plane`,
		)
	}
}

func BenchmarkDataCreation(b *testing.B) {
	for _, v := range testTable {
		b.Run(fmt.Sprintf("numOfPlanes: %d, pointsPerPlane: %d", v.numOfPlanes, v.pointsPerPlane), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				GenerateDataWithNoise(v.numOfPlanes, v.pointsPerPlane, utils.NormalDist{Mean: 0, Stddev: 0})
			}
		})
	}
}

// This is not a unit test. It can be used to generate test data and save it to a file.
func TestSaveTestDataToFile(t *testing.T) {
	flag.Parse()
	if *outputFile == "" {
		return
	}
	testData := GenerateDataWithNoise(*numOfPlanes, *pointsPerPlane, utils.NormalDist{Mean: *mean, Stddev: *stddev})

	file, err := os.Create(*outputFile)
	defer file.Close()

	if err != nil {
		t.Error(err)
		return
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"X", "Y", "Z"}); err != nil {
		t.Error(err)
		return
	}
	for _, point := range testData.points {
		csvRow := []string{strconv.FormatFloat(point.X, 'f', 15, 64), strconv.FormatFloat(point.Y, 'f', 15, 64), strconv.FormatFloat(point.Z, 'f', 15, 64)}
		if err := writer.Write(csvRow); err != nil {
			t.Error(err)
			return
		}
	}
}
