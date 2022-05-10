package partitioning3D

import (
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/algorithm"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
	"gonum.org/v1/gonum/mat"
)

type Output struct {
	NumOfPlanes int
	Planes      []geometry.Vector
}

// This functions takes in an arbitrary number of points and tries
// to find the best plane that goes through the origin that minimizes
// the sum of squared distances from the points to the plane
func FitPlane(points ...geometry.Vector) geometry.Vector {
	matAList := []float64{}
	matBList := []float64{}
	allZero := true

	for _, point := range points {
		matAList = append(matAList, point.X, point.Y)
		matBList = append(matBList, point.Z)
		if point.X != 0 || point.Y != 0 {
			allZero = false
		}
	}

	// if all x and y coordinates of the points are zero the SVD will fail,
	// so we return the y-z plane as it fits all the points
	if allZero {
		return geometry.Vector{X: 1, Y: 0, Z: 0}
	}

	size := len(points)
	A := mat.NewDense(size, 2, matAList)
	b := mat.NewVecDense(size, matBList)

	var res mat.Dense
	if err := res.Solve(A, b); err == nil {
		return geometry.CoordinateReprToNormalVec(res.At(0, 0), res.At(1, 0))
	}

	bMat := mat.NewDense(size, 1, matBList)

	var svd mat.SVD
	ok := svd.Factorize(A, mat.SVDFull)
	if !ok {
		panic("Failed to factorize")
	}

	rank := svd.Rank(1e-15)
	var x mat.Dense
	svd.SolveTo(&x, bMat, rank)

	return geometry.CoordinateReprToNormalVec(x.At(0, 0), x.At(1, 0))
}

// This function calculates output from a given partition, that is the number
// of planes that were found and the normal vectors for those planes
func CreateOutputFromPartitioning(partitioning algorithm.Partitioning[geometry.Vector]) Output {
	planes := []geometry.Vector{}

	// for every partition we create a plane normal vector
	for _, partition := range partitioning {
		planes = append(planes, FitPlane(partition...))
	}

	return Output{NumOfPlanes: len(planes), Planes: planes}
}
