import { Button, MenuItem, Select } from '@mui/material';
import { DefaultLayout } from '../components/Layout';
import React, { useState } from 'react';
import { finalDecisionDryRun } from '../apis/admin';

export function AdminPage(): JSX.Element {
  const [dryruned, setDryruned] = useState(false);
  const [year, setYear] = useState(2023);
  const [errorMessage, setErrorMessage] = useState<string | undefined>(
    undefined,
  );
  return (
    <DefaultLayout>
      <h1>Admin Page</h1>
      {errorMessage && <p style={{ color: 'red' }}>{errorMessage}</p>}
      <Select
        value={year}
        onChange={(e) => setYear(e.target.value as number)}
        sx={{ marginY: 2 }}
      >
        <MenuItem value={2023}>2023年</MenuItem>
        <MenuItem value={2022}>2022年</MenuItem>
      </Select>
      <br />
      <Button
        variant='contained'
        onClick={async () => {
          try {
            await finalDecisionDryRun(year);
            setDryruned(true);
          } catch (e) {
            console.error(e);
            if (e instanceof Error) {
              setErrorMessage(e.message);
            }
          }
        }}
      >
        最終配属(dry-run)
      </Button>
      <br />
      <Button color='warning' variant='contained' disabled={!dryruned}>
        最終配属
      </Button>
    </DefaultLayout>
  );
}
