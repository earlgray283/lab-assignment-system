import { AxiosError } from 'axios';
import { http } from '../libs/axios';
import { ApiError } from './models/api-error';

export async function finalDecisionDryRun(year: number) {
  try {
    const resp = await http.post(
      '/admin/final-decision-dryrun',
      {
        year: year,
      },
      { responseType: 'blob' },
    );
    const blob = new Blob([resp.data], { type: 'application/zip' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = 'csv.zip';
    link.click();
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
