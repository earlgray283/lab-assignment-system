import React, { useContext, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { signin } from '../apis/auth';
import { LoadingDispatchContext, UserDispatchContext } from '../App';
import { SigninForm, SigninFormInput } from '../components/forms';
import { DefaultLayout } from '../components/Layout';

function Signin(): JSX.Element {
  const [errorMessage, setErrorMessage] = useState<string | undefined>();
  const navigate = useNavigate();
  const setLoading = useContext(LoadingDispatchContext);
  const setCurrentUser = useContext(UserDispatchContext);

  const onSubmit = async (data: SigninFormInput) => {
    setLoading(true);
    try {
      const user = await signin(data.uid);
      setCurrentUser(user);
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
