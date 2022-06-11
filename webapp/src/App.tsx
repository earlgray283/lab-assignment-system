import { Box, CssBaseline, LinearProgress } from '@mui/material';
import { User, getAuth } from 'firebase/auth';
import React, { createContext, useEffect, useState } from 'react';
import { BrowserRouter, Route, Routes } from 'react-router-dom';

import { Appbar } from './components/Appbar';
import Dashboard from './pages/Dashboard';
import NotFound from './pages/NotFound';
import Signin from './pages/Signin';
import Signup from './pages/Signup';

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
    auth.onAuthStateChanged((user) => {
      setLoading(false);
      setCurrentUser(user);
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
            <Box
              sx={{
                display: 'flex',
                justifyContent: 'center',
              }}
            >
              <Routes>
                <Route path='/'>
                  <Route index element={<Dashboard />} />

                  <Route path='auth'>
                    <Route path='signup' element={<Signup />} />
                    <Route path='signin' element={<Signin />} />
                  </Route>
                </Route>
                <Route path='*' element={<NotFound />} />
              </Routes>
            </Box>
          </BrowserRouter>
        </UserContext.Provider>
      </LoadingDispatchContext.Provider>
    </LoadingStateContext.Provider>
  );
}

export default App;
