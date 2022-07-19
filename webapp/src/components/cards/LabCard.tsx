import { Box, Divider, Grid, Stack, Typography, Tooltip } from '@mui/material';
import React, { useContext, useEffect, useState } from 'react';
import { fetchLabs } from '../../apis/labs';
import { LabList } from '../../apis/models/lab';
import { Doughnut } from 'react-chartjs-2';
import { Chart as ChartJS, ArcElement, Legend } from 'chart.js';
import CheckIcon from '@mui/icons-material/Check';
import CloseIcon from '@mui/icons-material/Close';
import { DisplayGpa } from '../util';
import { UserContext } from '../../App';
import {
  NotificationsContext,
  NotificationsDispatchContext,
} from '../../pages/Dashboard';
import { Link } from 'react-router-dom';

ChartJS.register(ArcElement, Legend);

function LabCard(props: { labIds?: string[]; gpa: number }): JSX.Element {
  const [labList, setLabList] = useState<LabList | undefined>(undefined);
  const notifications = useContext(NotificationsContext);
  const setNotifications = useContext(NotificationsDispatchContext);
  const user = useContext(UserContext);
  useEffect(() => {
    (async () => {
      const labList2 = await fetchLabs(props.labIds, ['grade']);
      setLabList(labList2);
    })();
  }, []);

  if (!labList || user === undefined) {
    return <div>loading</div>;
  }
  if (!user) {
    return <div>unauthorized</div>;
  }
  if (!user.lab1) {
    console.log('check');
    const newNotifications = [...notifications];
    newNotifications.push({
      severity: 'error',
      message: (
        <div>
          研究室アンケートに回答されていません。
          <Link to='/profile'>研究室アンケートページ</Link>
          から回答を行って下さい。
        </div>
      ),
    });
    setNotifications([...newNotifications]);
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
        {labList.labs.map((lab) => {
          if (!lab.grades) {
            // unreachable
            return <div />;
          }
          const mingpa =
            lab.grades.gpas1.at(lab.capacity - 1) ??
            lab.grades.gpas1.at(lab.grades.gpas1.length - 1) ??
            -1;
          const gpaAveg =
            lab.grades.gpas1.length != 0
              ? lab.grades.gpas1.reduce((prev, cur) => prev + cur) /
                lab.grades.gpas1.length
              : -1;
          const labMag = (lab.firstChoice / lab.capacity) * 100;

          if (lab.grades.gpas1.length >= lab.capacity && mingpa >= user.gpa) {
            const newNotifications = [...notifications];
            newNotifications.push({
              severity: 'error',
              message: (
                <div>
                  第1希望の研究室({lab.name}
                  )への配属ができません。第一希望の研究室を変更してください。
                </div>
              ),
            });
            setNotifications([...newNotifications]);
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
                    lab.grades.gpas1.length < lab.capacity || mingpa < user.gpa
                      ? '配属可能です'
                      : '配属ができません'
                  }
                >
                  {lab.grades.gpas1.length < lab.capacity ||
                  mingpa < user.gpa ? (
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
                  <Box>志望者数: {lab.firstChoice}人</Box>
                  <Box>GPA(最小は上位{lab.capacity}名中)</Box>
                  <Box marginLeft='10px'>
                    {' '}
                    - 平均: <DisplayGpa gpa={gpaAveg} />
                  </Box>
                  <Box marginLeft='10px'>
                    {' '}
                    - 最大: <DisplayGpa gpa={lab.grades.gpas1.at(0) ?? -1} />
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
                          data: [
                            lab.firstChoice,
                            lab.secondChoice,
                            lab.thirdChoice,
                          ],
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
