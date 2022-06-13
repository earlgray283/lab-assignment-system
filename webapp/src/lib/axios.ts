import axios, { AxiosResponse } from 'axios';

export const http = axios.create({
  baseURL: `${import.meta.env.VITE_BACKEND_HOST}`,
  withCredentials: true,
  transformResponse: (data) => {
    if (typeof data === 'string') {
      try {
        return JSON.parse(data, (k: string, val: unknown) => {
          if (typeof val === 'string') {
            console.log(val);
            const date = new Date(val);
            if (Number.isNaN(date.getDate())) {
              return val;
            }
            return date;
          }
          return val;
        });
      } catch (e) {
        return data;
      }
    }
    return data;
  },
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
