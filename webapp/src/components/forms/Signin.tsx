import { Alert, Button, Stack, TextField, Typography } from '@mui/material';
import React from 'react';
import { Controller, useForm } from 'react-hook-form';

export interface SigninFormInput {
  uid: string;
}

export function SigninForm(props: {
  onSubmit: (data: SigninFormInput) => void;
  errorMessage?: string;
  onError?: (e: unknown) => void;
}): JSX.Element {
  const { control, handleSubmit } = useForm<SigninFormInput>();

  return (
    <Stack spacing={2} component='form' onSubmit={handleSubmit(props.onSubmit)}>
      <Typography variant='h4'>ログイン</Typography>

      {props.errorMessage && (
        <Alert severity='error'>{props.errorMessage}</Alert>
      )}

      <Controller
        name='uid'
        defaultValue=''
        control={control}
        rules={{
          required: 'IDを入力してください',
        }}
        render={({ field, fieldState }) => (
          <TextField
            required
            fullWidth
            label='uid'
            type='text'
            error={fieldState.error !== undefined}
            helperText={fieldState.error?.message}
            {...field}
          />
        )}
      />

      <Button color='primary' type='submit' variant='contained'>
        Sign in
      </Button>
    </Stack>
  );
}
