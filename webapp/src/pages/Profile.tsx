import { Box, Button, Typography } from '@mui/material';
import React, { useContext, useState } from 'react';
import { updateUser } from '../apis/user';
import { User, UserContext } from '../App';
import LabSurvey from '../components/forms/LabSurvey';
import { DefaultLayout } from '../components/layout';
import { sleep } from '../lib/util';

export interface LabSurveyFormInput {
  lab1: string;
  lab2: string;
  lab3: string;
}

let user: User | null | undefined;

function useUser(): User {
  user = useContext(UserContext);
  if (!user) {
    throw sleep(2000);
  }
  return user;
}

function Profile(): JSX.Element {
  const user = useUser();
  const [errorMessage, setErrorMessage] = useState<string | undefined>(
    undefined
  );
  const [labSurvey, setLabSurvey] = useState<LabSurveyFormInput>({
    lab1: user.apiUser.lab1,
    lab2: user.apiUser.lab2,
    lab3: user.apiUser.lab3,
  });

  return (
    <DefaultLayout>
      <Typography variant='h4'>Profile</Typography>
      <Box>
        <Typography variant='h6' marginBottom='10px'>
          研究室アンケートの変更
        </Typography>
        <LabSurvey
          onChange={(lab1, lab2, lab3) => {
            console.log(lab1, lab2, lab3);
            setLabSurvey({
              lab1,
              lab2,
              lab3,
            });
          }}
          defaultLab1={labSurvey.lab1}
          defaultLab2={labSurvey.lab2}
          defaultLab3={labSurvey.lab3}
        />
        <Button
          variant='contained'
          sx={{ marginY: '10px' }}
          onClick={async () => {
            const newUser = user.apiUser;
            newUser.lab1 = labSurvey.lab1;
            newUser.lab2 = labSurvey.lab2;
            newUser.lab3 = labSurvey.lab3;
            try {
              await updateUser(newUser);
            } catch (e) {
              if (e instanceof Error) {
                setErrorMessage(e.message);
              }
            }
          }}
        >
          保存する
        </Button>
      </Box>
    </DefaultLayout>
  );
}

export default Profile;
