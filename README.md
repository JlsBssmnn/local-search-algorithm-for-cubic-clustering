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