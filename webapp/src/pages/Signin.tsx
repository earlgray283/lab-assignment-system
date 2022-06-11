import React, { useContext, useState } from 'react';
import { getAuth, signInWithEmailAndPassword } from 'firebase/auth';
import { useNavigate } from 'react-router-dom';
import { DefaultLayout } from '../components/layout';
import { SigninForm, SigninFormInput } from '../components/forms';
import { LoadingDispatchContext } from '../App';

function Signin(): JSX.Element {
  const auth = getAuth();
  const [errorMessage, setErrorMessage] = useState<string | undefined>();
  const navigate = useNavigate();
  const setLoading = useContext(LoadingDispatchContext);
  const onSubmit = async (data: SigninFormInput) => {
    setLoading(true);
    try {
      await signInWithEmailAndPassword(auth, data.email, data.password);
    } catch (e: unknown) {
      setErrorMessage(`${e}`);
    }
    setLoading(false);
    navigate('/');
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
