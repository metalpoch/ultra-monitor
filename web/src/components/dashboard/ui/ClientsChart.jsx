import { useState, useEffect } from 'react'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  TimeScale
} from 'chart.js'
import { Line } from 'react-chartjs-2'
import 'chartjs-adapter-date-fns'
import dayjs from 'dayjs'
import { COLOR } from '../../../constants/colors'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  TimeScale
)

export default function ClientsChart({ data }) {
  const [chartData, setChartData] = useState(null)

  useEffect(() => {
    if (data && data.length > 0) {
      prepareChartData()
    }
  }, [data])

  const prepareChartData = () => {
    const labels = data.map(item => {
      const date = new Date(item.date)
      // Adjust for timezone offset to get correct local date
      return new Date(date.getTime() + date.getTimezoneOffset() * 60000)
    })

    const actives = data.map(item => item.actives)
    const provisionedOffline = data.map(item => item.provisioned_offline)
    const cutOff = data.map(item => item.cut_off)
    const inProgress = data.map(item => item.in_progress)

    setChartData({
      labels,
      datasets: [
        {
          label: 'Activos',
          data: labels.map((label, index) => ({
            x: label,
            y: actives[index]
          })),
          borderColor: COLOR[0],
          backgroundColor: COLOR[0],
          borderWidth: 2,
          fill: false,
          tension: 0.1,
          pointRadius: 0
        },
        {
          label: 'Activos offline',
          data: labels.map((label, index) => ({
            x: label,
            y: provisionedOffline[index]
          })),
          borderColor: COLOR[1],
          backgroundColor: COLOR[1],
          borderWidth: 2,
          fill: false,
          tension: 0.1,
          pointRadius: 0
        },
        {
          label: 'Cortados',
          data: labels.map((label, index) => ({
            x: label,
            y: cutOff[index]
          })),
          borderColor: COLOR[2],
          backgroundColor: COLOR[2],
          borderWidth: 2,
          fill: false,
          tension: 0.1,
          pointRadius: 0
        },
        {
          label: 'En proceso',
          data: labels.map((label, index) => ({
            x: label,
            y: inProgress[index]
          })),
          borderColor: COLOR[3],
          backgroundColor: COLOR[3],
          borderWidth: 2,
          fill: false,
          tension: 0.1,
          pointRadius: 0
        }
      ]
    })
  }

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    layout: {
      padding: {
        top: 10,
        right: 10,
        bottom: 10,
        left: 10
      }
    },
    plugins: {
      legend: {
        position: 'top',
        labels: {
          color: '#e2e8f0',
          font: {
            size: 12
          }
        }
      },
      title: {
        display: true,
        text: 'Crecimiento de Usuarios',
        color: '#e2e8f0',
        font: {
          size: 16,
          weight: 'bold'
        }
      },
      tooltip: {
        mode: 'index',
        intersect: false,
        backgroundColor: '#1a233a',
        titleColor: '#e0e6ed',
        bodyColor: '#e0e6ed',
        borderColor: '#2d3652'
      }
    },
    scales: {
      x: {
        type: 'time',
        time: {
          unit: 'day',
          displayFormats: {
            day: 'MMM dd'
          }
        },
        grid: {
          display: false
        },
        ticks: {
          color: '#94a3b8',
          maxRotation: 45
        }
      },
      y: {
        grid: {
          display: false
        },
        ticks: {
          color: '#94a3b8'
        }
      }
    },
    interaction: {
      mode: 'nearest',
      axis: 'x',
      intersect: false
    }
  }

  if (!data || data.length === 0) {
    return (
      <div className="h-96 flex items-center justify-center">
        <div className="text-gray-400 text-lg">No hay datos disponibles</div>
      </div>
    )
  }

  return (
    <div className="w-full h-full">
      <div className="h-96 w-full">
        {chartData && <Line data={chartData} options={chartOptions} />}
      </div>
    </div>
  )
}

