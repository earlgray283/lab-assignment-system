import { AxiosError } from 'axios';
import { http } from '../libs/axios';
import { ApiError } from './models/api-error';

export async function fetchGrades(year: number): Promise<number[]> {
  try {
    const resp = await http.get<number[]>('/grades', {
      params: {
        year: year,
      },
    });
    return resp.data;
  } catch (e: unknown) {
    if (e instanceof AxiosError) {
      if (e.response) {
        const errorJson = e.response.data as ApiError;
        throw new Error(errorJson.message);
      } else {
        throw new Error(e.message);
      }
    }
    throw new Error(e as string);
  }
}
