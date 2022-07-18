export function mergeGpas(
  gpas1: number[],
  gpas2: number[],
  gpas3: number[]
): number[] {
  const gpas = Array.prototype.concat(gpas1 ?? [], gpas2 ?? [], gpas3 ?? []);
  gpas.sort((a, b) => b - a);
  return gpas;
}
