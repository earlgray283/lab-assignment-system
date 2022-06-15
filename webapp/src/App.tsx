import { Box, CssBaseline, LinearProgress } from '@mui/material';
import { User as FirebaseUser, getAuth } from 'firebase/auth';
import { ApiUser } from './apis/models/user';
import React, {
  createContext,
  lazy,
  Suspense,
  useEffect,
  useState,
} from 'react';
import { BrowserRouter, Route, Routes } from 'react-router-dom';

import { Appbar } from './components/Appbar';
import RegisterGrades from './pages/RegisterGrades';
import { fetchUser } from './apis/user';
const Dashboard = lazy(() => import('./pages/Dashboard'));
const NotFound = lazy(() => import('./pages/NotFound'));
const Profile = lazy(() => import('./pages/Profile'));
const Signin = lazy(() => import('./pages/Signin'));
const Signup = lazy(() => import('./pages/Signup'));

interface User {
  firebaseUser: FirebaseUser;
  apiUser: ApiUser;
}

export const UserContext = createContext<User | null | undefined>(undefined);
export const LoadingStateContext = createContext(false);
export const LoadingDispatchContext = createContext<
  React.Dispatch<React.SetStateAction<boolean>>
>(() => undefined);

function App(): JSX.Element {
  const [loading, setLoading] = useState(true);
  const [currentUser, setCurrentUser] = useState<User | null | undefined>(
    undefined
  );
  const auth = getAuth();

  useEffect(() => {
    auth.onAuthStateChanged(async (firebaseUser) => {
      if (firebaseUser) {
        const apiUser = await fetchUser();
        setCurrentUser({ firebaseUser, apiUser });
      } else {
        setCurrentUser(null);
      }
      setLoading(false);
    });
  }, []);

  return (
    <LoadingStateContext.Provider value={loading}>
      <LoadingDispatchContext.Provider value={setLoading}>
        <UserContext.Provider value={currentUser}>
          {/* reset css */}
          <CssBaseline />

          <BrowserRouter>
            <Appbar />
            {loading && <LinearProgress />}
            <Suspense fallback={<LinearProgress />}>
              <Box
                sx={{
                  display: 'flex',
                  justifyContent: 'center',
                }}
              >
                <Routes>
                  <Route path='*' element={<NotFound />} />
                  <Route path='/'>
                    <Route index element={<Dashboard />} />

                    <Route path='auth'>
                      <Route path='signup' element={<Signup />} />
                      <Route path='signin' element={<Signin />} />
                    </Route>

                    <Route path='profile'>
                      <Route index element={<Profile />} />
                      <Route
                        path='register-grades'
                        element={<RegisterGrades />}
                      />
                    </Route>
                  </Route>
                </Routes>
              </Box>
            </Suspense>
          </BrowserRouter>
        </UserContext.Provider>
      </LoadingDispatchContext.Provider>
    </LoadingStateContext.Provider>
  );
}

export default App;
