import { Stack, TextField, MenuItem } from '@mui/material';
import React, { useEffect, useState } from 'react';
import { LabList } from '../../apis/models/lab';
import { fetchLabs } from '../../apis/labs';
import { sleep } from '../../libs/util';

export interface LabSurveyFormInput {
  lab1: string;
  lab2: string;
  lab3: string;
}

let labList: LabList | undefined;

function useLabList(): LabList {
  if (labList === undefined) {
    throw fetchLabs()
      .then((data) => (labList = data))
      .catch(() => sleep(2000));
  }
  return labList;
}

function LabSurvey(props: {
  onChange: (lab1: string, lab2: string, lab3: string) => void;
  defaultLab1?: string;
  defaultLab2?: string;
  defaultLab3?: string;
}): JSX.Element {
  const [lab1, setLab1] = useState(props.defaultLab1 ?? '');
  const [lab2, setLab2] = useState(props.defaultLab2 ?? '');
  const [lab3, setLab3] = useState(props.defaultLab3 ?? '');
  const labList = useLabList();

  useEffect(() => {
    props.onChange(lab1, lab2, lab3);
  }, [lab1, lab2, lab3]);

  return (
    <Stack spacing={2} width='90%' display='flex' flexDirection='column'>
      <TextField
        defaultValue={props.defaultLab1 ?? ''}
        select
        label='第1希望の研究室'
        onChange={(e) => setLab1(e.target.value)}
      >
        <MenuItem value={''}>未選択</MenuItem>
        {labList.labs.map((lab) => (
          <MenuItem value={lab.id} key={lab.id}>
            {lab.name}
          </MenuItem>
        ))}
      </TextField>

      <TextField
        defaultValue={props.defaultLab2 ?? ''}
        select
        label='第2希望の研究室'
        onChange={(e) => setLab2(e.target.value)}
      >
        <MenuItem value={''}>未選択</MenuItem>
        {labList.labs.map((lab) => (
          <MenuItem value={lab.id} key={lab.id}>
            {lab.name}
          </MenuItem>
        ))}
      </TextField>

      <TextField
        defaultValue={props.defaultLab3 ?? ''}
        select
        label='第3希望の研究室'
        onChange={(e) => setLab3(e.target.value)}
      >
        <MenuItem value={''}>未選択</MenuItem>
        {labList.labs.map((lab) => (
          <MenuItem value={lab.id} key={lab.id}>
            {lab.name}
          </MenuItem>
        ))}
      </TextField>
    </Stack>
  );
}

export default LabSurvey;
