package partitioning3D

import (
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/algorithm"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"
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
	size := len(points)
	matAList := make([]float64, size*3)

	for i := 0; i < size; i++ {
		matAList[i*3] = points[i].X
		matAList[i*3+1] = points[i].Y
		matAList[i*3+2] = points[i].Z
	}
	A := mat.NewDense(size, 3, matAList)

	var M mat.Dense
	M.Mul(A.T(), A)
	A = nil

	var svd mat.SVD
	ok := svd.Factorize(&M, mat.SVDFull)
	if !ok {
		panic("Failed to factorize")
	}
	singularValues := svd.Values(nil)

	var u mat.Dense
	svd.UTo(&u)

	i := utils.ArgMin(singularValues)
	x := u.At(0, i)
	y := u.At(1, i)
	z := u.At(2, i)

	return geometry.Vector{X: x, Y: y, Z: z}
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
