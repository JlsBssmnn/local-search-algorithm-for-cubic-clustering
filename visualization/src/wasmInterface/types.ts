import type { Vector } from "src/geometry/vector";

export interface TestData {
  NumOfPlanes: number;
  Planes: Vector[];
  Points: Vector[];
}

export interface Evaluation {
  NumOfPlanesError: number;
  Accuracy: number;
  TotalEdges: number;
  TruePositives: number;
  TrueNegatives: number;
  FalsePositives: number;
  FalseNegatives: number;
  ComputedPlanes: Vector[];
}
