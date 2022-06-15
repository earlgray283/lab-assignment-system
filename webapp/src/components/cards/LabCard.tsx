import { Box, Divider, Grid, Stack, Typography } from '@mui/material';
import React from 'react';
import { fetchLabs } from '../../apis/labs';
import { LabList } from '../../apis/models/lab';
import { sleep } from '../../lib/util';
import { Doughnut } from 'react-chartjs-2';
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js';
import CheckIcon from '@mui/icons-material/Check';
import CloseIcon from '@mui/icons-material/Close';

ChartJS.register(ArcElement, Tooltip, Legend);

let labList: LabList | undefined;
function useLabList(labIds: string[]): LabList {
  if (labList === undefined) {
    throw fetchLabs(labIds, ['grade'])
      .then((data) => (labList = data))
      .catch(() => sleep(2000));
  }
  return labList;
}

function LabCard(props: { labIds: string[]; gpa: number }): JSX.Element {
  console.log(props.labIds);

  const labList = useLabList(props.labIds);
  return (
    <Box padding='5px'>
      <Grid
        container
        direction={'row'}
        justifyContent='center'
        alignItems='center'
      >
        {labList.labs.map((lab, i) => {
          let rank = -1;
          if (lab.grades) {
            if (i == 0) {
              rank = lab.grades.gpas1.indexOf(props.gpa) + 1;
            } else if (i == 1) {
              rank = lab.grades.gpas2.indexOf(props.gpa) + 1;
            } else {
              rank = lab.grades.gpas3.indexOf(props.gpa) + 1;
            }
          }
          const labMag =
            ((lab.firstChoice + lab.secondChoice + lab.thirdChoice) /
              lab.capacity) *
            100;
          return (
            <Box
              key={lab.id}
              boxShadow={2}
              margin='10px'
              padding='10px'
              width='30%'
              minWidth='300px'
            >
              <Typography variant='h5'>
                {lab.name}{' '}
                {rank <= lab.capacity ? (
                  <CheckIcon fontSize='small' sx={{ color: 'green' }} />
                ) : (
                  <CloseIcon fontSize='small' sx={{ color: 'red' }} />
                )}
              </Typography>
              <Divider />
              <Box display='flex'>
                <Stack marginTop='5px'>
                  <Box>競争率: {<span>{labMag}</span>} %</Box>
                  <Box>定員: {lab.capacity}人</Box>
                  <Box>
                    志望者数:{' '}
                    {lab.firstChoice + lab.secondChoice + lab.thirdChoice}人
                  </Box>
                  <Box marginLeft='10px'> - 第1希望: {lab.firstChoice}人</Box>
                  <Box marginLeft='10px'> - 第2希望: {lab.secondChoice}人</Box>
                  <Box marginLeft='10px'> - 第3希望: {lab.thirdChoice}人</Box>
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
