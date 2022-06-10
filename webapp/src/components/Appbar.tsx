import { AppBar, Box, Toolbar } from '@mui/material';
import React, { useContext } from 'react';
import { UserContext } from '../App';
import { TypographyLink } from './util';

export function Appbar(): JSX.Element {
  const user = useContext(UserContext);

  return (
    <Box sx={{ flexGlow: 1 }}>
      <AppBar
        position='static'
        sx={{ backgroundColor: '#99b7dc', boxShadow: 'none' }}
      >
        <Toolbar>
          <TypographyLink to='/' variant='h6'>
            Lab assignment system
          </TypographyLink>

          <Box sx={{ flexGrow: 1 }} />

          <TypographyLink to='/auth/signin'>
            {user ? user.displayName : 'SIGN IN'}
          </TypographyLink>
        </Toolbar>
      </AppBar>
    </Box>
  );
}
