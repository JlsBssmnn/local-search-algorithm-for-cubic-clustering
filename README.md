# Local search algorithms for cubic clustering
This repository implements different local search algorithms for the cubic clustering problem. It focuses on the special case of clustering, where the underlying graph is fully connected, which is also called partitioning.
## Cubic Clustering/Partitioning
In cubic partitioning triples of elements are considered when calculating the cost for a partitioning instead of 2 elements. This makes sense for some applications.
For example consider some planes that contain the origin and a set of points in the 3D vector space where each point is sampled from one of the planes. When calculating costs for a triple of points, a plane that best fits the points and contains the origin is generated. Then the maximum distance from one of the 3 points to this plane is used to calculate the cost. Partitioning the points such that all points in one partition are sampled from the same plane only makes sense when considering 3 points at a time (as for 2 points there always exists a plane through the origin that contains both points).
## Usage
### General usage
The partitioning algorithms theoretically work on any problem, but only the above described problem is really implemented. You have to extend the code in order to apply it to your own problem.

### For points sampled from planes
For the partitioning problem of points sampled from planes you can input the data to `src/cmd/partitionByCsv/main.go` by providing a path to a csv that contains the input data. The program then outputs a partitioning on the standard output.

The csv must have to following structure:
- x, y and z coordinate of a point are described in one row
- the head of the csv may describe which columns contain x, y and z coordinate by giving the head cells the names 'x', 'y', and 'z'
- if the x, y and z columns are not specified it is assumed that the first column contains the x-values, the second the y-values and the third the z-values

Additionally you can specify `threshold` and `amplification`. For the cost calculation the maximum distance $d_{max}$ from one of the 3 points to the best fitting plane for the 3 points that contains the origin is calculated. The costs $c$ are then calculated with the threshold $t$ and the amplification $a$ as follows:
$$c = a \cdot (d_{max} - t)$$

### Fixed Evaluation
An evaluation which will track accuracy and execution time of an algorithm can be started via the `src/cmd/runFixedEvaluation/main.go` file. This evaluation will use the XY, XZ and YZ planes to sample data point. Threshold and amplification depend on the standard deviation $\sigma$ and are calculated like this:

$$t = 3\cdot \sigma \quad\text{and}\quad a = \frac{1}{\sigma}$$

Thus the cost calculation is basically a line which intersects the y-axis at $-3$ and intersects the x-axis at $3\sigma$.

The evaluation will be written to a json file, the path to this file has to be provided as an argument. To customize some of the parameters for the evaluation you can use a *configuration file*. Examples for these files are in `/src/temp/eval_configs`. Which file should be applied has to be specified as argument too. By default, the `default_config.json` file will be used.

It is also possible to continue the execution of an evaluation. For this just specify the already existing output file again as output file. The program will detect this and will ask whether to overwrite the file, continue execution or to abort. The input has to be given in the shell but can also be passed as a flag, so the program will bot wait for user input.

In summary these are the arguments which can be used:
- `output`: Path to a file where the evaluation output should be written to
- `algorithm`: The algorithm that will be evaluated
- `verbose`: If 0 nothing will be printed, if 1 a progress bar will indicate the progress
- `choice`: Whether to overwrite(o) an existing file, continue(c) the execution or abort(a) if the output file already exists
- `config`: The path to the configuration file

An example:
```sh
go run ./src/cmd/runFixedEvaluation/main.go -algorithm GreedyJoining -output ./temp/results/results.json
```

### The `helper_scripts` directory
In the `/helper_scripts` directory are python scripts for parsing and visualizing output from the actual application. See the documentation in these files for further use.

### Benchmarks
Benchmarks for the algorithms for the case of points sampled from planes can be found in the file `src/partitioning3D/evaluation/AlgorithmBenchmarks_test.go`. The data for the partitioning will be created automatically. The following parameters can be specified via command-line arguments:
- `algorithm1`: The algorithm that should be benchmarked
- `threshold`: The threshold for the cost calculation
- `amplification`: The amplification for the cost calculation
- `mean`: The mean for the noise
- `stddev`: The standard deviation for the noise
- `numberOfPlanes`: How many planes should be used to sample data points
- `pointsPerPlane`: How many points per plane should be sampled

To run a benchmark for the Greedy Joining algorithm use the following command:
```sh
go test ./src/partitioning3D/evaluation -run=^$ -bench=^BenchmarkAlgorithm$ -v -algorithm1 GreedyJoining
```

### Evaluation of Algorithms
The algorithms can also be evaluated according to their accuracy. This can be done via the file `src/partitioning3D/evaluation/Algorithm_evaluation_test.go`. You can specify the same arguments as for the benchmarks.

To evaluate Greedy Joining you could use:
```sh
go test ./src/partitioning3D/evaluation -run=^TestEvalAlgorithm$ -v -algorithm1 GreedyJoining -threshold 0.5 -numberOfPlanes 7 -pointsPerPlane 10
```

### Compare Algorithms
The `src/partitioning3D/evaluation/Compare_implementations_test.go` file can be used to compare if 2 algorithms work the same way. To do this use the flags `-algorithm1` and `-algorithm2` to specify which 2 algorithms should be compared. Additionally you can specify the following parameters:
-	`iterations`: How many iterations should be executed to test algorithms for equality, each iteration new test data is created and the 2 algorithms are applied to that data
- `seed`: The seed for the random number generation, this can be used for reproduction. If the seed is not specified, **it'll always be 5 which might be unwanted for testing**
- `randomizeParameters`: If `true` the parameters threshold, numOfPlanes, pointsPerPlane, stddev, mean will be randomly chosen in each iteration, otherwise the parameters will be according to the command-line arguments
- `algorithm1`: The first algorithm in the equality test
- `algorithm2`: The second algorithm in the equality test
- `verbose`: If this is 1 the result of every iteration will be printed out, if it's 2 it will also print which elements are differently partitioned by the 2 algorithms

All arguments mentioned for the evaluation of algorithms can also be specified but will only be applied if `randomizeParameters` is `false`. 

### Generating test data
The `src/partitioning3D/evaluation/TestDataGenerator_test.go` file can be used to write test data into a csv file. More specifically, the test `TestSaveTestDataToFile` will do this. Again the following parameters can be specified as command-line arguments to describe how the test data will be created:
- `numberOfPlanes`
- `pointsPerPlane`
- `mean`
- `stddev`

Additionally you have to specify where the csv file should be written to. This is done via the `-outputFile` argument. The value of this arguments must be the entire path to the output file and the file itself. To generate test data from 3 planes with 5 points per plane without any noise use this command:
```sh
go test ./src/partitioning3D/evaluation -run=^TestSaveTestDataToFile$ -numberOfPlanes 3 -pointsPerPlane 5 -mean 0 -stddev 0 -outputFile /path/to/file/testData.csv
```

### Calculating costs and save it to JSON
The `src/partitioning3D/evaluation/TestDataGenerator_test.go` file can also be used to write the triple costs to a json file in order to use it in some other context without having to calculate the cost yourself. The test `TestSaveCostToFile` is used for this. To get the costs you have to specify the input as a path to a csv file which contains the input data. You also must specify the destination, so where the output should be written to. Thus you can specify the following arguments for this usage:
- `inputFile`
- `outputFile`
- `threshold`
- `amplification`

An example:
```sh
go test ./src/partitioning3D/evaluation/ -run=^TestSaveCostToFile$ -inputFile /path/to/input.csv -outputFile /path/to/output.json -threshold 0.001 -amplification 3
```
