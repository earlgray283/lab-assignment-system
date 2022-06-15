import { Box } from '@mui/material';
import React from 'react';
import { fetchLabs } from '../../apis/labs';
import { LabList } from '../../apis/models/lab';
import { sleep } from '../../lib/util';

let labList: LabList | undefined;
function useLabList(labIds: string[]): LabList {
  if (labList === undefined) {
    throw fetchLabs(labIds)
      .then((data) => (labList = data))
      .catch(() => sleep(2000));
  }
  return labList;
}

function LabCard(props: { labIds: string[] }): JSX.Element {
  console.log(props.labIds);
  const labList = useLabList(props.labIds);
  return (
    <Box boxShadow={1} padding='5px'>
      <Box display='flex' flexDirection='row' alignItems='center'>
        {labList.labs.map((lab) => {
          return <Box key={lab.id}>{lab.name}</Box>;
        })}
      </Box>
    </Box>
  );
}

export default LabCard;
