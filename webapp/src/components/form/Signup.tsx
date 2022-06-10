import React from 'react';
import { Stack, TextField, Button, Alert, Typography } from '@mui/material';
import { Controller, useForm } from 'react-hook-form';

export interface SignupFormInput {
  email: string;
  password: string;
  confirmPassword: string;
  studentNumber: number;
  name: string;
}

export function SignupForm(props: {
  onSubmit: (data: SignupFormInput) => void;
  errorMessage?: string;
  onError?: (e: unknown) => void;
}): JSX.Element {
  const { control, handleSubmit, watch } = useForm<SignupFormInput>();

  return (
    <Stack
      spacing={2}
      component='form'
      onSubmit={handleSubmit<SignupFormInput>(props.onSubmit)}
    >
      <Typography variant='h4'>新規登録</Typography>

      <Controller
        name='email'
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
            helperText={
              fieldState.error?.message ??
              'メールアドレスは静大メール(@shizuoka.ac.jp)のみ使用可能です'
            }
            {...field}
          />
        )}
      />

      <Controller
        name='password'
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
            label='password'
            type='password'
            error={fieldState.error !== undefined}
            helperText={
              fieldState.error?.message ??
              'パスワードは8文字以上、半角英数字をそれぞれ1種類以上使用して入力してください。'
            }
            {...field}
          />
        )}
      />

      <Controller
        name='confirmPassword'
        control={control}
        rules={{
          required: 'パスワードを入力してください(確認用)',
          minLength: {
            value: 8,
            message: 'パスワードは8文字以上である必要があります',
          },
          validate: {
            message: (input) =>
              input === watch().password ? true : 'パスワードが一致しません',
          },
        }}
        render={({ field, fieldState }) => (
          <TextField
            required
            fullWidth
            label='password(確認用)'
            type='password'
            error={fieldState.error !== undefined}
            helperText={
              fieldState.error?.message ??
              '先ほど入力したパスワードを入力してください'
            }
            {...field}
          />
        )}
      />

      <Controller
        name='studentNumber'
        control={control}
        rules={{
          required: '学籍番号を入力してください',
          pattern: {
            value: /^\d{8,8}$/i,
            message: '学籍番号の形式が正しくありません',
          },
        }}
        render={({ field, fieldState }) => (
          <TextField
            required
            fullWidth
            label='学籍番号'
            type='text'
            placeholder='700100xx'
            error={fieldState.error !== undefined}
            helperText={
              fieldState.error?.message ??
              '学籍番号は8桁の数字です(例: 700100xx)'
            }
            {...field}
          />
        )}
      />

      <Controller
        name='name'
        control={control}
        rules={{
          required: '氏名を入力してください',
          maxLength: 20,
        }}
        render={({ field, fieldState }) => (
          <TextField
            required
            fullWidth
            label='氏名'
            type='text'
            placeholder='静大太郎'
            error={fieldState.error !== undefined}
            helperText={
              fieldState.error?.message ??
              'スペースを入れずに入力してください(例: 静大太郎)'
            }
            {...field}
          />
        )}
      />

      <Button color='primary' type='submit' variant='contained'>
        Sign up
      </Button>

      {props.errorMessage && (
        <Alert severity='error'>{props.errorMessage}</Alert>
      )}
    </Stack>
  );
}
