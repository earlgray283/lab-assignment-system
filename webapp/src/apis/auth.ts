import { ApiError } from './models/api-error';
import { AxiosError } from 'axios';
import { http, postJson } from '../lib/axios';
import { SignupData } from './models/signup';

export async function signup(data: SignupData): Promise<void> {
  try {
    await postJson('/auth/signup', data);
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

export async function signin(idToken: string): Promise<void> {
  try {
    await postJson('/auth/signin', { idToken: idToken });
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
