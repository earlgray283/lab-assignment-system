import { Stack, Typography, TextField, MenuItem } from '@mui/material';
import { Control, Controller, UseFormWatch } from 'react-hook-form';
import React from 'react';
import { LabList } from '../../apis/models/lab';
import { SignupFormInput } from './Signup';

export interface Props {
  labList: LabList;
  control: Control<SignupFormInput, unknown>;
  watch: UseFormWatch<SignupFormInput>;
}

function LabSurvey(props: Props): JSX.Element {
  return (
    <Stack spacing={2} width='90%' display='flex' flexDirection='column'>
      <Typography variant='h6'>研究室配属アンケート</Typography>
      <Controller
        name='lab1'
        defaultValue=''
        control={props.control}
        render={({ field }) => (
          <TextField defaultValue='' select label='第1希望の研究室' {...field}>
            <MenuItem value={''}>未選択</MenuItem>
            {props.labList.labs.map((lab) => (
              <MenuItem value={lab.id} key={lab.id}>
                {lab.name}
              </MenuItem>
            ))}
          </TextField>
        )}
      />
      <Controller
        name='lab2'
        control={props.control}
        defaultValue=''
        render={({ field }) => (
          <TextField defaultValue='' select label='第2希望の研究室' {...field}>
            <MenuItem value={''}>未選択</MenuItem>
            {props.labList.labs.map((lab) => (
              <MenuItem value={lab.id} key={lab.id}>
                {lab.name}
              </MenuItem>
            ))}
          </TextField>
        )}
      />
      <Controller
        name='lab3'
        control={props.control}
        defaultValue=''
        render={({ field }) => (
          <TextField defaultValue='' select label='第3希望の研究室' {...field}>
            <MenuItem value={''}>未選択</MenuItem>
            {props.labList.labs.map((lab) => (
              <MenuItem value={lab.id} key={lab.id}>
                {lab.name}
              </MenuItem>
            ))}
          </TextField>
        )}
      />
    </Stack>
  );
}

export default LabSurvey;
