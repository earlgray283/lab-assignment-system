import { Alert, Stack } from '@mui/material';
import { useContext, useEffect, useState } from 'react';
import React from 'react';
import { Link, useNavigate } from 'react-router-dom';

import { LoadingDispatchContext, UserContext } from '../App';
import GpaCard from '../components/cards/GpaCard';
import { FullLayout } from '../components/layout';
import { fetchGpa } from '../apis/grade';

function Dashboard(): JSX.Element {
  const setLoading = useContext(LoadingDispatchContext);
  const [gpa, setGpa] = useState<number | null | undefined>(undefined);
  const user = useContext(UserContext);
  const navigate = useNavigate();
  useEffect(() => {
    if (user === null) {
      navigate('/auth/signin');
    }
  }, [user]);

  useEffect(() => {
    setLoading(true);
    (async () => {
      try {
        const gpa2 = await fetchGpa();
        setGpa(gpa2);
      } catch (e) {
        if (e instanceof Error) {
          if (e.message === 'there are no grades') {
            setGpa(null);
          }
        }
      }
    })();

    return () => setLoading(false);
  }, []);

  if (!user) {
    return <div />; // /auth/signin にリダイレクトされることが保証される
  }

  return (
    <FullLayout>
      <Stack spacing={2}>
        {!user.emailVerified && (
          <Alert severity='error'>
            {user.email}{' '}
            宛に確認リンクを送信しました。メールアドレスの確認をしてください
          </Alert>
        )}
        {gpa === null && (
          <Alert severity='error'>
            成績情報が登録されていないようです。
            <Link to='/'>成績情報の登録ページ</Link>
            から登録作業を行なって下さい。
          </Alert>
        )}
        {gpa && <GpaCard data={[1, 2, 3, 4, 3, 2, 1, 0]} gpa={gpa} />}
      </Stack>
    </FullLayout>
  );
}

export default Dashboard;
