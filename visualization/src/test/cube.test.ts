import type { Vector } from "../geometry/vector";
import Cube from "../geometry/cube";

test("cube-plane intersections", () => {
  const cube1 = new Cube({ X: 0, Y: 0, Z: 0 }, 1, 1, 1);
  const cube2 = new Cube({ X: 1, Y: 1, Z: 2 }, 5, 1, 2);
  const cube3 = new Cube({ X: -1, Y: 0, Z: 1 }, -1, 4, -3);

  const planeNormal1: Vector = { X: 1, Y: 0, Z: 0 };
  const planeNormal2: Vector = { X: 1, Y: 0, Z: 1 };
  const planeNormal3: Vector = { X: -1, Y: 1, Z: 2 };

  // cube1

  let intersect = cube1.planeIntersect(planeNormal1, 0);
  expect(intersect).toHaveLength(4);
  expect(intersect).toContainEqual({ X: 0, Y: 0, Z: 0 });
  expect(intersect).toContainEqual({ X: 0, Y: 1, Z: 0 });
  expect(intersect).toContainEqual({ X: 0, Y: 0, Z: 1 });
  expect(intersect).toContainEqual({ X: 0, Y: 1, Z: 1 });

  intersect = cube1.planeIntersect(planeNormal1, 0.5);
  expect(intersect).toHaveLength(4);
  expect(intersect).toContainEqual({ X: 0.5, Y: 0, Z: 0 });
  expect(intersect).toContainEqual({ X: 0.5, Y: 1, Z: 0 });
  expect(intersect).toContainEqual({ X: 0.5, Y: 0, Z: 1 });
  expect(intersect).toContainEqual({ X: 0.5, Y: 1, Z: 1 });

  intersect = cube1.planeIntersect(planeNormal2, 0);
  expect(intersect).toHaveLength(2);
  expect(intersect).toContainEqual({ X: 0, Y: 0, Z: 0 });
  expect(intersect).toContainEqual({ X: 0, Y: 1, Z: 0 });

  intersect = cube1.planeIntersect(planeNormal2, 1);
  expect(intersect).toHaveLength(4);
  expect(intersect).toContainEqual({ X: 0, Y: 0, Z: 1 });
  expect(intersect).toContainEqual({ X: 0, Y: 1, Z: 1 });
  expect(intersect).toContainEqual({ X: 1, Y: 0, Z: 0 });
  expect(intersect).toContainEqual({ X: 1, Y: 1, Z: 0 });

  intersect = cube1.planeIntersect(planeNormal3, 6);
  expect(intersect).toHaveLength(0);

  intersect = cube1.planeIntersect(planeNormal3, 0);
  expect(intersect).toHaveLength(3);
  expect(intersect).toContainEqual({ X: 0, Y: 0, Z: 0 });
  expect(intersect).toContainEqual({ X: 1, Y: 0, Z: 0.5 });
  expect(intersect).toContainEqual({ X: 1, Y: 1, Z: 0 });

  // cube2
  intersect = cube2.planeIntersect(planeNormal1, 0);
  expect(intersect).toHaveLength(0);

  intersect = cube2.planeIntersect(planeNormal1, 1.5);
  expect(intersect).toHaveLength(4);
  expect(intersect).toContainEqual({ X: 1.5, Y: 1, Z: 2 });
  expect(intersect).toContainEqual({ X: 1.5, Y: 1, Z: 3 });
  expect(intersect).toContainEqual({ X: 1.5, Y: 3, Z: 2 });
  expect(intersect).toContainEqual({ X: 1.5, Y: 3, Z: 3 });

  intersect = cube2.planeIntersect(planeNormal2, 3.5);
  expect(intersect).toHaveLength(4);
  expect(intersect).toContainEqual({ X: 1, Y: 1, Z: 2.5 });
  expect(intersect).toContainEqual({ X: 1, Y: 3, Z: 2.5 });
  expect(intersect).toContainEqual({ X: 1.5, Y: 1, Z: 2 });
  expect(intersect).toContainEqual({ X: 1.5, Y: 3, Z: 2 });

  // cube3

  intersect = cube3.planeIntersect(planeNormal1, -1.5);
  expect(intersect).toHaveLength(4);
  expect(intersect).toContainEqual({ X: -1.5, Y: 0, Z: 1 });
  expect(intersect).toContainEqual({ X: -1.5, Y: 0, Z: 5 });
  expect(intersect).toContainEqual({ X: -1.5, Y: -3, Z: 1 });
  expect(intersect).toContainEqual({ X: -1.5, Y: -3, Z: 5 });
});

test("sorted intersections", () => {
  const cube1 = new Cube({ X: 0, Y: 0, Z: 0 }, 1, 1, 1);

  const planeNormal1: Vector = { X: 1, Y: 0, Z: 0 };

  let intersect = cube1.planeIntersectSorted(planeNormal1, 0);
  expect(intersect).toEqual([
    { X: 0, Y: 1, Z: 1 },
    { X: 0, Y: 1, Z: 0 },
    { X: 0, Y: 0, Z: 0 },
    { X: 0, Y: 0, Z: 1 },
  ]);
});
