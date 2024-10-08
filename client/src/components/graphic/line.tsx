import { useEffect } from 'react';
import { Chart } from 'chart.js/auto';
import type { Measurement } from '../../models/measurement';

interface LineProps {
  title: string;
  canvasID: string;
  data?: Measurement[];
}

export default function Line({ title, canvasID, data }: LineProps) {
  useEffect(() => {
      const canvas = document.getElementById(canvasID) as HTMLCanvasElement;
      if ((canvas) && (data)) {
        new Chart(canvas, {
          type: 'line',
          data: {
            labels: data.map((measurement: Measurement) => {
              let currentDate = new Date(measurement.date);
              return `${currentDate.getDate().toString().padStart(2, '0')}/${currentDate.getMonth().toString().padStart(2, '0')}/${currentDate.getFullYear()}`;
            }),
            datasets: [
              {
                label: 'In',
                data: data.map((measurement: Measurement) => measurement.in_bps),
                fill: false,
                borderColor: 'rgb(75, 192, 25)',
                tension: 0.1
              },
              {
                label: 'Out',
                data: data.map((measurement: Measurement) => measurement.out_bps),
                fill: false,
                borderColor: 'rgb(75, 192, 192)',
                tension: 0.1
              },
              {
                label: 'Bandwith',
                data: data.map((measurement: Measurement) => measurement.bandwidth_bps),
                fill: false,
                borderColor: 'rgb(75, 25, 200)',
                tension: 0.1
              },
            ]
          },
          options: {
            plugins: {
                legend: {
                    position: 'bottom',
                    align: 'center',
                }
            }
          }
      });
    }
  }, []);

  return(
    <div className="w-full min-w-fit h-fit px-12 py-10 bg-white rounded-xl shadow-2xl flex flex-col justify-center items-center gap-4">
      <h1 className="self-start text-2xl font-semibold text-gray-700">{title}</h1>
      <section className="w-full h-fit">
          <canvas id={canvasID} className='w-full h-full'></canvas>
      </section>
    </div>
  );
}
