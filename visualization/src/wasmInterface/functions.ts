import type { Evaluation, TestData } from "./types";

export function generateData(
  numOfPlanes: number,
  pointsPerPlane: number,
  mean: number,
  stddev: number
): TestData {
  // @ts-ignore
  return window.generateData(numOfPlanes, pointsPerPlane, mean, stddev);
}

export function partitionData(
  algorithm: string,
  threshold: number,
  amplification: number
): Evaluation {
  // @ts-ignore
  return window.partitionData(algorithm, threshold, amplification);
}
