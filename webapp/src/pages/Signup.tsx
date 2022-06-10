import React, { useState } from 'react';
import {
  createUserWithEmailAndPassword,
  getAuth,
  sendEmailVerification,
} from 'firebase/auth';
import { DefaultLayout } from '../components/layout';
import { SignupForm, SignupFormInput } from '../components/form';

function Signup(): JSX.Element {
  const auth = getAuth();
  const [errorMessage, setErrorMessage] = useState<string | undefined>();
  const onSubmit = async (data: SignupFormInput) => {
    try {
      const credential = await createUserWithEmailAndPassword(
        auth,
        data.email,
        data.password
      );
      await sendEmailVerification(credential.user, {
        url: `${import.meta.env.VITE_HOST}`,
      });
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
