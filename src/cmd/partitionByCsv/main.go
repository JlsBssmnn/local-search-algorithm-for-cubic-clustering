package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/algorithm"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/partitioning3D"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"
)

func main() {
	fileName := flag.String("fileName", "", "The path to the csv file with the input data")
	threshold := flag.Float64("threshold", 1.0, "The threshold for the cost calculation")
	amplification := flag.Float64("amplification", 1.0, "The amplification for the cost calculation")
	selectedAlgorithm := flag.String("algorithm", "", "The algorithm which should be used for the partitioning")
	constraintFile := flag.String("constraintFile", "", "The path to a file which containts constraints for a partitioning")

	flag.Parse()

	if *fileName == "" {
		panic("The path to a file with geometry data that should be partitioned must be provided as argument")
	} else if !strings.HasSuffix(*fileName, ".csv") {
		panic("Input file must be a csv file!")
	}

	points := parsePoints(*fileName)
	calc := partitioning3D.CostCalculator{Threshold: *threshold, Amplification: *amplification}

	var partitioningArray algorithm.PartitioningArray
	switch *selectedAlgorithm {
	case "GreedyJoining":
		partitioningArray = algorithm.GreedyJoining[geometry.Vector](&points, &calc)
	case "GreedyMoving":
		if *constraintFile != "" {
			partitioningArray = algorithm.GreedyMovingWithConstraints[geometry.Vector](&points, &calc, *constraintFile)
		} else {
			partitioningArray = algorithm.GreedyMoving[geometry.Vector](&points, &calc)
		}
	default:
		panic("The provided algorithm isn't supported!")
	}

	// Order elements by their partition
	partitioning := make(map[int]*utils.LinkedList[int])
	for i, partition := range partitioningArray {
		val, ok := partitioning[partition]
		if ok {
			val.Add(i)
		} else {
			list := utils.LinkedList[int]{}
			list.Add(i)
			partitioning[partition] = &list
		}
	}

	// Output partitioning
	fmt.Println("--------------")
	i := 0
	for _, elements := range partitioning {
		fmt.Printf("Partition_%d\n", i)
		iterator := elements.Iterator()
		for iterator.HasNext() {
			point := iterator.Next()
			fmt.Printf("X: %f, Y: %f, Z: %f\n", points[point].X, points[point].Y, points[point].Z)
		}
		fmt.Println("--------------")
		i++
	}
}

// Parses the given csv file into a slice of vectors
// if this fails this function panics.
func parsePoints(fileName string) []geometry.Vector {
	file, err := os.Open(fileName)
	if err != nil {
		panic("The specified file wasn't found!")
	}
	reader := csv.NewReader(file)

	row, err := reader.Read()
	if err == io.EOF {
		os.Exit(0)
	} else if err != nil {
		panic(err)
	} else if len(row) < 3 {
		panic("The csv file must contain at least 3 columns for x, y and z coordinates!")
	}

	data := []geometry.Vector{}

	// find out column indicies of x, y and z coordinate
	xIdx, yIdx, zIdx := -1, -1, -1
	for i, value := range row {
		switch strings.ToLower(value) {
		case "x":
			xIdx = i
		case "y":
			yIdx = i
		case "z":
			zIdx = i
		}
	}
	if xIdx+yIdx+zIdx == -3 {
		// the first row doesn't contain x, y or z, so assume this row
		// already contains data and the first row is the x-coordinate,
		// second is the y-coordinate and third row is z-coordinate
		xIdx = 0
		yIdx = 1
		zIdx = 2
		data = append(data, rowToVector(row[xIdx], row[yIdx], row[zIdx]))
	} else if xIdx == -1 || yIdx == -1 || zIdx == -1 {
		// some x, y or z is specified but not all so panic
		panic("Csv head doesn't specify each of the x, y and z coordinates!")
	}

	for {
		row, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		data = append(data, rowToVector(row[xIdx], row[yIdx], row[zIdx]))
	}
	return data
}

// Try to convert 3 string into a vector by converting each
// into a float. It this fails the function panics.
func rowToVector(v1, v2, v3 string) geometry.Vector {
	x, errX := strconv.ParseFloat(v1, 64)
	y, errY := strconv.ParseFloat(v2, 64)
	z, errZ := strconv.ParseFloat(v3, 64)

	if errX != nil || errY != nil || errZ != nil {
		panic("Couldn't convert a cell in the csv into a float")
	}

	return geometry.Vector{X: x, Y: y, Z: z}
}
