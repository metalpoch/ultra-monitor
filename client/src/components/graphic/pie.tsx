import { useEffect } from 'react';
import { Chart } from 'chart.js/auto';
import type { Measurement } from '../../models/measurement';

interface PieProps {
    title: string;
    canvasID: string;
    data?: Measurement[];
}

export default function Pie({ title, canvasID, data }: PieProps) {
        
    useEffect(() => {
        const canvas = document.getElementById(canvasID) as HTMLCanvasElement;

        if ((canvas) && (data)) {
            new Chart(canvas, {
                type: 'pie',
                data: {
                    labels: data.map((measurement: Measurement) => {
                        let currentDate = new Date(measurement.date);
                        return `${currentDate.getDate().toString().padStart(2, '0')}/${currentDate.getMonth().toString().padStart(2, '0')}/${currentDate.getFullYear()}`;
                    }),
                    datasets: [{
                        label: 'Count',
                        data: data.map((measurement: Measurement) => measurement.out_bps),
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

    return (
        <div className="w-1/4 min-w-fit h-fit px-12 py-10 bg-white rounded-xl shadow-2xl flex flex-col justify-center items-center">
            <h1 className="self-start text-2xl font-semibold text-gray-700">{title}</h1>
            <section className="w-72 h-fit">
                <canvas id={canvasID} className='w-full h-full'></canvas>
            </section>
        </div>
    )
}
