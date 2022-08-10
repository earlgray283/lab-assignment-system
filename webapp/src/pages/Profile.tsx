import { Alert, Box, Button, Typography } from '@mui/material';
import React, { useContext, useState } from 'react';
import { ApiUser, UserLab } from '../apis/models/user';
import { updateUserLab } from '../apis/user';
import { UserContext, UserDispatchContext } from '../App';
import LabSurvey from '../components/forms/LabSurvey';
import { DefaultLayout } from '../components/Layout';
import { sleep } from '../libs/util';

export interface LabSurveyFormInput {
  lab1: string;
  lab2: string;
  lab3: string;
}

let user: ApiUser | null | undefined;

function useUser(): ApiUser {
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
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [labSurvey, setLabSurvey] = useState<LabSurveyFormInput>({
    lab1: user.lab1 ?? '',
    lab2: user.lab2 ?? '',
    lab3: user.lab3 ?? '',
  });
  const setCurrentUser = useContext(UserDispatchContext);

  return (
    <DefaultLayout>
      <Typography variant='h4'>Profile</Typography>
      <Box>
        {errorMessage && <Alert severity='error'>errorMessage</Alert>}
        {successMessage && <Alert severity='success'>{successMessage}</Alert>}
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
            const userLab: UserLab = {
              lab1: labSurvey.lab1,
              lab2: labSurvey.lab2,
              lab3: labSurvey.lab3,
            };
            try {
              const user = await updateUserLab(userLab);
              console.log(user);
              setCurrentUser(user);
              setSuccessMessage('更新に成功しました');
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
