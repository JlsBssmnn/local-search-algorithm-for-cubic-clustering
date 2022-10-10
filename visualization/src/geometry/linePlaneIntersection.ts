import { dot, scalarMultiplication, vecAdd, type Vector } from "./vector";

export type Intersection = "all" | "none" | Vector;

/**
 * Computes the intersection between a line segement and a plane.
 * @param planeNormal The normal vector of the plane
 * @param d The d in the plane equation n*x=d
 * @param segStart One of the 2 starting point of the line segement
 * @param segDir The direction vector of the line segment. `segStart + segDir` has to be the other endpoint of the
 * line segment
 * @returns If there is one point of intersection it returns this point, otherwise `all` if every point of the
 * line segement is in the plane or `none` if the is no intersection.
 */
export default function linesegPlaneIntersec(
  planeNormal: Vector,
  d: number,
  segStart: Vector,
  segDir: Vector
): Intersection {
  const numerator = d - dot(planeNormal, segStart);
  const denominator = dot(planeNormal, segDir);

  if (denominator === 0) {
    if (numerator === 0) {
      return "all";
    } else {
      return "none";
    }
  }

  const t = numerator / denominator;
  if (t < 0 || t > 1) {
    return "none";
  }
  return vecAdd(segStart, scalarMultiplication(segDir, t));
}
