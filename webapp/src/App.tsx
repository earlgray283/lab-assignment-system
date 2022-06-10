import { Box, CssBaseline } from '@mui/material';
import { getAuth, User } from 'firebase/auth';
import React, { createContext, useEffect, useState } from 'react';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { Appbar } from './components/Appbar';
import Dashboard from './pages/Dashboard';
import Signin from './pages/Signin';
import Signup from './pages/Signup';

export const UserContext = createContext<User | null>(null);

function App(): JSX.Element {
  const [user, setUser] = useState<User | null>(null);
  const auth = getAuth();

  useEffect(() => {
    auth.onAuthStateChanged((currentUser) => {
      setUser(currentUser);
    });
  }, []);

  return (
    <UserContext.Provider value={user}>
      {/* reset css */}
      <CssBaseline />

      <BrowserRouter>
        <Appbar />
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
          </Routes>
        </Box>
      </BrowserRouter>
    </UserContext.Provider>
  );
}

export default App;
