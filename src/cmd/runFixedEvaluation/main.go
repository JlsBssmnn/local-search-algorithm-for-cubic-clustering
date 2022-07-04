package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/algorithm"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/geometry"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/partitioning3D"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/partitioning3D/evaluation"
	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

const ITERATIONS = 45

var STDDEV_VALUES = []float64{0.005, 0.01, 0.015, 0.02, 0.025, 0.03, 0.035, 0.04, 0.045, 0.05, 0.055, 0.06, 0.065, 0.07, 0.075, 0.08, 0.085, 0.09, 0.095, 0.1}

const POINTS_PER_PLANE = 33

var bar = mpb.New()

type Result struct {
	GitCommit       string
	Algorithm       string
	Iterations      int
	StddevValues    []float64
	PointsPerPlane  int
	Seed            int64
	AccuracyResults []AccuracyResult
	outputFile      *os.File
}

type AccuracyResult struct {
	Accuracies []float64
	Time       int64
}

func (result *Result) write() {
	jsonString, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("Cannot transform go struct to json for the following go struct: %#v", result))
	}
	ioutil.WriteFile(result.outputFile.Name(), jsonString, os.ModePerm)
}

// This function checks whether the parameters stored in the result struct
// coincide with the specified parameters in this file. If not the function panics.
// This can be used if the evaluation is continued from an existing evaluation.
func (result *Result) checkParameters(algorithm, commit string) {
	var wrongParameter string
	switch {
	case result.Iterations != ITERATIONS:
		wrongParameter = "Iterations"
	case result.Algorithm != algorithm:
		wrongParameter = "Algorithm"
	case !utils.EqualSlices(result.StddevValues, STDDEV_VALUES):
		wrongParameter = "StddevValues"
	case result.PointsPerPlane != POINTS_PER_PLANE:
		wrongParameter = "PointsPerPlane"
	}

	if wrongParameter != "" {
		panic(fmt.Sprintf("%s in existing result file don't coincide with specified %s in main.go", wrongParameter, wrongParameter))
	}
}

