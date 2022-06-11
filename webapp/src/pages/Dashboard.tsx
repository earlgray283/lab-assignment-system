import { Alert, Stack } from '@mui/material';
import { useContext, useEffect } from 'react';
import React from 'react';
import { useNavigate } from 'react-router-dom';

import { UserContext } from '../App';
import GpaCard from '../components/cards/GpaCard';
import { FullLayout } from '../components/layout';

function Dashboard(): JSX.Element {
  const user = useContext(UserContext);
  const navigate = useNavigate();
  useEffect(() => {
    if (user === null) {
      navigate('/auth/signin');
    }
  }, [user]);
  if (!user) {
    return <div />; // /auth/signin にリダイレクトされることが保証される
  }

  return (
    <FullLayout>
      <Stack spacing={2}>
        {user.emailVerified && (
          <Alert severity='error'>
            {user.email}{' '}
            宛に確認リンクを送信しました。メールアドレスの確認をしてください
          </Alert>
        )}
        <GpaCard data={[1, 2, 3, 4, 3, 2, 1, 0]} />
      </Stack>
    </FullLayout>
  );
}

export default Dashboard;
