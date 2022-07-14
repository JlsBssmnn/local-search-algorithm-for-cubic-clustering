'''
This utility offers the possibility to calculate costs for data points. The data points are 3D vectors.
The x-,y- and z-coordinates have to be input comma-separated. Multiple points are also just separated.

Example: x1,y1,z1,x2,y2,z2,x3,y3,z3

Threshold and amplification for the calculation may be provided as arguments.
'''
import numpy as np
import argparse

def getPoints():
  return input('Input the data points (separated by comma) for which the costs should be calculated \
    or type "exit" to exit the utility: ')

if __name__ == '__main__':
  parser = argparse.ArgumentParser(description='Utility for calculating costs between data points, type "exit" to exit this utility')
  parser.add_argument('-t', '--threshold', type=float, default=1, help='The threshold for the cost calculation')
  parser.add_argument('-a', '--amplification', type=float, default=1, help='The amplification for the cost calculation')
  args = parser.parse_args()

  user_input = getPoints()
  while user_input.lower() != 'exit':
    points = user_input.split(',')
    points = [x.replace(' ', '') for x in points]
    try:
      points = [float(x) for x in points]
    except Exception:
      print('The input data points could not be parsed into floats, make sure the input has the correct format!')
      user_input = getPoints()
      continue

    size = len(points)
    if size / 3 != int(size / 3):
      print('You did not input the points correctly, each point must have three number representing X, Y and Z coordinate')
      user_input = getPoints()
      continue

    size = int(size / 3)
    A = np.empty((size, 3))

    for i in range(0, size):
      A[i][0] = points[i*3]
      A[i][1] = points[i*3 + 1]
      A[i][2] = points[i*3 + 2]

    M = np.matmul(A.T, A)
    svd = np.linalg.svd(M)

    i = np.argmin(svd[1])
    normal_vector = svd[0][:, i]

    distances = []
    for i in range(A.shape[0]):
      row = A[i]
      distances.append(abs(normal_vector[0]*row[0]+normal_vector[1]*row[1]+normal_vector[2]*row[2]) / np.linalg.norm(normal_vector))

    print('Input points as matrix:', A, sep='\n')
    print('Costs:', args.amplification * (max(distances) - args.threshold))

    user_input = getPoints()