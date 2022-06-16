import { FirebaseError } from 'firebase/app';
import {
  AuthErrorCodes,
  getAuth,
  signInWithEmailAndPassword,
} from 'firebase/auth';
import React, { useContext, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { LoadingDispatchContext } from '../App';
import { SigninForm, SigninFormInput } from '../components/forms';
import { DefaultLayout } from '../components/layout';

function Signin(): JSX.Element {
  const auth = getAuth();
  const [errorMessage, setErrorMessage] = useState<string | undefined>();
  const navigate = useNavigate();
  const setLoading = useContext(LoadingDispatchContext);
  const onSubmit = async (data: SigninFormInput) => {
    setLoading(true);
    try {
      await signInWithEmailAndPassword(auth, data.email, data.password);
      navigate('/');
    } catch (e: unknown) {
      console.error(e);
      if (e instanceof FirebaseError) {
        switch (e.code) {
          case AuthErrorCodes.INVALID_PASSWORD:
            setErrorMessage('パスワードが正しくありません');
            break;
          default:
            setErrorMessage(
              `未知のエラーです。管理者に連絡をしてください。(code: ${e.code})`
            );
        }
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
