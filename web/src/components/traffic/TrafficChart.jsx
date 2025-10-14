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
import { es } from 'date-fns/locale'
import { COLOR } from '../../constants/colors'

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

export default function TrafficChart({ data, dataType }) {
  const [chartData, setChartData] = useState(null)

  useEffect(() => {
    if (data && data.length > 0) {
      prepareChartData()
    }
  }, [data, dataType])

  const prepareChartData = () => {
    // Use the original date strings - Chart.js will handle timezone conversion
    const labels = data.map(item => new Date(item.time))

    const bpsIn = data.map(item => item.bps_in || 0)
    const bpsOut = data.map(item => item.bps_out || 0)
    const bytesIn = data.map(item => item.bytes_in || 0)
    const bytesOut = data.map(item => item.bytes_out || 0)

    setChartData({
      labels,
      datasets: [
        {
          label: 'Entrante',
          data: labels.map((label, index) => ({
            x: label,
            y: dataType === 'traffic' ? bpsIn[index] : bytesIn[index]
          })),
          borderColor: dataType === 'traffic' ? COLOR[9] : COLOR[1],
          backgroundColor: dataType === 'traffic' ? COLOR[9] : COLOR[1],
          borderWidth: 2,
          fill: false,
          tension: 0.1,
          pointRadius: 0
        },
        {
          label: 'Saliente',
          data: labels.map((label, index) => ({
            x: label,
            y: dataType === 'traffic' ? bpsOut[index] : bytesOut[index]
          })),
          borderColor: dataType === 'traffic' ? COLOR[5] : COLOR[3],
          backgroundColor: dataType === 'traffic' ? COLOR[5] : COLOR[3],
          borderWidth: 2,
          fill: false,
          tension: 0.1,
          pointRadius: 0
        }
      ]
    })
  }

  const formatBytes = (bytes) => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  const formatBps = (bps) => {
    if (bps === 0) return '0 bps'
    const k = 1000
    const sizes = ['bps', 'Kbps', 'Mbps', 'Gbps', 'Tbps']
    const i = Math.floor(Math.log(bps) / Math.log(k))
    return parseFloat((bps / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
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
        text: dataType === 'traffic' ? 'Tráfico de Red' : 'Volumen de la Red',
        color: '#e2e8f0',
        font: {
          size: 16,
          weight: 'bold'
        }
      },
      tooltip: {
        mode: 'index',
        intersect: false,
        callbacks: {
          label: function(context) {
            let label = context.dataset.label || ''
            if (label) {
              label += ': '
            }
            if (context.parsed.y !== null) {
              label += dataType === 'traffic'
                ? formatBps(context.parsed.y)
                : formatBytes(context.parsed.y)
            }
            return label
          }
        },
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
          },
          // Use local timezone instead of UTC
          tooltipFormat: 'PPpp',
        },
        adapters: {
          date: {
            locale: es // Use Spanish locale for date formatting
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
          color: '#94a3b8',
          callback: function(value) {
            return dataType === 'traffic'
              ? formatBps(value)
              : formatBytes(value)
          }
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

