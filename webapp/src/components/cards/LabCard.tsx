import { Box, Divider, Grid, Stack, Typography, Tooltip } from '@mui/material';
import React, { useContext, useEffect, useState } from 'react';
import { fetchLabs } from '../../apis/labs';
import { Lab, UserGPA } from '../../apis/models/lab';
import { Chart as ChartJS, ArcElement, Legend } from 'chart.js';
import CheckIcon from '@mui/icons-material/Check';
import CloseIcon from '@mui/icons-material/Close';
import { DisplayGpa } from '../util';
import { UserContext } from '../../App';
import { Link } from 'react-router-dom';
import { cmpLessThan } from '../../libs/util';
import { Notification } from '../../pages/Dashboard';

ChartJS.register(ArcElement, Legend);

function LabCard(props: {
  year: number;
  pushNotification: (n: Notification) => void;
}): JSX.Element {
  const [labList, setLabList] = useState<(Lab | undefined)[] | undefined>(
    undefined,
  );
  const user = useContext(UserContext);

  useEffect(() => {
    if (!user) {
      return;
    }
    const wishLab = user.wishLab;
    if (!wishLab) {
      // MEMO: useEffect の中で親の state をいじらないと
      // setState-in-render が起きてしまう
      props.pushNotification({
        id: 'no-wish-lab',
        severity: 'error',
        message: (
          <div>
            研究室アンケートに回答されていません。
            <Link to='/profile'>研究室アンケートページ</Link>
            から回答を行って下さい。
          </div>
        ),
      });
      setLabList([]);
      return;
    }

    const labIDs = [wishLab];
    const watchLab1 = localStorage.getItem('watchLab1');
    if (watchLab1) {
      labIDs.push(watchLab1);
    }
    const watchLab2 = localStorage.getItem('watchLab2');
    if (watchLab2) {
      labIDs.push(watchLab2);
    }

    (async () => {
      const labList2 = await fetchLabs(props.year, labIDs);
      const labs2: (Lab | undefined)[] = [];
      for (const lab of labList2.labs) {
        labs2.push(lab);
      }
      for (let i = labList2.labs.length; i < 3; i++) {
        labs2.push(undefined);
      }
      setLabList([...labs2]);
    })();
  }, [user]);

  if (!user) {
    return <div>unauthorized</div>;
  }
  if (!labList || user === undefined) {
    return <div />;
  }

  return (
    <Box padding='5px'>
      <Grid
        container
        boxShadow={2}
        direction={'row'}
        justifyContent='center'
        alignItems='center'
      >
        {labList.map((lab, i) => {
          if (!lab) {
            return (
              <Box
                key={i}
                boxShadow={2}
                margin='20px 10px 20px 10px'
                padding='10px'
                width='30%'
                minWidth='300px'
                height='250px'
              ></Box>
            );
          }

          if (lab.userGPAs == null) {
            lab.userGPAs = [] as UserGPA[];
          }
          const userGPAsLength = lab.userGPAs.length;
          const averageGPA =
            userGPAsLength === 0
              ? 0
              : lab.userGPAs.map((u) => u.gpa).reduce((p, c) => p + c, 0) /
                userGPAsLength;

          lab.userGPAs.sort((a, b) => b.gpa - a.gpa);
          const borderLine = lab.userGPAs.at(lab.capacity - 1)?.gpa;
          const assignable =
            userGPAsLength < lab.capacity ||
            cmpLessThan(borderLine ?? 0, user.gpa);

          if (i == 0 && !assignable) {
            props.pushNotification({
              id: 'no-assignable-lab',
              severity: 'error',
              message: (
                <div>
                  希望の研究室({lab.name}
                  )への配属ができません。変更してください。
                </div>
              ),
            });
          }

          return (
            <Box
              key={lab.id}
              boxShadow={2}
              margin='20px 10px 20px 10px'
              padding='10px'
              width='30%'
              minWidth='300px'
              height='250px'
            >
              <Typography variant='h5'>
                {lab.name}{' '}
                <Tooltip
                  title={assignable ? '配属可能です' : '配属ができません'}
                >
                  {assignable ? (
                    <CheckIcon fontSize='small' sx={{ color: 'green' }} />
                  ) : (
                    <CloseIcon fontSize='small' sx={{ color: 'red' }} />
                  )}
                </Tooltip>
              </Typography>

              <Divider />

              <Box display='flex'>
                <Stack marginTop='5px'>
                  <Box>定員: {lab.capacity}人</Box>
                  <Box>志望者数: {userGPAsLength}人</Box>
                  <Box>
                    競争率:{' '}
                    {<span>{(userGPAsLength / lab.capacity) * 100}</span>}%
                  </Box>

                  <Box>GPA</Box>
                  <Box marginLeft='10px'>
                    {' '}
                    - 平均: <DisplayGpa gpa={averageGPA} />
                  </Box>

                  <Box marginLeft='10px'>
                    {' '}
                    - ボーダーライン
                    {borderLine ? <DisplayGpa gpa={borderLine} /> : 'BF'}
                  </Box>
                </Stack>
              </Box>
            </Box>
          );
        })}
      </Grid>
    </Box>
  );
}

export default LabCard;
