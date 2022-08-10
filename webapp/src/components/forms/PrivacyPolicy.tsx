import { Box, Button } from '@mui/material';
import { DefaultLayout } from '../Layout';
import React from 'react';
import { PrivacyPolicy as PrivacyPolicyCard } from '../cards/PrivacyPolicy';

function PrivacyPolicy(props: { onAgree: () => void }): JSX.Element {
  return (
    <DefaultLayout>
      <PrivacyPolicyCard />
      <Box
        sx={{
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
        }}
      >
        <Button
          variant='contained'
          onClick={() => props.onAgree()}
          sx={{ width: '100px' }}
        >
          同意する
        </Button>
      </Box>
    </DefaultLayout>
  );
}

export default PrivacyPolicy;
