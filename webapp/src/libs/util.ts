export function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

const EPS = 1e-6;

// check a <= b
export function cmpLessThan(a: number, b: number): boolean {
  return a - b <= EPS;
}
