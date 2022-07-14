package evaluation

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"testing"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/algorithm"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/partitioning3D"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"
	"github.com/stretchr/testify/assert"
)

var delta = 0.00000001

var outputFile, inputFile *string

func init() {
	outputFile = flag.String("outputFile", "", "The path and file name where the test data should be written to")
	inputFile = flag.String("inputFile", "", "The path and file name where the test data should be read from")
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

func TestDataFromPlanes(t *testing.T) {
	planes := []geometry.Vector{{X: 1, Y: 0, Z: 0}, {X: 0, Y: 1, Z: 0}, {X: 0, Y: 0, Z: 1}}
	testData := GenerateDataFromPlanesWithNoise(planes, 5, utils.NormalDist{Mean: 0, Stddev: 0})

	for i := 0; i < 5; i++ {
		point := testData.points[i]
		assert.Equal(t, 0.0, point.X)
		assert.True(t, point.Y != math.NaN() && point.Y != math.Inf(1) && point.Y != math.Inf(-1))
		assert.True(t, point.Z != math.NaN() && point.Z != math.Inf(1) && point.Z != math.Inf(-1))
	}
	for i := 5; i < 10; i++ {
		point := testData.points[i]
		assert.Equal(t, 0.0, point.Y)
		assert.True(t, point.X != math.NaN() && point.X != math.Inf(1) && point.X != math.Inf(-1))
		assert.True(t, point.Z != math.NaN() && point.Z != math.Inf(1) && point.Z != math.Inf(-1))
	}
	for i := 10; i < 15; i++ {
		point := testData.points[i]
		assert.Equal(t, 0.0, point.Z)
		assert.True(t, point.X != math.NaN() && point.X != math.Inf(1) && point.X != math.Inf(-1))
		assert.True(t, point.Y != math.NaN() && point.Y != math.Inf(1) && point.Y != math.Inf(-1))
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

// Writes the costs for the points in the given input file to a json file
func TestSaveCostToFile(t *testing.T) {
	flag.Parse()
	if *inputFile == "" || *outputFile == "" {
		return
	}

	points, err := partitioning3D.ParsePoints(*inputFile)
	if err != nil {
		t.Error(err)
		return
	}

	calc := partitioning3D.CostCalculator{Threshold: *threshold, Amplification: *amplification}
	alg := algorithm.CreateGreedyMovingAlgorithm[geometry.Vector](points, &calc, nil, nil, nil, nil, nil)
	alg.InitializeTripleCosts()

	tripleCost := alg.GetTripleCostArray()

	file, err := os.Create(*outputFile)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	jsonString, err := json.Marshal(tripleCost)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(*outputFile, jsonString, os.ModePerm)
}
