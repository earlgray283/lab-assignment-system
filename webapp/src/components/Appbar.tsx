import {
  AppBar,
  Box,
  Button,
  Divider,
  Menu,
  MenuItem,
  Toolbar,
} from '@mui/material';
import React, { useContext, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { UserContext, UserDispatchContext } from '../App';
import { signout } from '../apis/auth';
import { TypographyLink } from './util';

export function Appbar(): JSX.Element {
  const user = useContext(UserContext);
  const setCurrentUser = useContext(UserDispatchContext);
  const navigate = useNavigate();
  const [anchorEl, setAnchorEl] = useState<HTMLElement | null>(null);
  const handleSignout = async () => {
    setAnchorEl(null);
    await signout();
    setCurrentUser(null);
    navigate('/auth/signin');
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

          <div style={{ marginLeft: '20px' }} />

          <TypographyLink to='/labs' variant='subtitle1'>
            研究室一覧
          </TypographyLink>

          <Box sx={{ flexGrow: 1 }} />

          {user !== undefined &&
            (user ? (
              <Button
                onClick={(event) => setAnchorEl(event.currentTarget)}
                sx={{ color: 'white' }}
              >
                {user.uid ?? '<名前未設定>'}
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
            <MenuItem onClick={() => navigate('/profile')}>Profile</MenuItem>
            <MenuItem onClick={() => navigate('/admin')}>Admin</MenuItem>
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
