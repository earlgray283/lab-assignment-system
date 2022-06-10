import { Box, CssBaseline, LinearProgress } from '@mui/material';
import { getAuth, User } from 'firebase/auth';
import React, { createContext, useEffect, useState } from 'react';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { Appbar } from './components/Appbar';
import Dashboard from './pages/Dashboard';
import Signin from './pages/Signin';
import Signup from './pages/Signup';

interface UserContextType {
  currentUser: User | null;
  loading: boolean;
}

export const UserContext = createContext<UserContextType>({
  currentUser: null,
  loading: true,
});

function App(): JSX.Element {
  const [loading, setLoading] = useState(true);
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const auth = getAuth();

  useEffect(() => {
    auth.onAuthStateChanged((user) => {
      setLoading(false);
      setCurrentUser(user);
    });
  }, []);

  return (
    <UserContext.Provider value={{ currentUser, loading }}>
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
          </Routes>
        </Box>
      </BrowserRouter>
    </UserContext.Provider>
  );
}

export default App;
