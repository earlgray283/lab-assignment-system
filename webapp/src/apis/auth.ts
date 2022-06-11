import { ApiError } from './models/api-error';
import { AxiosError } from 'axios';
import { postJson } from '../lib/axios';
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
  }
}
