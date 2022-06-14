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
import React from 'react';
import { Bar } from 'react-chartjs-2';

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend
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

interface Props {
  data: number[];
  gpa: number;
}

function GpaCard(props: Props): JSX.Element {
  const data = {
    labels,
    datasets: [
      {
        data: props.data,
        backgroundColor: 'rgba(153, 183, 220, 0.6)',
      },
    ],
  };

  return (
    <Box boxShadow={1} padding='5px'>
      <p>あなたの GPA は 4.0 です</p>
      <Box display='flex' flexDirection='column' alignItems='center'>
        <Box width='50%' minWidth='300px'>
          <Bar data={data} options={options} />
        </Box>
      </Box>
    </Box>
  );
}

export default GpaCard;