func main() {
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	var gitCommit string
	if err != nil {
		gitCommit = string(err.Error())
	} else {
		gitCommit = string(out[:len(out)-1])
	}

	selectedAlgorithm := flag.String("algorithm", "", "The algorithm which should be used for the partitioning")
	output := flag.String("output", "", "Where the output of the evaluation should be written to, the output will be in json")
	verbose := flag.Int("verbose", 0, "0: nothing will be printed, 1: Progress bars will indicate the progress of the evaluation")
	choice := flag.String("choice", "", "What to do if the file already exists, if specified the user will not be requested to give input")
	flag.Parse()

	seed := time.Now().Unix()
	rand.Seed(seed)

	if _, err := os.Stat(*output); *choice == "" && !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("The output file seems to exist already, do you want to overwrite it or continue your work or abort [o/c/a]?: ")
		fmt.Scanf("%s", choice)
		if (strings.ToLower(*choice) == "o") || strings.ToLower(*choice) == "overwrite" {
			*choice = "o"
		} else if (strings.ToLower(*choice) == "c") || strings.ToLower(*choice) == "continue" {
			*choice = "c"
		} else {
			os.Exit(0)
		}
	}

	result := Result{}
	if *choice == "c" {
		// read the existing file and load the results contained in it
		file, _ := os.Open(*output)
		content, _ := ioutil.ReadAll(file)
		json.Unmarshal(content, &result)
		result.checkParameters(*selectedAlgorithm, gitCommit)

		file.Close()
	} else if *choice != "o" {
		os.Exit(0)
	}

	file, err := os.Create(*output)
	if err != nil {
		panic("Could not open the specified output file")
	}
	defer file.Close()

	if *choice == "c" {
		result.outputFile = file
		result.write()
		rand.Seed(result.Seed)
	} else {
		result = Result{
			GitCommit:       gitCommit,
			Algorithm:       *selectedAlgorithm,
			Iterations:      ITERATIONS,
			StddevValues:    STDDEV_VALUES,
			PointsPerPlane:  POINTS_PER_PLANE,
			Seed:            seed,
			AccuracyResults: make([]AccuracyResult, 0, len(STDDEV_VALUES)),
			outputFile:      file,
		}
	}

	algorithm := algorithm.AlgorithmStringToFunc[geometry.Vector](*selectedAlgorithm)
	planes := []geometry.Vector{{X: 1, Y: 0, Z: 0}, {X: 0, Y: 1, Z: 0}, {X: 0, Y: 0, Z: 1}}
	defer printErrors(result)

	// Determine the starting point if execution of evalution is continued
	stddevOffset := len(result.AccuracyResults)
	if stddevOffset > 0 && len(result.AccuracyResults[stddevOffset-1].Accuracies) < ITERATIONS {
		stddevOffset--
	}

	var iterationOffset int
	if len(result.AccuracyResults)-1 == stddevOffset {
		iterationOffset = len(result.AccuracyResults[stddevOffset].Accuracies)
	} else {
		result.AccuracyResults = append(result.AccuracyResults, AccuracyResult{Accuracies: make([]float64, 0, ITERATIONS), Time: 0})
	}

	fixSeed(stddevOffset, iterationOffset, planes)

	mainPb := createPb(int64(len(STDDEV_VALUES)), "Stddev values:", *verbose, stddevOffset)

	for i := stddevOffset; i < len(STDDEV_VALUES); i++ {
		stddev := result.StddevValues[i]

		var startIteration int
		if i == stddevOffset {
			startIteration = iterationOffset
		} else {
			result.AccuracyResults = append(result.AccuracyResults, AccuracyResult{Accuracies: make([]float64, 0, ITERATIONS), Time: 0})
		}

		var secondaryPb = createPb(int64(ITERATIONS), "iterations:", *verbose, startIteration)

		for j := startIteration; j < ITERATIONS; j++ {
			start := time.Now()
			testData := evaluation.GenerateDataFromPlanesWithNoise(planes, POINTS_PER_PLANE, utils.NormalDist{Mean: 0, Stddev: stddev})
			calc := partitioning3D.CostCalculator{Threshold: 3 * stddev, Amplification: 1 / stddev}
			eval := evaluation.EvaluateAlgorithm(algorithm, &calc, &testData)

			result.AccuracyResults[i].Accuracies = append(result.AccuracyResults[i].Accuracies, eval.Accuracy)
			result.AccuracyResults[i].Time += time.Since(start).Milliseconds()
			result.write()
			if secondaryPb != nil {
				secondaryPb.Increment()
			}
		}
		if mainPb != nil {
			mainPb.Increment()
		}
	}
}

func printErrors(result Result) {
	err := recover()
	if err != nil {
		fmt.Printf("Error received for the following evaluation: %#v\n", result)
		fmt.Println("Error:", err)
		panic(err)
	}
}

func createPb(total int64, name string, verbose, startProgress int) *mpb.Bar {
	var mainPb *mpb.Bar
	if verbose >= 1 {
		mainPb = bar.AddBar(total, mpb.PrependDecorators(
			decor.Name(name, decor.WCSyncSpaceR),
			decor.CountersNoUnit("%d / %d", decor.WCSyncSpaceR),
		),
			mpb.AppendDecorators(
				decor.AverageETA(decor.ET_STYLE_GO),
				decor.Name(" - "),
				decor.Percentage(decor.WC{W: 5}),
			))
		mainPb.IncrBy(startProgress)
	}
	return mainPb
}

// This function generates random data so often that the seed which was
// used in a prior execution now matches with this execution
func fixSeed(startStddevValue, iterationStart int, planes []geometry.Vector) {
	for i := 0; i < startStddevValue*ITERATIONS+iterationStart; i++ {
		evaluation.GenerateDataFromPlanesWithNoise(planes, POINTS_PER_PLANE, utils.NormalDist{Mean: 0, Stddev: 0})
	}
}
