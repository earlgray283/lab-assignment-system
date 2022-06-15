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

export const DisplayGpa = (props: { gpa: number }): JSX.Element => {
  return (
    <strong>
      <span
        style={{
          color:
            props.gpa < 1.5
              ? 'black'
              : props.gpa < 2.0
              ? 'gray'
              : props.gpa < 2.5
              ? 'brown'
              : props.gpa < 2.8
              ? 'green'
              : props.gpa < 3.0
              ? 'cyan'
              : props.gpa < 3.1
              ? 'blue'
              : props.gpa < 3.2
              ? 'yellow'
              : props.gpa < 3.3
              ? 'orange'
              : props.gpa < 3.4
              ? 'red'
              : 'gold',
        }}
      >
        {Math.floor(props.gpa * 1000) / 1000}
      </span>
    </strong>
  );
};
