import { Chart } from 'chart.js/auto';
import { useEffect } from 'react';
import { LoadStatus } from '../../constant/loadStatus';
import { unit } from '../../utils/transform';
import type { Measurement } from '../../models/measurement';
import type { LoadingStateValue } from '../../types/loadingState';
import SpinnerBasicComponent from '../spinner/basic';
import React from 'react';

interface LineProps {
  loading: LoadingStateValue;
  title: string;
  canvasID: string;
  data: Measurement[];
}

export default function LineGraphComponent({ loading, title, canvasID, data }: LineProps) {

  const ViewWaitingData = () => {
    return(<h1 className="text-2xl font-semibold text-gray-300 text-center">Sin búsqueda</h1>);
  }

  const ViewWithoutData = () => {
    return(<h1 className="text-2xl font-semibold text-gray-300 text-center">No se encontró información para gráficar</h1>);
  }

  const ViewData = () => {
    return (
      <>
        <h1 className="self-start text-2xl font-semibold text-gray-700">{title}</h1>
        <canvas id={canvasID}></canvas>
      </>
    )
  }

  useEffect(() => {
    const canvas = document.getElementById(canvasID) as HTMLCanvasElement;
    if (canvas) {
      new Chart(canvas, {
        type: 'line',
        data: {
          labels: data.map((measurement: Measurement) => {
            let currentDate = new Date(measurement.date);
            return `${currentDate.getDate().toString().padStart(2, '0')}/${currentDate.getMonth().toString().padStart(2, '0')}/${currentDate.getFullYear()} ${currentDate.getHours().toString().padStart(2, '0')}:${currentDate.getMinutes().toString().padStart(2, '0')}`;
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
              tension: 0.1, 
              hidden: true
            },
          ]
        },
        options: {
          plugins: {
              legend: {
                  position: 'bottom',
                  align: 'center',
              }
          },
          scales: {
            x: { display: true },
            y: {
              ticks: {
                callback: (value) => {
                  if (typeof value === 'number') return unit(value);
                },
              },
            },
          },
        }
      });
    }
  }, [data]);

  return (
    <div className="w-full lg:w-5/6 max-w-full min-h-96 px-4 py-6 bg-white rounded-xl flex flex-col justify-center items-center gap-4">
      {loading === LoadStatus.EMPTY && ViewWaitingData()}
      {loading === LoadStatus.LOADING && <SpinnerBasicComponent />}
      {loading === LoadStatus.LOADED && data && data.length > 0 && ViewData()}
      {loading === LoadStatus.LOADED && data && data.length <= 0 && ViewWithoutData()}
      {loading === LoadStatus.LOADED && !data && ViewWithoutData()}
    </div>
  );
}
