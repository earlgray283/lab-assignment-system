import { ApiError } from './models/api-error';
import { AxiosError } from 'axios';
import { http, postJson } from '../lib/axios';
import { ApiUser } from './models/user';

export async function signin(uid: string): Promise<ApiUser> {
  try {
    const resp = await postJson<ApiUser>('/auth/signin', { uid: uid });
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

export async function signout(): Promise<void> {
  try {
    await http.post('/auth/signout');
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
