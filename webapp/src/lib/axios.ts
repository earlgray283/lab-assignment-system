import axios, { AxiosResponse } from 'axios';

export const http = axios.create({
  baseURL: `${import.meta.env.VITE_BACKEND_HOST}`,
});

export async function postJson<
  T = unknown,
  R = AxiosResponse<T, unknown>,
  D = unknown
>(uri: string, data: D): Promise<R> {
  const resp = http.post<T, R, D>(uri, data, {
    headers: {
      'Content-Type': 'application/json',
    },
  });
  return resp;
}

export function isErrorStatus(code: number): boolean {
  return code >= 400;
}
