import { Alert, Stack } from '@mui/material';
import { createContext, useContext, useEffect, useState } from 'react';
import React from 'react';
import { useNavigate } from 'react-router-dom';

import { UserContext } from '../App';
import GpaCard from '../components/cards/GpaCard';
import { FullLayout } from '../components/layout';
import { Notification } from '../components/Notification';
import LabCard from '../components/cards/LabCard';

export const NotificationsContext = createContext<Notification[]>([]);
export const NotificationsDispatchContext = createContext<
  React.Dispatch<React.SetStateAction<Notification[]>>
>(() => undefined);

function Dashboard(): JSX.Element {
  const user = useContext(UserContext);
  const [notifications, setNotifications] = useState<Notification[]>([]);
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
  if (labIds === undefined) {
    return <div>loading...</div>;
  }

  return (
    <FullLayout>
      <NotificationsDispatchContext.Provider value={setNotifications}>
        <Stack spacing={2}>
          {notifications.map((notification, i) => (
            <Alert key={i} severity={notification.severity}>
              {notification.message}
            </Alert>
          ))}

          {user.gpa && <GpaCard gpa={user.gpa} />}
          {<LabCard labIds={labIds ?? undefined} gpa={user.gpa} />}
        </Stack>
      </NotificationsDispatchContext.Provider>
    </FullLayout>
  );
}

export default Dashboard;
