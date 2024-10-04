import { useEffect } from 'react';

interface LineProps {
  title: string;
  canvasID: string;
}

export default function Line({ title, canvasID }: LineProps) {
  const data = [
      { year: 2010, count: 10 },
      { year: 2011, count: 20 },
      { year: 2012, count: 15 },
  ];

  useEffect(() => {
      const canvas = document.getElementById(canvasID) as HTMLCanvasElement;

      if (canvas) {
          new Chart(canvas, {
              type: 'line',
              data: {
                  labels: data.map(row => row.year),
                  datasets: [{
                      label: 'Count',
                      data: data.map(row => row.count),
                      fill: false,
                      borderColor: 'rgb(75, 192, 192)',
                      tension: 0.1
                  }]
              },
              options: {
                  plugins: {
                      legend: {
                          position: 'bottom',
                          align: 'start',
                      }
                  }
              }
          });
      }
  }, []);

  return(
    <div className="w-1/4 min-w-fit h-fit px-12 py-10 bg-white rounded-xl shadow-2xl flex flex-col justify-center items-center">
      <h1 className="self-start text-2xl font-semibold text-gray-700">{title}</h1>
      <section className="w-72 h-fit">
          <canvas id={canvasID} className='w-full h-full'></canvas>
      </section>
    </div>
  );
}
