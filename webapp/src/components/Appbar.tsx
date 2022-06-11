import {
  AppBar,
  Box,
  Button,
  Divider,
  Menu,
  MenuItem,
  Toolbar,
} from '@mui/material';
import { getAuth } from 'firebase/auth';
import React, { useContext, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { UserContext } from '../App';
import { TypographyLink } from './util';

export function Appbar(): JSX.Element {
  const user = useContext(UserContext);
  const auth = getAuth();
  const navigation = useNavigate();
  const [anchorEl, setAnchorEl] = useState<HTMLElement | null>(null);
  const handleSignout = async () => {
    setAnchorEl(null);
    await auth.signOut();
  };

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

          {user !== undefined &&
            (user ? (
              <Button
                onClick={(event) => setAnchorEl(event.currentTarget)}
                sx={{ color: 'white' }}
              >
                {user.displayName ?? '<名前未設定>'}
              </Button>
            ) : (
              <TypographyLink to='/auth/signin'>SIGN IN</TypographyLink>
            ))}
          <Menu
            id='basic-menu'
            anchorEl={anchorEl}
            open={anchorEl !== null}
            onClose={() => setAnchorEl(null)}
            MenuListProps={{
              'aria-labelledby': 'basic-button',
            }}
          >
            <MenuItem onClick={() => navigation('/profile')}>Profile</MenuItem>
            <Divider />
            <MenuItem onClick={async () => await handleSignout()}>
              Logout
            </MenuItem>
          </Menu>
        </Toolbar>
      </AppBar>
    </Box>
  );
}
