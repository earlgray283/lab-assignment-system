import {
  AppBar,
  Box,
  Button,
  Divider,
  IconButton,
  Menu,
  MenuItem,
  Toolbar,
} from '@mui/material';
import React, { useContext, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import GitHubIcon from '@mui/icons-material/GitHub';
import { UserContext } from '../App';
import { TypographyLink } from './util';
import { signout } from '../apis/auth';

const repoLink = 'https://github.com/earlgray283/lab-assignment-system';

export function Appbar(): JSX.Element {
  const user = useContext(UserContext);
  const navigation = useNavigate();
  const [anchorEl, setAnchorEl] = useState<HTMLElement | null>(null);
  const handleSignout = async () => {
    setAnchorEl(null);
    await signout();
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
          <IconButton
            onClick={() => window.open(repoLink, '_blank')}
            sx={{ color: 'white' }}
          >
            <GitHubIcon />
          </IconButton>

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
