import { Typography } from '@mui/material';
import { DataGrid, GridColDef } from '@mui/x-data-grid';
import React, { useEffect, useState } from 'react';
import { fetchLabs } from '../apis/labs';
import { LabList, Lab } from '../apis/models/lab';
import { DefaultLayout } from '../components/layout';
import { DisplayGpa } from '../components/util';

const columns: GridColDef<Lab>[] = [
  {
    field: 'name',
    headerName: '研究室',
    align: 'center',
    headerAlign: 'center',
    width: 120,
  },
  {
    field: 'capacity',
    headerName: '定員',
    width: 20,
    align: 'center',
    headerAlign: 'center',
  },
  {
    field: 'firstChoice',
    headerName: '希望者数',
    width: 80,
    align: 'center',
    headerAlign: 'center',
  },
  {
    field: 'mag',
    headerName: '競争率',
    width: 60,
    align: 'center',
    headerAlign: 'center',
    valueGetter: (params) => {
      if (!params.row.grades) {
        return -1;
      }
      const gpas = params.row.grades.gpas1;
      return (gpas.length / params.row.capacity) * 100;
    },
    renderCell: (params) => {
      if (!params.row.grades) {
        return -1;
      }
      const gpas = params.row.grades.gpas1;
      return (
        <span>
          {Math.round((gpas.length / params.row.capacity) * 10000) / 100}%
        </span>
      );
    },
  },
  {
    field: 'minGpa',
    headerName: '最小GPA',
    width: 80,
    align: 'center',
    headerAlign: 'center',
    valueGetter: (params) => {
      if (!params.row.grades) {
        return -1;
      }
      const gpas = params.row.grades.gpas1;
      gpas.sort((a, b) => b - a);
      return (
        gpas.at(params.row.capacity - 1) ??
        gpas.at(params.row.grades.gpas1.length - 1) ??
        -1
      );
    },
    renderCell: (params) => {
      if (!params.row.grades) {
        return -1;
      }
      const gpas = params.row.grades.gpas1;
      gpas.sort((a, b) => b - a);
      return (
        <DisplayGpa
          gpa={
            gpas.at(params.row.capacity - 1) ??
            gpas.at(params.row.grades.gpas1.length - 1) ??
            -1
          }
        />
      );
    },
  },
];

function LabListPage(): JSX.Element {
  const [labList, setLabList] = useState<LabList | undefined>(undefined);
  useEffect(() => {
    (async () => {
      const labList2 = await fetchLabs(undefined, ['grade']);
      setLabList(labList2);
    })();
  }, []);

  if (!labList) {
    return <div />;
  }

  return (
    <DefaultLayout>
      <Typography variant='h4'>研究室一覧</Typography>

      <DataGrid
        rows={labList.labs}
        columns={columns}
        sx={{ height: 800, width: '100%' }}
      />
    </DefaultLayout>
  );
}

export default LabListPage;
