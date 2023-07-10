import {
  Alert,
  Box,
  Button,
  MenuItem,
  Stack,
  TextField,
  Typography,
} from '@mui/material';
import React, { useContext, useEffect, useState } from 'react';
import { ApiUser } from '../apis/models/user';
import { updateUserLab } from '../apis/user';
import { UserContext, UserDispatchContext } from '../App';
import { DefaultLayout } from '../components/Layout';
import { sleep } from '../libs/util';
import { fetchLabs } from '../apis/labs';
import { Lab } from '../apis/models/lab';
import { useSearchParams } from 'react-router-dom';

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
  const [wishLab, setWishLab] = useState<string | undefined>(undefined);
  const setCurrentUser = useContext(UserDispatchContext);
  const [labList, setLabList] = useState<Lab[]>([]);
  const [searchParams] = useSearchParams();

  useEffect(() => {
    const year = searchParams.get('year');
    (async () => {
      const labList2 = await fetchLabs(
        year ? Number(year) : user.year,
        undefined
      );
      setLabList(labList2.labs ?? []);
    })();
  }, [searchParams]);

  return (
    <DefaultLayout>
      <Typography variant='h4'>Profile</Typography>
      <Box>
        {errorMessage && <Alert severity='error'>{errorMessage}</Alert>}
        {successMessage && <Alert severity='success'>{successMessage}</Alert>}
        <Typography variant='h6' marginBottom='10px'>
          研究室アンケートの変更
        </Typography>
        <Stack spacing={2} width='90%' display='flex' flexDirection='column'>
          <TextField
            defaultValue={wishLab ?? ''}
            select
            label='希望する研究室'
            onChange={(e) => setWishLab(e.target.value)}
          >
            <MenuItem value={''}>未選択</MenuItem>
            {labList.map((lab) => (
              <MenuItem value={lab.id} key={lab.id}>
                {lab.name}
              </MenuItem>
            ))}
          </TextField>
        </Stack>
        <Button
          variant='contained'
          sx={{ marginY: '10px' }}
          onClick={async () => {
            if (!wishLab) {
              setErrorMessage('研究室を選択してください');
              return;
            }
            try {
              const user = await updateUserLab(wishLab);
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
          更新する
        </Button>
      </Box>
    </DefaultLayout>
  );
}

export default Profile;
