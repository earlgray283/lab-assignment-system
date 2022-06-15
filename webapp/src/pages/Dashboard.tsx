import { Alert, Stack } from '@mui/material';
import { useContext, useEffect } from 'react';
import React from 'react';
import { Link, useNavigate } from 'react-router-dom';

import { UserContext } from '../App';
import GpaCard from '../components/cards/GpaCard';
import { FullLayout } from '../components/layout';
import { fetchGpa } from '../apis/grade';
import { sleep } from '../lib/util';
import LabCard from '../components/cards/LabCard';

let gpa: number | null | undefined;
function useGpa(): number | null {
  if (gpa === undefined) {
    throw fetchGpa()
      .then((data) => (gpa = data))
      .catch((e) => {
        if (e instanceof Error) {
          if (e.message === 'there are no grades') {
            gpa = null;
            return;
          }
        }
        return sleep(2000);
      });
  }
  return gpa;
}

function Dashboard(): JSX.Element {
  const gpa = useGpa();
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

  console.log(user.apiUser, [
    user.apiUser.lab1,
    user.apiUser.lab2,
    user.apiUser.lab3,
  ]);

  return (
    <FullLayout>
      <Stack spacing={2}>
        {!user.firebaseUser.emailVerified && (
          <Alert severity='error'>
            {user.firebaseUser.email}{' '}
            宛に確認リンクを送信しました。メールアドレスの確認をしてください
          </Alert>
        )}
        {gpa === null && (
          <Alert severity='error'>
            成績情報が登録されていないようです。
            <Link to='/profile/register-grades'>成績情報の登録ページ</Link>
            から登録作業を行って下さい。
          </Alert>
        )}
        {gpa && <GpaCard data={[1, 2, 3, 4, 3, 2, 1, 0]} gpa={gpa} />}
        <LabCard
          labIds={[user.apiUser.lab1, user.apiUser.lab2, user.apiUser.lab3]}
        />
      </Stack>
    </FullLayout>
  );
}

export default Dashboard;
