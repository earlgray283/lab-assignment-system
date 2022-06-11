import {
  createUserWithEmailAndPassword,
  getAuth,
  sendEmailVerification,
  updateProfile,
} from 'firebase/auth';
import React, { useState } from 'react';

import { SignupForm, SignupFormInput } from '../components/forms';
import { DefaultLayout } from '../components/layout';
import { postJson } from '../lib/axios';

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
      await updateProfile(credential.user, {
        displayName: `(${data.studentNumber})${data.name}`,
      });
      data.idToken = await credential.user.getIdToken();
      await postJson('/auth/signup', data);
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
