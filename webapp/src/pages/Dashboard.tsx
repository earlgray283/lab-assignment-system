import { useContext, useEffect } from 'react';
import { UserContext } from '../App';
import { DefaultLayout } from '../components/layout';
import React from 'react';
import { useNavigate } from 'react-router-dom';

function Dashboard(): JSX.Element {
  const user = useContext(UserContext);
  const navigate = useNavigate();
  useEffect(() => {
    if (user === null) {
      navigate('/auth/signin');
    }
  }, [navigate]);

  return <DefaultLayout>s</DefaultLayout>;
}

export default Dashboard;
