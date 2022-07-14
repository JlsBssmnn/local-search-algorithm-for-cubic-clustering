package main

import (
	"flag"
	"fmt"
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
	constraintFile := flag.String("constraintFile", "", "The path to a file which constraints constraints for a partitioning")

	flag.Parse()

	if *fileName == "" {
		panic("The path to a file with geometry data that should be partitioned must be provided as argument")
	} else if !strings.HasSuffix(*fileName, ".csv") {
		panic("Input file must be a csv file!")
	}

	points, err := partitioning3D.ParsePoints(*fileName)

	if err != nil {
		panic(err)
	}

	calc := partitioning3D.CostCalculator{Threshold: *threshold, Amplification: *amplification}

	var partitioningArray algorithm.PartitioningArray
	if *constraintFile != "" && *selectedAlgorithm == "GreedyMoving" {
		partitioningArray = algorithm.GreedyMovingWithConstraints[geometry.Vector](points, &calc, *constraintFile)
	} else {
		partitioningArray = algorithm.AlgorithmStringToFunc[geometry.Vector](*selectedAlgorithm)(points, &calc)
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
			fmt.Printf("X: %f, Y: %f, Z: %f\n", (*points)[point].X, (*points)[point].Y, (*points)[point].Z)
		}
		fmt.Println("--------------")
		i++
	}
}
