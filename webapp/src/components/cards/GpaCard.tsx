import { Box } from '@mui/material';
import {
  BarElement,
  CategoryScale,
  Chart as ChartJS,
  Legend,
  LinearScale,
  Title,
  Tooltip,
} from 'chart.js';
import React, { useContext, useEffect, useState } from 'react';
import { Bar } from 'react-chartjs-2';
import { fetchGrades } from '../../apis/grade';
import { sleep } from '../../libs/util';
import { DisplayGpa } from '../util';
import { UserContext } from '../../App';

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
);

const labels = [
  '0.0-0.5',
  '0.5-1.0',
  '1.0-1.5',
  '1.5-2.0',
  '2.0-2.5',
  '2.5-3.0',
  '3.0-3.5',
  '3.5-4.0',
];

const options = {
  responsive: true,
  plugins: {
    legend: {
      display: false,
      //position: 'top' as const,
    },
  },
};

let gpas: number[] | undefined;
function useGpas(year: number): number[] {
  if (gpas === undefined) {
    throw fetchGrades(year)
      .then((data) => (gpas = data))
      .catch(() => sleep(2000));
  }
  return gpas;
}

function GpaCard(props: { year: number }): JSX.Element {
  const user = useContext(UserContext);
  const [gpaClasses, setGpaClasses] = useState<number[]>([]);
  const gpas = useGpas(props.year);
  useEffect(() => {
    const list = new Array<number>(0, 0, 0, 0, 0, 0, 0, 0);
    for (const gpa of gpas) {
      if (gpa <= 0.5) list[0]++;
      else if (gpa <= 1.0) list[1]++;
      else if (gpa <= 1.5) list[2]++;
      else if (gpa <= 2.0) list[3]++;
      else if (gpa <= 2.5) list[4]++;
      else if (gpa <= 3.0) list[5]++;
      else if (gpa <= 3.5) list[6]++;
      else list[7]++;
    }
    setGpaClasses([...list]);
  }, [gpas]);
  const data = {
    labels,
    datasets: [
      {
        data: gpaClasses,
        backgroundColor: 'rgba(153, 183, 220, 0.6)',
      },
    ],
  };

  if (!user) {
    return <div />;
  }

  return (
    <Box boxShadow={2} padding='5px'>
      <p>
        あなたの GPA は <DisplayGpa gpa={user.gpa} /> です
      </p>
      <Box display='flex' flexDirection='column' alignItems='center'>
        <Box width='50%' minWidth='300px'>
          <Bar data={data} options={options} />
        </Box>
      </Box>
    </Box>
  );
}

export default GpaCard;
