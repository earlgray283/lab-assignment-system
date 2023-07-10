import { Box, Divider, Grid, Stack, Typography, Tooltip } from '@mui/material';
import React, { useContext, useEffect, useState } from 'react';
import { fetchLabs } from '../../apis/labs';
import { Lab, LabList } from '../../apis/models/lab';
import { Doughnut } from 'react-chartjs-2';
import { Chart as ChartJS, ArcElement, Legend } from 'chart.js';
import CheckIcon from '@mui/icons-material/Check';
import CloseIcon from '@mui/icons-material/Close';
import { DisplayGpa } from '../util';
import { UserContext } from '../../App';
import { Link } from 'react-router-dom';
import { cmpLessThan } from '../../libs/util';
import { Notification } from '../Notification';

ChartJS.register(ArcElement, Legend);

function calcMinGpa(lab: Lab): number {
  let minUserGPA = undefined;
  for (const userGPA of lab.userGPAs) {
    if (minUserGPA === undefined) {
      minUserGPA = userGPA;
      continue;
    }
    if (minUserGPA.gpa > userGPA.gpa) {
      minUserGPA = userGPA;
    }
  }
  return minUserGPA?.gpa ?? -1;
}

function calcLabGpaAveg(lab: Lab): number {
  let sum = 0;
  for (const userGPA of lab.userGPAs) {
    sum += userGPA.gpa;
  }
  return lab.userGPAs.length === 0 ? -1 : sum / lab.userGPAs.length;
}

function isAssignable(lab: Lab, userGpa: number): boolean {
  const mingpa = calcMinGpa(lab);
  return lab.userGPAs.length < lab.capacity || cmpLessThan(mingpa, userGpa);
}

function LabCard(props: {
  year?: number;
  pushNotification: (n: Notification) => void;
}): JSX.Element {
  const [labList, setLabList] = useState<LabList | undefined>(undefined);
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
      setLabList({ labs: [] });
      return;
    }
    (async () => {
      const labList2 = await fetchLabs(props.year ?? new Date().getFullYear(), [
        wishLab,
      ]);
      setLabList(labList2);
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
        {labList.labs.map((lab, i) => {
          lab.userGPAs.sort((a, b) => b.gpa - a.gpa);
          const mingpa = calcMinGpa(lab);
          const gpaAveg = calcLabGpaAveg(lab);
          const labMag = (lab.userGPAs.length / lab.capacity) * 100;

          if (i == 0 && !isAssignable(lab, user.gpa)) {
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
            >
              <Typography variant='h5'>
                {lab.name}{' '}
                <Tooltip
                  title={
                    isAssignable(lab, user.gpa)
                      ? '配属可能です'
                      : '配属ができません'
                  }
                >
                  {isAssignable(lab, user.gpa) ? (
                    <CheckIcon fontSize='small' sx={{ color: 'green' }} />
                  ) : (
                    <CloseIcon fontSize='small' sx={{ color: 'red' }} />
                  )}
                </Tooltip>
              </Typography>
              <Divider />
              <Box display='flex'>
                <Stack marginTop='5px'>
                  <Box>競争率: {<span>{labMag}</span>}%</Box>
                  <Box>定員: {lab.capacity}人</Box>
                  <Box>志望者数: {lab.userGPAs.length}人</Box>
                  <Box>GPA(最小は上位{lab.capacity}名中)</Box>
                  <Box marginLeft='10px'>
                    {' '}
                    - 平均: <DisplayGpa gpa={gpaAveg} />
                  </Box>
                  <Box marginLeft='10px'>
                    {' '}
                    - 最大: <DisplayGpa gpa={lab.userGPAs.at(0)?.gpa ?? -1} />
                  </Box>
                  <Box marginLeft='10px'>
                    {' '}
                    - 最小: <DisplayGpa gpa={mingpa} />
                  </Box>
                </Stack>
              </Box>
              <Box display='flex' flexDirection='column' alignItems='center'>
                <Box width='75%'>
                  <Doughnut
                    data={{
                      labels: ['第1希望', '第2希望', '第3希望'],
                      datasets: [
                        {
                          data: [lab.userGPAs.length, 0, 0],
                          backgroundColor: [
                            '#AFD7F7CC',
                            '#84A4D4CC',
                            '#222E80CC',
                          ],
                        },
                      ],
                    }}
                  />
                </Box>
              </Box>
            </Box>
          );
        })}
      </Grid>
    </Box>
  );
}

export default LabCard;
