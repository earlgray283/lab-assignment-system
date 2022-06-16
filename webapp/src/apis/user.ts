import { ApiError } from './models/api-error';
import { AxiosError } from 'axios';
import { http } from '../lib/axios';
import { ApiUser } from './models/user';

export async function fetchUser(): Promise<ApiUser> {
  try {
    const resp = await http.get<ApiUser>('/users');
    return resp.data;
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
