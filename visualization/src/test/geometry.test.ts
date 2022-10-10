import linesegPlaneIntersec from "../geometry/linePlaneIntersection";
import {
  angle2D,
  normalize,
  vecLen,
  type Vector,
  type Vector2D,
} from "../geometry/vector";

test("basic vector operations", () => {
  const v1: Vector = {
    X: 1,
    Y: 0,
    Z: 0,
  };
  const v2: Vector = {
    X: 2,
    Y: 2,
    Z: 1,
  };
  const v3: Vector = {
    X: 0,
    Y: -3,
    Z: -4,
  };

  expect(vecLen(v1)).toBe(1);
  expect(vecLen(v2)).toBe(3);
  expect(vecLen(v3)).toBe(5);

  const normalized = normalize(v3);
  expect(normalized.X).toBe(0);
  expect(normalized.Y).toBe(-0.6);
  expect(normalized.Z).toBe(-0.8);

  expect(vecLen(normalize(v2))).toBe(1);
  expect(vecLen(normalize(v3))).toBe(1);
});

test("the intersection between a line segment and a plane", () => {
  const plane1: Vector = {
    X: 0,
    Y: 1,
    Z: 0,
  };
  const lineStart1: Vector = {
    X: 0,
    Y: 0,
    Z: 0,
  };
  const lineDir1: Vector = {
    X: 1,
    Y: 0,
    Z: 0,
  };
  expect(linesegPlaneIntersec(plane1, 0, lineStart1, lineDir1)).toBe("all");

  const lineStart2: Vector = {
    X: 0,
    Y: 1,
    Z: 0,
  };
  const lineDir2: Vector = {
    X: 1,
    Y: 0,
    Z: 0,
  };
  expect(linesegPlaneIntersec(plane1, 0, lineStart2, lineDir2)).toBe("none");

  const plane2: Vector = {
    X: -1,
    Y: 0,
    Z: 2,
  };
  expect(linesegPlaneIntersec(plane2, 0, lineStart1, lineDir1)).toEqual({
    X: 0,
    Y: 0,
    Z: 0,
  });

  const lineStart3: Vector = {
    X: 5,
    Y: -4,
    Z: 10,
  };
  const lineDir3: Vector = {
    X: -2,
    Y: 8,
    Z: -5,
  };
  expect(linesegPlaneIntersec(plane1, 0, lineStart3, lineDir3)).toEqual({
    X: 4,
    Y: 0,
    Z: 7.5,
  });

  const lineDir4: Vector = {
    X: -2,
    Y: 3,
    Z: -5,
  };
  expect(linesegPlaneIntersec(plane1, 0, lineStart3, lineDir4)).toBe("none");
});

test("2D vector angle", () => {
  const line1: Vector2D = { X: 1, Y: 0 };
  const line2: Vector2D = { X: 0, Y: 1 };
  const line3: Vector2D = { X: 1, Y: 1 };
  const line4: Vector2D = { X: -1, Y: 1 };
  const line5: Vector2D = { X: -1, Y: 0 };
  const line6: Vector2D = { X: 0, Y: -1 };

  expect(angle2D(line1, line1)).toBe(0);
  expect(angle2D(line1, line2)).toBe(-Math.PI / 2);
  expect(angle2D(line1, line3)).toBe(-Math.PI / 4);
  expect(angle2D(line1, line4)).toBe((-3 * Math.PI) / 4);
  expect(angle2D(line1, line5)).toBe(-Math.PI);
  expect(angle2D(line1, line6)).toBe(Math.PI / 2);

  expect(angle2D(line3, line1)).toBe(Math.PI / 4);
  expect(angle2D(line3, line2)).toBe(-Math.PI / 4);
  expect(angle2D(line3, line3)).toBe(0);
  expect(angle2D(line3, line4)).toBe(-Math.PI / 2);
  expect(angle2D(line3, line5)).toBe((-3 * Math.PI) / 4);
  expect(angle2D(line3, line6)).toBe((3 * Math.PI) / 4);

  expect(angle2D(line4, line1)).toBe((3 * Math.PI) / 4);
  expect(angle2D(line4, line2)).toBe(Math.PI / 4);
  expect(angle2D(line4, line3)).toBe(Math.PI / 2);
  expect(angle2D(line4, line4)).toBe(0);
  expect(angle2D(line4, line5)).toBe(-Math.PI / 4);
  expect(angle2D(line4, line6)).toBe((-3 * Math.PI) / 4);

  expect(angle2D(line6, line1)).toBe(-Math.PI / 2);
  expect(angle2D(line6, line2)).toBe(-Math.PI);
  expect(angle2D(line6, line3)).toBe((-3 * Math.PI) / 4);
  expect(angle2D(line6, line4)).toBe((3 * Math.PI) / 4);
  expect(angle2D(line6, line5)).toBe(Math.PI / 2);
  expect(angle2D(line6, line6)).toBe(0);
});
