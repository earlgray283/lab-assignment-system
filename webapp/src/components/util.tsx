import { Typography } from '@mui/material';
import React, { ReactNode, useEffect, useState } from 'react';
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

export const useWindowSize = (): [number, number] => {
  const [innerWidth, setInnerWidth] = useState(window.innerWidth);
  const [innerHeight, setInnerHeight] = useState(window.innerHeight);
  useEffect(() => {
    const onResize = () => {
      setInnerWidth(window.innerWidth);
      setInnerHeight(window.innerHeight);
    };
    window.addEventListener('resize', onResize);
    return () => window.removeEventListener('resize', onResize);
  }, []);

  return [innerWidth, innerHeight];
};
