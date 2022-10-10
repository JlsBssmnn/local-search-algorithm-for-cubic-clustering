import linesegPlaneIntersec from "./linePlaneIntersection";
import {
  angle2D,
  dot,
  dot2D,
  normalize2D,
  scalarMultiplication,
  vecAdd,
  vecDiff2D,
  type Vector,
  type Vector2D,
} from "./vector";

export interface LineSegment {
  start: Vector;
  direction: Vector;
}

/**
 * Represents a cube in 3D space. The sides may have different
 * lengths. Every side of the cube will be parallel to one of the
 * x-, y- or z-axis, so it can't be rotated.
 *
 * A cube is defined buy the lower-left-back corner, so the edge with
 * minimum x-, y- and z-coordinates, and the lengths of the sides.
 */
export default class Cube {
  corner: Vector;
  width: Vector;
  height: Vector;
  depth: Vector;

  constructor(corner: Vector, width: number, height: number, depth: number) {
    this.corner = corner;
    this.width = {
      X: width,
      Y: 0,
      Z: 0,
    };
    this.height = {
      X: 0,
      Y: 0,
      Z: height,
    };
    this.depth = {
      X: 0,
      Y: depth,
      Z: 0,
    };
  }

  /**
   * Retruns the sides of the cube as an arry of line segements.
   */
  getSides(): LineSegment[] {
    const sides: LineSegment[] = [];

    // the segments are defined by 4 anchor points, `this.corner` is the first
    const p1 = this.corner;
    const p2 = vecAdd(vecAdd(p1, this.width), this.depth); // lower-right-front
    const p3 = vecAdd(vecAdd(p1, this.width), this.height); // upper-right-back
    const p4 = vecAdd(vecAdd(p1, this.height), this.depth); // upper-left-front

    sides.push({ start: p1, direction: this.width });
    sides.push({ start: p1, direction: this.height });
    sides.push({ start: p1, direction: this.depth });

    sides.push({ start: p2, direction: scalarMultiplication(this.width, -1) });
    sides.push({ start: p2, direction: this.height });
    sides.push({ start: p2, direction: scalarMultiplication(this.depth, -1) });

    sides.push({ start: p3, direction: scalarMultiplication(this.width, -1) });
    sides.push({ start: p3, direction: scalarMultiplication(this.height, -1) });
    sides.push({ start: p3, direction: this.depth });

    sides.push({ start: p4, direction: this.width });
    sides.push({ start: p4, direction: scalarMultiplication(this.height, -1) });
    sides.push({ start: p4, direction: scalarMultiplication(this.depth, -1) });

    return sides;
  }

  /**
   * Computes a list of all intersection points of the cube with the given plane.
   */
  planeIntersect(planeNormal: Vector, d: number): Vector[] {
    const sides = this.getSides();
    const intersects: Set<string> = new Set();

    for (const side of sides) {
      const intersect = linesegPlaneIntersec(
        planeNormal,
        d,
        side.start,
        side.direction
      );
      if (intersect === "none") {
        continue;
      } else if (intersect == "all") {
        intersects.add(JSON.stringify(side.start));
        intersects.add(JSON.stringify(vecAdd(side.start, side.direction)));
      } else {
        intersects.add(JSON.stringify(intersect));
      }

      // There can only be at most 6 intersections
      if (intersects.size >= 6) {
        break;
      }
    }

    return [...intersects].map((intersect) => JSON.parse(intersect));
  }

  /**
   * Returns the intersection points with a plane sorted, such that if
   * you connect consecutive points in the array you will get the shape
   * that is formed by those points.
   */
  planeIntersectSorted(planeNormal: Vector, d: number): Vector[] {
    const intersects = this.planeIntersect(planeNormal, d);

    const axisToKeep: (keyof Vector)[] = []; // The coordinates that we don't want to squash

    if (dot(planeNormal, { X: 1, Y: 0, Z: 0 }) == 0) {
      axisToKeep.push("X");
    }
    if (dot(planeNormal, { X: 0, Y: 1, Z: 0 }) == 0) {
      axisToKeep.push("Y");
    }
    if (dot(planeNormal, { X: 0, Y: 0, Z: 1 }) == 0) {
      axisToKeep.push("Z");
    }

    // by default we use X and Y as the axis to keep
    // only the first 2 elements of this array are relevant
    if (axisToKeep.length < 2) {
      if (!axisToKeep.includes("X")) {
        axisToKeep.push("X");
      }
      axisToKeep.push("Y");
    }

    const projected: Vector2D[] = intersects.map((vector) => ({
      X: vector[axisToKeep[0]],
      Y: vector[axisToKeep[1]],
    })); // we project the points on a coordinate plane

    // find a rectangle that surrounds the projected area
    let minX = Infinity,
      minY = Infinity,
      maxX = -Infinity,
      maxY = -Infinity;
    for (const vec of projected) {
      if (vec.X < minX) {
        minX = vec.X;
      }
      if (vec.Y < minY) {
        minY = vec.Y;
      }
      if (vec.X > maxX) {
        maxX = vec.X;
      }
      if (vec.Y > maxY) {
        maxY = vec.Y;
      }
    }

    // the middle of this rectangle
    const middle: Vector2D = {
      X: minX + (maxX - minX) / 2,
      Y: minY + (maxY - minY) / 2,
    };

    // determine the angles for each point (we use the angle between [the line through the middle point
    // and the current point] and the line that we get for the first point this way)
    const reference = normalize2D(vecDiff2D(projected[0], middle));
    const angles = projected.map((vec) => {
      const line = normalize2D(vecDiff2D(vec, middle));
      return angle2D(reference, line);
    });

    return projected
      .map((vec, i) => ({ vec, i, angle: angles[i] }))
      .sort((e1, e2) => e1.angle - e2.angle)
      .map((el) => intersects[el.i]);
  }
}
