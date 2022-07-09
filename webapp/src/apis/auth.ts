import { ApiUser } from './models/user';
import { ApiError } from './models/api-error';
import { AxiosError } from 'axios';
import { http, postJson } from '../lib/axios';

export async function signin(id: string): Promise<void> {
  try {
    await postJson('/auth/signin', { id: id });
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

export async function confirmSession(): Promise<ApiUser> {
  try {
    return await http.get('/auth/signin');
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
