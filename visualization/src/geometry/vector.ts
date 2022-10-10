export interface Vector {
  X: number;
  Y: number;
  Z: number;
}

export function dot(v1: Vector, v2: Vector): number {
  return v1.X * v2.X + v1.Y * v2.Y + v1.Z * v2.Z;
}

export function vecLen(v: Vector): number {
  return Math.sqrt(v.X ** 2 + v.Y ** 2 + v.Z ** 2);
}

export function normalize(v: Vector): Vector {
  const len = vecLen(v);
  return {
    X: v.X / len,
    Y: v.Y / len,
    Z: v.Z / len,
  };
}

export function scalarMultiplication(v: Vector, s: number): Vector {
  return {
    X: v.X * s,
    Y: v.Y * s,
    Z: v.Z * s,
  };
}

export function vecAdd(v1: Vector, v2: Vector): Vector {
  return {
    X: v1.X + v2.X,
    Y: v1.Y + v2.Y,
    Z: v1.Z + v2.Z,
  };
}

export function vecDiff(v1: Vector, v2: Vector): Vector {
  return {
    X: v1.X - v2.X,
    Y: v1.Y - v2.Y,
    Z: v1.Z - v2.Z,
  };
}

export interface Vector2D {
  X: number;
  Y: number;
}

export function dot2D(v1: Vector2D, v2: Vector2D): number {
  return v1.X * v2.X + v1.Y * v2.Y;
}

export function vecLen2D(v: Vector2D): number {
  return Math.sqrt(v.X ** 2 + v.Y ** 2);
}

export function normalize2D(v: Vector2D): Vector2D {
  const len = vecLen2D(v);
  return {
    X: v.X / len,
    Y: v.Y / len,
  };
}

export function scalarMultiplication2D(v: Vector2D, s: number): Vector2D {
  return {
    X: v.X * s,
    Y: v.Y * s,
  };
}

export function vecAdd2D(v1: Vector2D, v2: Vector2D): Vector2D {
  return {
    X: v1.X + v2.X,
    Y: v1.Y + v2.Y,
  };
}

export function vecDiff2D(v1: Vector2D, v2: Vector2D): Vector2D {
  return {
    X: v1.X - v2.X,
    Y: v1.Y - v2.Y,
  };
}

export function angle2D(v1: Vector2D, v2: Vector2D): number {
  let angle = Math.atan2(v1.Y, v1.X) - Math.atan2(v2.Y, v2.X);

  if (angle > Math.PI) {
    angle -= 2 * Math.PI;
  } else if (angle < -Math.PI) {
    angle += 2 * Math.PI;
  }

  return angle;
}
