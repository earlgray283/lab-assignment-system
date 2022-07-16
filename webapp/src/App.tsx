import { Box, CssBaseline, LinearProgress } from '@mui/material';
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
import { fetchUser } from './apis/user';
const Dashboard = lazy(() => import('./pages/Dashboard'));
const NotFound = lazy(() => import('./pages/NotFound'));
const Profile = lazy(() => import('./pages/Profile'));
const Signin = lazy(() => import('./pages/Signin'));

export const UserContext = createContext<ApiUser | null | undefined>(undefined);
export const UserDispatchContext = createContext<
  React.Dispatch<React.SetStateAction<ApiUser | null | undefined>>
>(() => undefined);
export const LoadingStateContext = createContext(false);
export const LoadingDispatchContext = createContext<
  React.Dispatch<React.SetStateAction<boolean>>
>(() => undefined);

function App(): JSX.Element {
  const [loading, setLoading] = useState(true);
  const [currentUser, setCurrentUser] = useState<ApiUser | null | undefined>(
    undefined
  );

  useEffect(() => {
    setCurrentUser(undefined);
    (async () => {
      try {
        const apiUser = await fetchUser();
        setCurrentUser(apiUser);
      } catch (e) {
        setCurrentUser(null);
      }
      setLoading(false);
    })();
  }, []);

  return (
    <LoadingStateContext.Provider value={loading}>
      <LoadingDispatchContext.Provider value={setLoading}>
        <UserContext.Provider value={currentUser}>
          <UserDispatchContext.Provider value={setCurrentUser}>
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
                        <Route path='signin' element={<Signin />} />
                      </Route>

                      <Route path='profile'>
                        <Route index element={<Profile />} />
                      </Route>
                    </Route>
                  </Routes>
                </Box>
              </Suspense>
            </BrowserRouter>
          </UserDispatchContext.Provider>
        </UserContext.Provider>
      </LoadingDispatchContext.Provider>
    </LoadingStateContext.Provider>
  );
}

export default App;
