import { Alert, Stack } from '@mui/material';
import { useContext, useEffect, useState } from 'react';
import React from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';

import { UserContext } from '../App';
import GpaCard from '../components/cards/GpaCard';
import { FullLayout } from '../components/Layout';
import { Notification } from '../components/Notification';
import LabCard from '../components/cards/LabCard';

function Dashboard(): JSX.Element {
  const user = useContext(UserContext);
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  useEffect(() => {
    if (user === null) {
      navigate('/auth/signin');
      return;
    }
  }, [user]);

  if (!user) {
    return <div />; // /auth/signin にリダイレクトされることが保証される
  }

  return (
    <FullLayout>
      <Stack spacing={2}>
        {user.confirmedLab && (
          <Alert severity='success'>
            おめでとうございます。あなたの配属先が確定しました。
          </Alert>
        )}
        {notifications.map((notification, i) => (
          <Alert key={i} severity={notification.severity}>
            {notification.message}
          </Alert>
        ))}

        {user.gpa && (
          <GpaCard
            year={
              searchParams.get('year')
                ? Number(searchParams.get('year'))
                : user.year
            }
          />
        )}
        <LabCard
          year={
            searchParams.get('year')
              ? Number(searchParams.get('year'))
              : user.year
          }
          pushNotification={(n) => {
            if (notifications.find((x) => x.id === n.id)) {
              return;
            }
            setNotifications([...notifications, n]);
          }}
        />
      </Stack>
    </FullLayout>
  );
}

export default Dashboard;
