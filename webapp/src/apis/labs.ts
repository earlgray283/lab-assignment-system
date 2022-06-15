import { LabList } from './models/lab';
import { AxiosError } from 'axios';
import { http } from '../lib/axios';
import { ApiError } from './models/api-error';

export async function fetchLabs(labIds?: string[]): Promise<LabList> {
  try {
    const resp = await http.get<LabList>('/labs', {
      params: {
        labIds: labIds && labIds.join('+'),
      },
    });
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
