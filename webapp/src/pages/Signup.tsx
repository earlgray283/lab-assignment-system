import {
  createUserWithEmailAndPassword,
  getAuth,
  sendEmailVerification,
  updateProfile,
} from 'firebase/auth';
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { signup } from '../apis/auth';

import { SignupForm, SignupFormInput } from '../components/forms';
import { DefaultLayout } from '../components/layout';

function Signup(): JSX.Element {
  const auth = getAuth();
  const [errorMessage, setErrorMessage] = useState<string | undefined>();
  const navigate = useNavigate();
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
        await sendEmailVerification(credential.user);
        navigate('/');
      } catch (e: unknown) {
        await credential.user.delete();
        setErrorMessage(`${e}`);
      }
    } catch (e: unknown) {
      setErrorMessage(`${e}`);
    }
  };

  return (
    <DefaultLayout>
      <SignupForm
        onSubmit={onSubmit}
        errorMessage={errorMessage}
        onError={(e) => setErrorMessage(`${e}`)}
      />
    </DefaultLayout>
  );
}

export default Signup;
