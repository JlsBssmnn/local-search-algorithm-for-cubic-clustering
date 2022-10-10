package main

import (
	"fmt"
	"syscall/js"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/algorithm"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/partitioning3D"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/partitioning3D/evaluation"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"
)

var currentData evaluation.TestData

// The parameters for a call from javascript are:
// generateData(numOfPlanes int, pointsPerPlane int, mean float64, stddev float64)
func generateData(this js.Value, inputs []js.Value) any {
	dist := utils.NormalDist{Mean: inputs[2].Float(), Stddev: inputs[3].Float()}
	currentData = evaluation.GenerateDataWithNoise(inputs[0].Int(), inputs[1].Int(), dist)
	return toWasmType(currentData)
}

func test(val ...int) {
	fmt.Println(val)
}

// The parameters for a call from javascript are:
// partitionData(algorithm string, threshold float64, amplification float64)
func partitionData(this js.Value, inputs []js.Value) any {
	calc := partitioning3D.CostCalculator{Threshold: inputs[1].Float(), Amplification: inputs[2].Float()}

	switch inputs[0].String() {
	case "GreedyMoving":
		return toWasmType(evaluation.EvaluateAlgorithm(algorithm.GreedyMoving[geometry.Vector], &calc, &currentData))
	case "GreedyJoining":
		return toWasmType(evaluation.EvaluateAlgorithm(algorithm.GreedyJoining[geometry.Vector], &calc, &currentData))
	default:
		fmt.Println("Provided algorithm not supported!")
	}
	return nil
}

func main() {
	c := make(chan int)
	js.Global().Set("generateData", js.FuncOf(generateData))
	js.Global().Set("partitionData", js.FuncOf(partitionData))
	<-c
}
