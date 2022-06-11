import {
  Alert,
  Box,
  Button,
  Stack,
  TextField,
  Typography,
} from '@mui/material';
import React from 'react';
import { Controller, useForm } from 'react-hook-form';

import { TypographyLink } from '../util';

export interface SigninFormInput {
  email: string;
  password: string;
}

export function SigninForm(props: {
  onSubmit: (data: SigninFormInput) => void;
  errorMessage?: string;
  onError?: (e: unknown) => void;
}): JSX.Element {
  const { control, handleSubmit } = useForm<SigninFormInput>();

  return (
    <Stack
      spacing={2}
      component='form'
      onSubmit={() => handleSubmit<SigninFormInput>(props.onSubmit)}
    >
      <Typography variant='h4'>ログイン</Typography>
      <Controller
        name='email'
        defaultValue=''
        control={control}
        rules={{
          required: 'メールアドレスを入力してください',
          pattern: {
            value: /^[\w\-._]+@(shizuoka.ac.jp|inf.shizuoka.ac.jp)$/i,
            message:
              'メールアドレスの形式が正しくないか、静大メールではありません。\n\r(例: shizudai.taro.20@shizuoka.ac.jp)',
          },
        }}
        render={({ field, fieldState }) => (
          <TextField
            required
            fullWidth
            label='email address'
            type='email'
            placeholder='shizudai.taro.20@shizuoka.ac.jp'
            error={fieldState.error !== undefined}
            helperText={fieldState.error?.message}
            {...field}
          />
        )}
      />

      <Controller
        name='password'
        defaultValue=''
        control={control}
        rules={{
          required: 'パスワードを入力してください',
          minLength: {
            value: 8,
            message: 'パスワードは8文字以上である必要があります',
          },
        }}
        render={({ field, fieldState }) => (
          <TextField
            required
            fullWidth
            autoComplete='current-passowrd'
            label='password'
            type='password'
            error={fieldState.error !== undefined}
            helperText={fieldState.error?.message}
            {...field}
          />
        )}
      />

      <Box
        sx={{
          display: 'flex',
          justifyContent: 'left',
          alignItems: 'left',
          flexDirection: 'row',
        }}
      >
        <TypographyLink to='/auth/signup' color='#1976d2'>
          新規登録をする
        </TypographyLink>
      </Box>

      <Button color='primary' type='submit' variant='contained'>
        Sign in
      </Button>

      {props.errorMessage && (
        <Alert severity='error'>{props.errorMessage}</Alert>
      )}
    </Stack>
  );
}
