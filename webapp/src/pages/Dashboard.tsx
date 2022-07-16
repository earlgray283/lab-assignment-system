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
        {user.gpa === undefined && (
          <Alert severity='error'>
            成績情報が登録されていないようです。
            <Link to='/profile/register-grades'>成績情報の登録ページ</Link>
            から登録作業を行って下さい。
          </Alert>
        )}
        {user.gpa && <GpaCard gpa={user.gpa} />}
        {user.gpa && (
          <LabCard
            labIds={[user.lab1, user.lab2, user.lab3]}
            gpa={user.gpa}
          />
        )}
      </Stack>
    </FullLayout>
  );
}

export default Dashboard;
