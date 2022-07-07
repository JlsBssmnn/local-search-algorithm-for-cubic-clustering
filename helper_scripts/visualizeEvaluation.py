'''
This file visualizes evaluation that were created via the /src/cmd/runFixedEvaluation/main.go file.
'''

import argparse
import heapq
import itertools
import json
import os.path as path
from matplotlib import pyplot as plt
from matplotlib import style

import numpy as np

if 'myStyle' in plt.style.available:
  style.use('myStyle')

def visualizeResult(path: str, *otherFiles: str):
  '''
  This function visualizes results as box plots. It takes the path to the json file which contains the results.
  Optionally, other files can be specified. The function will try to merge the optional files in the results
  of the main file. 
  '''
  data, stddevValues, algorithm, pointsPerPlane, iterations = mergeResults(path, *otherFiles)

  # get the right scaling on x-axis
  differences = np.array(stddevValues[1:]) - np.array(stddevValues[:-1])
  m = np.min(differences)
  differences = differences / m
  differences = np.insert(differences, 0, 1)
  positions = np.cumsum(differences)

  _, ax = plt.subplots()
  ax.boxplot(data, labels=stddevValues, positions=positions)
  ax.set_xlabel('stddev')
  ax.set_ylabel('accuracy')
  
  plt.suptitle(algorithm, fontsize=30)
  title = f'{f"Datapoints: {pointsPerPlane*3}, " if pointsPerPlane else ""}{f"Iterations: {iterations}" if iterations else ""}'
  plt.title(title, fontsize=20)

  plt.show()

def mergeResults(mainFile, *otherFiles):
  '''
  Takes a list of paths to files and tries to merge their results together. If the results came from
  different algorithms, this function will raise an exception. It returns 5 values.

  Returns
  -------
  data: An array which contains the merged data. The first dimension represents the stddev values,
  the second dimension represents the iterations

  stddevValues: An 1D array which contains the stddev values that correspond to the data array

  algorithm: A string which is the algorithm that was used to produce all of the results
  
  pointsPerPlane: The pointsPerPlane parameter if it's the same for all results, otherwise it'll be None
  
  iterations: The iterations parameter if it's the same for all results, otherwise it'll be None
  '''
  jsonDataList = []
  for file in [mainFile, *otherFiles]:
    f = open(file)
    jsonDataList.append(json.load(f))
    f.close()

  algorithm = set([jsonData['Algorithm'] for jsonData in jsonDataList])
  if len(algorithm) > 1:
    raise Exception('Tried to merge results of different algorithms')
  else:
    algorithm = algorithm.pop()
  pointsPerPlane = set([jsonData['PointsPerPlane'] for jsonData in jsonDataList])
  pointsPerPlane = pointsPerPlane.pop() if len(pointsPerPlane) == 1 else None
  iterations = set([jsonData['Iterations'] for jsonData in jsonDataList])
  iterations = iterations.pop() if len(iterations) == 1 else None

  stddevOrder = []
  for i, jsonData in enumerate(jsonDataList):
    for j, stddevValue in enumerate(jsonData['StddevValues']):
      if not any(map(lambda x: x[0] == stddevValue, stddevOrder)):
        heapq.heappush(stddevOrder, (stddevValue, i, j))

  assert all(map(lambda x: jsonDataList[x[0][1]]['Iterations'] == jsonDataList[x[1][1]]['Iterations'], itertools.combinations(stddevOrder, 2)))

  data = np.empty((len(stddevOrder), jsonDataList[0]['Iterations']))
  i = 0
  stddevValues = np.empty(len(stddevOrder))
  while stddevOrder:
    tup = heapq.heappop(stddevOrder)
    for j, acc in enumerate(jsonDataList[tup[1]]['AccuracyResults'][tup[2]]['Accuracies']):
      data[i, j] = acc
    stddevValues[i] = tup[0]
    i += 1

  data = data.T
  return data, stddevValues, algorithm, pointsPerPlane, iterations


def visualizeTime(*files: str):
  '''
  This function visualizes the execution times of the given results in a plot. It will also show a
  legend which contains some meta information about the results.
  '''
  jsonData = []
  times = []
  for file in files:
    f = open(file)
    jsonData.append(json.load(f))
    f.close()

    timeArray = np.empty(len(jsonData[-1]['AccuracyResults']))
    for i, ar in enumerate(jsonData[-1]['AccuracyResults']):
      timeArray[i] = ar['Time']
    timeArray = timeArray / jsonData[-1]['Iterations'] / 1000
    times.append(timeArray)

  _, ax = plt.subplots()

  for i, timeArray in enumerate(times):
    title = jsonData[i]['Algorithm'] + ': '
    title += f'Datapoints: {jsonData[i]["PointsPerPlane"]*3}, Iterations: {jsonData[i]["Iterations"]}'
    if jsonData[i]['Cores'] != None:
      title += ', CPU-Cores: ' + str(jsonData[i]['Cores'])
    ax.plot(jsonData[i]['StddevValues'], timeArray, label=title)

  ax.set_xlabel('stddev')
  ax.set_ylabel('Execution Time in s')
  ax.legend(loc='best', fontsize=15)

  plt.suptitle('Execution Time Comparison', fontsize=20)

  plt.show()


def visualizeAverages(*files: str):
  '''
  This function visualizes the average accuracies of the given results in a plot. It will also show
  some meta information about the results in a legend.
  '''
  jsonData = []
  times = []
  for file in files:
    f = open(file)
    jsonData.append(json.load(f))
    f.close()

    accArray = np.empty(len(jsonData[-1]['AccuracyResults']))
    for i, ar in enumerate(jsonData[-1]['AccuracyResults']):
      accArray[i] = np.mean(ar['Accuracies'])
    times.append(accArray)

  fig, ax = plt.subplots()

  for i, accArray in enumerate(times):
    title = jsonData[i]['Algorithm'] + ': '
    title += f'Datapoints: {jsonData[i]["PointsPerPlane"]*3}, Iterations: {jsonData[i]["Iterations"]}'
    ax.plot(jsonData[i]['StddevValues'], accArray, label=title)

  ax.set_xlabel('stddev')
  ax.set_ylabel('Average accuracy')
  ax.legend(loc='upper right', fontsize=15)

  plt.suptitle('Accuracy Comparison', fontsize=20)

  plt.show()

if __name__ == '__main__':
  parser = argparse.ArgumentParser(description='Utility for visualizing result files')
  parser.add_argument('fileName', type=str, help='The name of the result file')
  parser.add_argument('-d', '--directory', help='The directory of the result file')
  parser.add_argument('-t', '--type', default='box', help='What should be visualized: box - The results as box plot, time - The execution time')
  parser.add_argument('-f', '--files', default=[], nargs='*', help='Further files that can be visualized for execution time')
  args = parser.parse_args()

  directory = args.directory if args.directory != None and args.directory != '' else './../temp/results'
  fullPath = path.abspath(path.join(path.dirname(__file__), directory, args.fileName))

  furtherFiles = []
  for file in args.files:
    furtherFiles.append(path.abspath(path.join(path.dirname(__file__), directory, file)))
  
  if args.type.lower() == 'box':
    visualizeResult(fullPath, *furtherFiles)
  elif args.type.lower() == 'time':
    visualizeTime(fullPath, *furtherFiles)
  elif args.type.lower() == 'accuracy':
    visualizeAverages(fullPath, *furtherFiles)
