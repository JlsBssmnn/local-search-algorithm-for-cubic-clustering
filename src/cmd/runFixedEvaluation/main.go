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

type Result struct {
	GitCommit       string
	Algorithm       string
	Iterations      int
	StddevValues    []float64
	PointsPerPlane  int
	Seed            int64
	AccuracyResults []AccuracyResult
}

type AccuracyResult struct {
	Accuracies []float64
	Time       int64
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
	flag.Parse()

	seed := time.Now().Unix()
	rand.Seed(seed)

	if _, err := os.Stat(*output); !errors.Is(err, os.ErrNotExist) {
		var replace string
		fmt.Printf("The output file seems to exist already, do you want to overwrite it [y/n]?: ")
		fmt.Scanf("%s", &replace)
		if !(strings.ToLower(replace) == "y") && !(strings.ToLower(replace) == "yes") {
			os.Exit(0)
		}
	}
	file, err := os.Create(*output)
	if err != nil {
		panic("Could not open the specified output file")
	}
	defer file.Close()

	algorithm := algorithm.AlgorithmStringToFunc[geometry.Vector](*selectedAlgorithm)

	planes := []geometry.Vector{{X: 1, Y: 0, Z: 0}, {X: 0, Y: 1, Z: 0}, {X: 0, Y: 0, Z: 1}}

	result := Result{
		GitCommit:       gitCommit,
		Algorithm:       *selectedAlgorithm,
		Iterations:      ITERATIONS,
		StddevValues:    STDDEV_VALUES,
		PointsPerPlane:  POINTS_PER_PLANE,
		Seed:            seed,
		AccuracyResults: make([]AccuracyResult, len(STDDEV_VALUES)),
	}
	defer printErrors(result)

	var bar *mpb.Progress
	var mainPb *mpb.Bar
	if *verbose >= 1 {
		bar = mpb.New()
		mainPb = bar.AddBar(int64(len(STDDEV_VALUES)), mpb.PrependDecorators(
			decor.Name("Stddev values:", decor.WCSyncSpaceR),
			decor.CountersNoUnit("%d / %d", decor.WCSyncSpaceR),
		),
			mpb.AppendDecorators(
				decor.AverageETA(decor.ET_STYLE_GO),
				decor.Name(" - "),
				decor.Percentage(decor.WC{W: 5}),
			))
	}

	for i, stddev := range STDDEV_VALUES {
		accuracies := make([]float64, 0, ITERATIONS)
		start := time.Now()
		var secondaryPb *mpb.Bar
		if *verbose >= 1 {
			secondaryPb = bar.AddBar(int64(ITERATIONS), mpb.PrependDecorators(
				decor.Name("iterations:", decor.WCSyncSpaceR),
				decor.CountersNoUnit("%d / %d", decor.WCSyncSpaceR),
			),
				mpb.AppendDecorators(
					decor.AverageETA(decor.ET_STYLE_GO),
					decor.Name(" - "),
					decor.Percentage(decor.WC{W: 5}),
				))
		}
		for j := 0; j < ITERATIONS; j++ {
			testData := evaluation.GenerateDataFromPlanesWithNoise(planes, POINTS_PER_PLANE, utils.NormalDist{Mean: 0, Stddev: stddev})
			calc := partitioning3D.CostCalculator{Threshold: 3 * stddev, Amplification: 3 / stddev}
			eval := evaluation.EvaluateAlgorithm(algorithm, &calc, &testData)
			accuracies = append(accuracies, eval.Accuracy)
			if secondaryPb != nil {
				secondaryPb.Increment()
			}
		}
		elapsed := time.Since(start)
		result.AccuracyResults[i] = AccuracyResult{Accuracies: accuracies, Time: elapsed.Milliseconds()}
		if mainPb != nil {
			mainPb.Increment()
		}
	}

	jsonString, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("Cannot transform go struct to json for the following go struct: %#v", result))
	}
	ioutil.WriteFile(*output, jsonString, os.ModePerm)
}

func printErrors(result Result) {
	err := recover()
	if err != nil {
		fmt.Printf("Error received for the following evaluation: %#v\n", result)
		fmt.Println("Error:", err)
	}
}
