import { Stack, Typography, TextField, Button } from '@mui/material';
import React, { useContext } from 'react';
import { Controller, useForm } from 'react-hook-form';
import { emailVerification } from '../apis/auth';
import { LoadingDispatchContext } from '../App';
import { DefaultLayout } from '../components/layout';

interface FormInput {
  email: string;
}

function EmailVerification(): JSX.Element {
  const { control, handleSubmit } = useForm<FormInput>();
  const setLoading = useContext(LoadingDispatchContext);
  return (
    <DefaultLayout>
      <Stack
        spacing={2}
        component='form'
        onSubmit={handleSubmit(async (data) => {
          setLoading(true);
          try {
            await emailVerification(data.email);
            alert(`確認メールを ${data.email} 宛に送信しました。`);
          } catch (e) {
            console.error(e);
            alert('確認メールの送信に失敗しました。');
          } finally {
            setLoading(false);
          }
        })}
      >
        <Typography variant='h4'>メールアドレス認証</Typography>
        <Typography variant='body1'>
          貴方が静岡大学の学生であることを確認するため、入力頂いたメールアドレス宛に確認メールを送信します。
        </Typography>
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

export default EmailVerification;
