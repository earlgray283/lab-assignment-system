import { Box } from '@mui/material';
import { ReactNode } from 'react';
import React from 'react';

export function DefaultLayout(props: { children?: ReactNode }): JSX.Element {
  return (
    <Box
      sx={{
        margin: '10px',
        minWidth: '350px',
        maxWidth: '600px',
        width: '60%',
        boxShadow: 2,
        padding: '20px',
      }}
    >
      {props.children}
    </Box>
  );
}

export function FullLayout(props: { children?: ReactNode }): JSX.Element {
  return (
    <Box
      margin='10px 0px 0px 0px'
      width='95%'
      maxWidth='1200px'
      boxShadow={2}
      padding='20px'
    >
      {props.children}
    </Box>
  );
}
