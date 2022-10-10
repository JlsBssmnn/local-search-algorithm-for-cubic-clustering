export function createRandomColors(n: number): string[] {
  const colors = new Array(n);
  for (let i = 0; i < n; i++) {
    colors[
      i
    ] = `rgb(${randomHexValue()}, ${randomHexValue()}, ${randomHexValue()})`;
  }
  return colors;
}

const whiteThreshold = 210;

export function createRandomColorsWithoutWhite(n: number): string[] {
  const colors = new Array(n);
  for (let i = 0; i < n; i++) {
    const c1 = randomHexValue();
    const c2 = randomHexValue();
    let c3;

    if (c1 > whiteThreshold && c2 > whiteThreshold) {
      c3 = Math.round(Math.random() * whiteThreshold);
    } else {
      c3 = randomHexValue();
    }

    colors[i] = `rgb(${c1}, ${c2}, ${c3})`;
  }
  return colors;
}

function randomHexValue(): number {
  return Math.round(Math.random() * 255);
}
