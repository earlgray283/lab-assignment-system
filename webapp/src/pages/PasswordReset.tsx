import { DefaultLayout } from '../components/layout';
import React from 'react';
import { Button, Stack, TextField, Typography } from '@mui/material';
import { useForm, Controller } from 'react-hook-form';
import { getAuth, sendPasswordResetEmail } from 'firebase/auth';
import { useNavigate } from 'react-router-dom';

interface PasswordResetFormInput {
  email: string;
}

function PasswordReset(): JSX.Element {
  const { control, handleSubmit } = useForm<PasswordResetFormInput>();
  const auth = getAuth();
  const navigate = useNavigate();
  return (
    <DefaultLayout>
      <Stack
        spacing={2}
        component='form'
        onSubmit={handleSubmit(async (data) => {
          try {
            await sendPasswordResetEmail(auth, data.email);
            alert(
              `再設定メールを ${data.email} 宛に送信しました。再設定後に再度サインインしてください。`
            );
            navigate('/auth/signin');
          } catch (e) {
            console.error(e);
            alert('再設定メールの送信に失敗しました。');
          }
        })}
      >
        <Typography variant='h4'>パスワード再設定</Typography>
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

        <Button color='primary' type='submit' variant='contained'>
          送信
        </Button>
      </Stack>
    </DefaultLayout>
  );
}

export default PasswordReset;
