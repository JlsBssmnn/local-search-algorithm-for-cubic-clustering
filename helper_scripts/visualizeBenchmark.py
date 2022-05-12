from typing import Dict
import matplotlib.pylab as plt
from os import path
import pandas as pd
import argparse

def visualizeBenchmark(path: str, x: str, y: str, parameters: Dict[str, str] = {}):
  '''
  Visualize a benchmark that is saved in a csv file.

  Parameters
  ----------
  path: The path to the csv file which will be used
  x: The column name of the column of which the values will be used for the x-axis
  y: The column name of the column of which the values will be used for the y-axis
  parameters: A dict where the keys are column names and the values are numbers (as string)
  that occur in the csv. It is required if the relation of values from
  x to y is not a map (one x value is mapped to 2 different y values). In this case only those x-values
  that are in the same row with the specified parameter values will be plotted.
  '''
  df = pd.read_csv(path)
  if (x not in df) or (y not in df):
    raise ValueError('The specified x or y is not a column name')
  xAxis = []
  yAxis = []
  for i in set(df[x]):
    query = f'{x} == {i}'
    for k, v in parameters.items():
      query += f' and {k} == {v} '
    row = df.query(query)
    if len(row.index) > 1:
      raise ValueError('Not enough parameters specified, there are still ambiguous rows where it \
        is not clear which of them should be plotted')
    elif len(row.index) == 0:
      continue
    xAxis.append(i)
    yAxis.append(row[y].iloc[0])

  data = sorted(zip(xAxis, yAxis), key=lambda t: t[0])
  xAxis = [x for x,_ in data]
  yAxis = [y for _,y in data]

  plt.plot(xAxis, yAxis)
  plt.xlabel(x)
  plt.ylabel(y)
  plt.title(str(parameters))
  plt.show()


if __name__ == '__main__':
  parser = argparse.ArgumentParser(description='Utility for visualizing a parsed go benchmark file')
  parser.add_argument('file', type=str, help='The path to the parsed benchmark file')
  parser.add_argument('x', type=str, help='The column name which will be the x-axis')
  parser.add_argument('y', type=str, help='The column name which will be the y-axis')
  parser.add_argument('-p', '--parameter', action='append', type=str, default=[],
    help='Other parameter and it\'s value in the form {parameter: value}. This is required for \
    creating a plot, see documentation for visualizeBenchmark()')
  args = parser.parse_args()

  parameters = {}
  for parameter in args.parameter:
    parameter = parameter.replace(' ', '')
    parameters[parameter[:parameter.rfind(':')]] = parameter[parameter.rfind(':') + 1:]

  visualizeBenchmark(path.realpath(args.file), args.x, args.y, parameters)
