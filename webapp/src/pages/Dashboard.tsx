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
        {(!user.lab1 || !user.lab2 || !user.lab3) && (
          <Alert severity='error'>
            研究室アンケートに回答されていません。
            <Link to='/profile'>研究室アンケートページ</Link>
            から回答を行って下さい。
          </Alert>
        )}
        {user.gpa && <GpaCard gpa={user.gpa} />}
        {user.lab1 && user.lab2 && user.lab3 && (
          <LabCard labIds={[user.lab1, user.lab2, user.lab3]} gpa={user.gpa} />
        )}
      </Stack>
    </FullLayout>
  );
}

export default Dashboard;
