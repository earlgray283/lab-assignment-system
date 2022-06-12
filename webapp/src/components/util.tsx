import { Typography } from '@mui/material';
import React, { ReactNode } from 'react';
import { Link } from 'react-router-dom';

export const TypographyLink = (props: {
  to: string;
  variant?:
    | 'h1'
    | 'h2'
    | 'h3'
    | 'h4'
    | 'h5'
    | 'h6'
    | 'subtitle1'
    | 'subtitle2'
    | 'body1'
    | 'body2'
    | 'caption'
    | 'button'
    | 'overline'
    | 'inherit';
  children?: ReactNode;
  color?: string;
}): JSX.Element => (
  <Typography
    component={Link}
    to={props.to}
    variant={props.variant}
    sx={{
      textDecoration: 'none',
      boxShadow: 'none',
      color: props.color ?? 'white',
    }}
  >
    {props.children}
  </Typography>
);
