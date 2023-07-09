import { LabList } from './models/lab';
import { AxiosError } from 'axios';
import { http } from '../libs/axios';
import { ApiError } from './models/api-error';

export async function fetchLabs(
  year: number,
  labIds?: string[],
  optFields?: string[]
): Promise<LabList> {
  try {
    const resp = await http.get<LabList>('/labs', {
      params: {
        year: year,
        labIds: labIds && labIds.join('+'),
        optFields: optFields && optFields.join('+'),
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
