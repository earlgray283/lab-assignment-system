import React, { useContext, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { signin } from '../apis/auth';
import { LoadingDispatchContext } from '../App';
import { SigninForm, SigninFormInput } from '../components/forms';
import { DefaultLayout } from '../components/layout';

function Signin(): JSX.Element {
  const [errorMessage, setErrorMessage] = useState<string | undefined>();
  const navigate = useNavigate();
  const setLoading = useContext(LoadingDispatchContext);
  const onSubmit = async (data: SigninFormInput) => {
    setLoading(true);
    try {
      await signin(data.id);
      navigate('/');
    } catch (e: unknown) {
      console.error(e);
      if (e instanceof Error) {
        setErrorMessage(e.message);
      }
    }
    setLoading(false);
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
