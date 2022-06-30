package partitioning3D

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
)

// Parses the 3D points which are stored in the given file and returns the data
// as a pointer to a geometry.Vector slice. If the parsing fails, an error is returned.
func ParsePoints(path string) (*[]geometry.Vector, error) {
	err := verifyPath(path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		return nil, errors.New("The specified file wasn't found!")
	}
	return parseCsv(file)
}

func verifyPath(path string) error {
	if path == "" {
		return errors.New("The path to a file with geometry data that should be partitioned must be provided as argument")
	} else if !strings.HasSuffix(path, ".csv") {
		return errors.New("Input file must be a csv file!")
	}
	return nil
}

func parseCsv(file *os.File) (*[]geometry.Vector, error) {
	reader := csv.NewReader(file)

	row, err := reader.Read()
	if err == io.EOF {
		return nil, errors.New("The file is empty")
	} else if err != nil {
		return nil, err
	} else if len(row) < 3 {
		return nil, errors.New("The csv file must contain at least 3 columns for x, y and z coordinates!")
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
		if vector, err := rowToVector(row[xIdx], row[yIdx], row[zIdx]); err != nil {
			return nil, err
		} else {
			data = append(data, *vector)
		}
	} else if xIdx == -1 || yIdx == -1 || zIdx == -1 {
		// some x, y or z is specified but not all so panic
		return nil, errors.New("Csv head doesn't specify each of the x, y and z coordinates!")
	}

	for {
		row, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if vector, err := rowToVector(row[xIdx], row[yIdx], row[zIdx]); err != nil {
			return nil, err
		} else {
			data = append(data, *vector)
		}
	}
	return &data, nil
}

// Try to convert 3 string into a vector by converting each
// into a float. It this fails the function returns an error.
func rowToVector(v1, v2, v3 string) (*geometry.Vector, error) {
	x, errX := strconv.ParseFloat(v1, 64)
	y, errY := strconv.ParseFloat(v2, 64)
	z, errZ := strconv.ParseFloat(v3, 64)

	if errX != nil || errY != nil || errZ != nil {
		return nil, errors.New("Couldn't convert a cell in the csv into a float")
	}

	return &geometry.Vector{X: x, Y: y, Z: z}, nil
}
