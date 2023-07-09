import { MenuItem, Select, Typography } from '@mui/material';
import { DataGrid, GridColDef } from '@mui/x-data-grid';
import React, { useEffect, useState } from 'react';
import { fetchLabs } from '../apis/labs';
import { Lab } from '../apis/models/lab';
import { DefaultLayout } from '../components/Layout';
import { DisplayGpa } from '../components/util';
import CheckIcon from '@mui/icons-material/Check';

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
    valueGetter: (params) => {
      if (!params.row.userGPAs) {
        return 0;
      }
      return params.row.userGPAs.length;
    },
  },
  {
    field: 'mag',
    headerName: '競争率',
    width: 60,
    align: 'center',
    headerAlign: 'center',
    valueGetter: (params) => {
      if (!params.row.userGPAs) {
        return 0;
      }
      return (params.row.userGPAs.length / params.row.capacity) * 100;
    },
    renderCell: (params) => {
      let mag = 0;
      if (params.row.userGPAs) {
        mag = Math.round(
          (params.row.userGPAs.length / params.row.capacity) * 100
        );
      }
      return <span>{mag}%</span>;
    },
  },
  {
    field: 'minGpa',
    headerName: '最小GPA',
    width: 80,
    align: 'center',
    headerAlign: 'center',
    valueGetter: (params) => {
      if (!params.row.userGPAs) {
        return -1;
      }
      params.row.userGPAs.sort((a, b) => b.gpa - a.gpa);
      return (
        params.row.userGPAs.at(params.row.capacity - 1) ??
        params.row.userGPAs.at(params.row.userGPAs.length - 1) ??
        -1
      );
    },
    renderCell: (params) => {
      let gpa = -1;
      if (params.row.userGPAs) {
        params.row.userGPAs.sort((a, b) => b.gpa - a.gpa);
        gpa =
          params.row.userGPAs.at(params.row.capacity - 1)?.gpa ??
          params.row.userGPAs.at(params.row.userGPAs.length - 1)?.gpa ??
          -1;
      }
      return <DisplayGpa gpa={gpa} />;
    },
  },
  {
    field: 'confirmed',
    headerName: '確定',
    width: 80,
    align: 'center',
    headerAlign: 'center',
    renderCell: (params) => {
      if (!params.row.userGPAs) {
        return <div />;
      }
      return params.row.capacity == params.row.userGPAs.length && <CheckIcon />;
    },
  },
];

function LabListPage(): JSX.Element {
  const [labList, setLabList] = useState<Lab[]>([]);
  const [year, setYear] = useState(2023);
  useEffect(() => {
    (async () => {
      const labList2 = await fetchLabs(year, undefined, ['grade']);
      setLabList(labList2.labs ?? []);
    })();
  }, [year]);

  return (
    <DefaultLayout>
      <Typography variant='h4'>研究室一覧</Typography>

      <Select
        value={year}
        onChange={(e) => setYear(e.target.value as number)}
        sx={{ marginY: 2 }}
      >
        <MenuItem value={2023}>2023年</MenuItem>
        <MenuItem value={2022}>2022年</MenuItem>
      </Select>

      <DataGrid
        rows={labList}
        columns={columns}
        sx={{ height: 800, width: '100%' }}
      />
    </DefaultLayout>
  );
}

export default LabListPage;
