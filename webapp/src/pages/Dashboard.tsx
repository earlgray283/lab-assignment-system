import { Alert, Stack } from '@mui/material';
import { useContext, useEffect } from 'react';
import React from 'react';
import { Link, useNavigate } from 'react-router-dom';

import { UserContext } from '../App';
import GpaCard from '../components/cards/GpaCard';
import { FullLayout } from '../components/layout';
import LabCard from '../components/cards/LabCard';

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
        {!user.firebaseUser.emailVerified && (
          <Alert severity='error'>
            {user.firebaseUser.email}{' '}
            宛に確認リンクを送信しました。メールアドレスの確認をしてください
          </Alert>
        )}
        {user.apiUser.gpa === undefined && (
          <Alert severity='error'>
            成績情報が登録されていないようです。
            <Link to='/profile/register-grades'>成績情報の登録ページ</Link>
            から登録作業を行って下さい。
          </Alert>
        )}
        {user.apiUser.gpa && <GpaCard gpa={user.apiUser.gpa} />}
        {user.apiUser.gpa && (
          <LabCard
            labIds={[user.apiUser.lab1, user.apiUser.lab2, user.apiUser.lab3]}
            gpa={user.apiUser.gpa}
          />
        )}
      </Stack>
    </FullLayout>
  );
}

export default Dashboard;
