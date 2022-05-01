package algorithm

import (
	"fmt"
	"math"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
	"gonum.org/v1/gonum/mat"
)

type Output struct {
	NumOfPlanes int
	Planes      []geometry.Vector
}

// Represents one partitioning of 3D points. First dimension are the
// different partitions, second dimension are the elements in one partition.
// Example partitioning: [[a,b,c], [d,e,f,g,h], [i,j]]
type Partitioning [][]geometry.Vector

// This function checks if a given partitioning is valid by checking if
// it contains duplicates (doesn't work yet)
func ValidPartitioning(partitioning Partitioning) bool {
	allKeys := make(map[*geometry.Vector]bool)
	for i := 0; i < len(partitioning); i++ {
		for j := 0; j < len(partitioning[i]); j++ {
			fmt.Printf("current item has address %p", &partitioning[i][j])
			if _, ok := allKeys[&partitioning[i][j]]; !ok {
				allKeys[&partitioning[i][j]] = true
			} else {
				fmt.Println(allKeys)
				return false
			}
		}
	}
	return true
}

// The greedy joining algorithm
func GeedyJoining(data []geometry.Vector, calc CostCalculator) Output {
	// the initial partitioning starts with singleton sets
	var part Partitioning = [][]geometry.Vector{}
	for _, point := range data {
		part = append(part, []geometry.Vector{point})
	}
	cost := calc.CubicPartitionCost(part)
	newCost := math.Inf(-1)

	for {
		// find {B,C} that minimize the cost when joining them
		newCost = math.Inf(1)
		B, C := -1, -1
		for i := 0; i < len(part); i++ {
			for j := i + 1; j < len(part); j++ {
				join := append(part[i], part[j]...)
				newPart := append(append(append(part[:i], part[i+1:j]...), part[j+1:]...), join)
				joinCost := calc.CubicPartitionCost(newPart)

				if joinCost < newCost {
					newCost = joinCost
					B = i
					C = j
				}
			}
		}
		if B == -1 || C == -1 {
			panic("All join costs are infinite")
		}

		if newCost-cost < 0 {
			cost = newCost
			join := append(part[B], part[C]...)
			part = append(append(append(part[:B], part[B+1:C]...), part[C+1:]...), join)
			continue
		} else {
			break
		}
	}

	return CreateOutputFromPartitioning(part)
}

// This function calculates output from a given partition, that is the number
// of planes that were found and the normal vectors for those planes
func CreateOutputFromPartitioning(partitioning Partitioning) Output {
	planes := []geometry.Vector{}

	// for every partition we create a plane normal vector
	for _, partition := range partitioning {
		// the plane regression is done like described here:
		// https://stackoverflow.com/questions/20699821/find-and-draw-regression-plane-to-a-set-of-points
		// where we want to find the x for Ax=B

		matAList := []float64{}
		matBList := []float64{}
		for _, point := range partition {
			matAList = append(matAList, point.X, point.Y, 1)
			matBList = append(matBList, point.Z)
		}
		size := len(partition)
		mat1 := mat.NewDense(size, 3, matAList)
		mat2 := mat.NewVecDense(size, matBList)

		res1 := mat.NewDense(3, 3, nil)
		res1.Product(mat1.T(), mat1)

		res2 := mat.NewDense(3, 3, nil)
		res2.Inverse(res1)

		res3 := mat.NewDense(3, size, nil)
		res3.Product(res2, mat1.T())

		res4 := mat.NewVecDense(3, nil)
		res4.MulVec(res3, mat2)

		planes = append(planes, geometry.Vector{X: res4.At(0, 0), Y: res4.At(1, 0), Z: res4.At(2, 0)})
	}

	return Output{NumOfPlanes: len(planes), Planes: planes}
}
