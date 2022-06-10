import React, { useState } from 'react';
import { getAuth, signInWithEmailAndPassword } from 'firebase/auth';
import { useNavigate } from 'react-router-dom';
import { DefaultLayout } from '../components/layout';
import { SigninForm, SigninFormInput } from '../components/forms';

function Signin(): JSX.Element {
  const auth = getAuth();
  const [errorMessage, setErrorMessage] = useState<string | undefined>();
  const navigate = useNavigate();
  const onSubmit = async (data: SigninFormInput) => {
    try {
      await signInWithEmailAndPassword(auth, data.email, data.password);
      navigate('/');
    } catch (e: unknown) {
      setErrorMessage(`${e}`);
    }
  };

  return (
    <DefaultLayout>
      <SigninForm
        onSubmit={onSubmit}
        errorMessage={errorMessage}
        onError={(e) => setErrorMessage(`${e}`)}
      />
    </DefaultLayout>
  );
}

export default Signin;
