# Local search algorithms for cubic clustering
This repository implements different local search algorithms for the cubic clustering problem. It focuses on the special case of clustering, where the underlying graph is fully connected, which is also called partitioning.
## Cubic Clustering/Partitioning
In cubic partitioning triples of elements are considered when calculating the cost for a partitioning instead of 2 elements. This makes sense for some applications.
For example consider some planes that contain the origin and a set of points in the 3D vector space where each point is sampled from one of the planes. When calculating costs for a triple of points, a plane that best fits the points and contains the origin is generated. Then the maximum distance from one of the 3 points to this plane is used to calculate the cost. Partitioning the points such that all points in one partition are sampled from the same plane only makes sense when considering 3 points at a time (as for 2 points there always exists a plane through the origin that contains both points).
## Usage
### General usage
The partitioning algorithms theoretically work on any problem, but only the above described problem is really implemented. You have to extend the code in order to apply it to your own problem.

### For points sampled from planes
For the partitioning problem of points sampled from planes you can input the data to `main.go` by providing a path to a csv that contains the input data. The program then outputs a partitioning on the standard output.

The csv must have to following structur:
- x, y and z coordinate of a point are described in one row
- the head of the csv may describe which columns contain x, y and z coordinate by giving the head cells the names 'x', 'y', and 'z'
- if the x, y and z columns are not specified it is assumed that the first column contains the x-values, the second the y-values and the third the z-values

Additionally you can specifiy `threshold` and `amplification`. For the cost calculation the maximum distance $d_{max}$ from one of the 3 points to the best fitting plane for the 3 points that contains the origin is calculated. The costs $c$ are then calculated with the threshold $t$ and the amplification $a$ as follows:
$$c = a \cdot (d_{max} - t)$$

### The `helper_scripts` directory
In the `/helper_scripts` directory are python scripts for parsing and visualizing output from the actual application. See the documentation in these files for further use.

### Benchmarks
Benchmarks for the algorithms for the case of points sampled from planes can be found in the file `src/partitioning3D/evaluation/AlgorithmBenchmarks_test.go`. An algorithm can be benchmarked either with or without noise. The data for the partitioning will be created automatically and the number of datapoints depends on the `testTable` which is located in the same file.

To run a benchmark with `benchmarkName` use the following command (with the braces):
`go test ./src/partitioning3D/evaluation -run=^$ -bench=^{benchmarkName}$ -v`

### Evaluation of Algorithms
The algorithms can also be evaluated according to their accuracy. This can be done via the file `src/partitioning3D/evaluation/Algorithm_evaluation_test.go`. You can define the following parameters for the test via command-line arguments:
- `threshold`: The threshold for the cost calculation
- `amplification`: The amplification for the cost calculation
- `mean`: The mean for the noise
- `stddev`: The standard deviation for the noise
- `numberOfPlanes`: How many planes should be used to sample data points
- `pointsPerPlane`: How many points per plane should be sampled
Again, these tests can be done with or without noise, depending on which function is tested. To evaluate an algorithm with the function `evalFunction` that uses a threshold of 0.5 and samples 10 points from 7 different planes respectively use the following command:
`go test ./src/partitioning3D/evaluation -run=^{evalFunction}$ -v -threshold 0.5 -numberOfPlanes 7 -pointsPerPlane 10`

### Compare Algorithms
The `src/partitioning3D/evaluation/Compare_implementations_test.go` file can be used to compare if 2 algorithms work the same way. To do this use the flags `-algorithm1` and `-algorithm2` to specify which 2 algorithms should be compared. Additionally you can specify the following parameters:
-	`iterations`: How many iterations should be executed to test algorithms for equality, each iteration new test data is created and the 2 algorithms are applied to that data
- `seed`: The seed for the random number generation, this can be used for reproduction. If the seed is not specified, **it'll always be 5 which might be unwanted for testing**
- `randomizeParameters`: If `true` the parameters threshold, numOfPlanes, pointsPerPlane, stddev, mean will be randomly choosen in each iteration, otherwise the parameters will be according to the command-line arguments
- `algorithm1`: The first algorithm in the equality test
- `algorithm2`: The second algorithm in the equality test
- `verbose`: If this is 1 the result of every iteration will be printed out, if it's 2 it will also print which elements are differenlty partitioned by the 2 algorithms
All arguments mentioned for the evaluation of algorithms can also be specified but will only be applied if `randomizeParameters` is `false`. 

### Generating test data
The `src/partitioning3D/evaluation/TestDataGenerator_test.go` file can be used to write test data into a csv file. More specifically, the test `TestSaveTestDataToFile` will do this. Again the following parameters can be specified as command-line arguments to describe how the test data will be created:
- `numberOfPlanes`
- `pointsPerPlane`
- `mean`
- `stddev`
Additionally you have to specify where the csv file should be written. This is done via the `outputFile` argument. The value of this arguments must be the entire path to the output file and the file itself. To generate test data from 3 planes with 5 points per plane without any noise use this command:
`go test .\src\partitioning3D\evaluation\ -run=^TestSaveTestDataToFile$ -numberOfPlanes 3 -pointsPerPlane 5 -outputFile /path/to/file/testData.csv`
