import { Gpa, GradeRegisterToken } from './models/grade';
import { AxiosError } from 'axios';
import { http } from './../lib/axios';
import { ApiError } from './models/api-error';

export async function fetchGpa(): Promise<number> {
  try {
    const resp = await http.get<Gpa>('/grades/me');
    return resp.data.gpa;
  } catch (e: unknown) {
    if (e instanceof AxiosError) {
      if (e.response) {
        console.log(e.response.data);
        const errorJson = e.response.data as ApiError;
        throw new Error(errorJson.message);
      } else {
        throw new Error(e.message);
      }
    }
    throw new Error(e as string);
  }
}

export async function generateToken(): Promise<GradeRegisterToken> {
  try {
    const resp = await http.post<GradeRegisterToken>('/grades/generate-token');
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
