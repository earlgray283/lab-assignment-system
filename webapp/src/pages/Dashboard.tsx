import { Alert, Stack } from '@mui/material';
import { useContext, useEffect, useState } from 'react';
import React from 'react';
import { Link, useNavigate } from 'react-router-dom';

import { UserContext } from '../App';
import GpaCard from '../components/cards/GpaCard';
import { FullLayout } from '../components/layout';
import LabCard from '../components/cards/LabCard';

function Dashboard(): JSX.Element {
  const user = useContext(UserContext);
  const [labIds, setLabIds] = useState<string[] | null | undefined>(undefined);
  const navigate = useNavigate();
  useEffect(() => {
    if (user === null) {
      navigate('/auth/signin');
    }
    if (user) {
      if (user.lab1 && user.lab2 && user.lab3) {
        setLabIds([user.lab1, user.lab2, user.lab3]);
      } else {
        setLabIds(null);
      }
    }
  }, [user]);

  if (!user) {
    return <div />; // /auth/signin にリダイレクトされることが保証される
  }

  console.log(labIds);

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
        {labIds && <LabCard labIds={[...labIds]} gpa={user.gpa} />}
      </Stack>
    </FullLayout>
  );
}

export default Dashboard;
