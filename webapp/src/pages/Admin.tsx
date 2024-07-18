import { Button, MenuItem, Select } from '@mui/material';
import { DefaultLayout } from '../components/layout';
import React, { useState } from 'react';
import { finalDecisionDryRun } from '../apis/admin';

export function AdminPage(): JSX.Element {
  const [dryruned, setDryruned] = useState(false);
  const [year, setYear] = useState(new Date().getFullYear());
  const [errorMessage, setErrorMessage] = useState<string | undefined>(
    undefined,
  );
  return (
    <DefaultLayout>
      <h1>Admin Page</h1>

      <h2>年</h2>
      <div>
        <Select
          value={year}
          onChange={(e) => setYear(e.target.value as number)}
          sx={{ marginY: 2 }}
        >
          {new Array(10).fill(0).map((_, i) => (
            <MenuItem key={2022 + i} value={2022 + i}>
              {2022 + i}年
            </MenuItem>
          ))}
        </Select>
      </div>

      <h2>データの流し込み</h2>
      <div>TODO: 研究室と学生データ</div>

      <h2>最終配属</h2>
      <div>
        {errorMessage && <p style={{ color: 'red' }}>{errorMessage}</p>}

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
      </div>
    </DefaultLayout>
  );
}
