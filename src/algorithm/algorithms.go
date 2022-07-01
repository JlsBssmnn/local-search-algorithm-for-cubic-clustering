package algorithm

import "fmt"

// This file maps the names of the algorithms to the actual functions
func AlgorithmStringToFunc[data any](algorithm string) PartitioningAlgorithm[data] {
	switch algorithm {
	case "":
		panic("The algorithm was not specified")
	case "GreedyJoining":
		return GreedyJoining[data]
	case "GreedyMoving":
		return GreedyMoving[data]
	case "NaiveGreedyJoining":
		return NaiveGreedyJoining[data]
	case "NaiveGreedyMoving":
		return NaiveGreedyMoving[data]
	default:
		panic(fmt.Sprintf("Algorithm %s not supported", algorithm))
	}
}
