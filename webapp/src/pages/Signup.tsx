import { Alert } from '@mui/material';
import {
  createUserWithEmailAndPassword,
  getAuth,
  updateProfile,
} from 'firebase/auth';
import React, { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { signup } from '../apis/auth';

import { SignupForm, SignupFormInput } from '../components/forms';
import { DefaultLayout } from '../components/layout';

function Signup(): JSX.Element {
  const auth = getAuth();
  const [token, setToken] = useState<string | null | undefined>(undefined);
  const [errorMessage, setErrorMessage] = useState<string | undefined>();
  const [params] = useSearchParams();
  const navigate = useNavigate();

  useEffect(() => {
    setToken(params.get('token'));
  }, [params]);

  const onSubmit = async (data: SignupFormInput) => {
    try {
      const credential = await createUserWithEmailAndPassword(
        auth,
        data.email,
        data.password
      );
      await updateProfile(credential.user, {
        displayName: `(${data.studentNumber})${data.name}`,
      });
      try {
        const idToken = await credential.user.getIdToken();
        await signup({
          email: data.email,
          studentNumber: data.studentNumber,
          name: data.name,
          lab1: data.lab1,
          lab2: data.lab2,
          lab3: data.lab3,
          idToken: idToken,
          password: data.password,
        });
        navigate('/');
      } catch (e: unknown) {
        await credential.user.delete();
        setErrorMessage(`${e}`);
      }
    } catch (e: unknown) {
      setErrorMessage(`${e}`);
    }
  };

  if (token === undefined) {
    return <DefaultLayout />;
  } else if (token === null) {
    return (
      <DefaultLayout>
        <Alert severity='error'>メールアドレスの確認が完了していません。</Alert>
      </DefaultLayout>
    );
  }

  return (
    <DefaultLayout>
      <SignupForm
        onSubmit={onSubmit}
        errorMessage={errorMessage}
        onError={(e) => setErrorMessage(`${e}`)}
        token={token}
      />
    </DefaultLayout>
  );
}

export default Signup;
