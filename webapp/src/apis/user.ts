import { ApiError } from './models/api-error';
import { AxiosError } from 'axios';
import { http, putJson } from '../libs/axios';
import { ApiUser, UserLab } from './models/user';

export async function fetchUser(): Promise<ApiUser> {
  try {
    const resp = await http.get<ApiUser>('/user');
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

export async function updateUserLab(user: UserLab): Promise<ApiUser> {
  try {
    const resp = await putJson<ApiUser>('/user/lab', user);
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
