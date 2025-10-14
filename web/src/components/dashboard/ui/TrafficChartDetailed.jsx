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

export default function TrafficChartDetailed({ data, dataType }) {
  const [chartData, setChartData] = useState(null)

  useEffect(() => {
    if (data && Object.keys(data).length > 0) {
      prepareChartData()
    }
  }, [data, dataType])

  const prepareChartData = () => {
    // Get all unique times from all datasets
    const allTimes = [
      ...new Set(
        Object.values(data)
          .flat()
          .map((d) => d.time)
      ),
    ].sort()

    const labels = allTimes.map(time => time)

    const datasets = Object.keys(data).map((name, index) => {
      const datasetData = labels.map((_, labelIndex) => {
        const time = allTimes[labelIndex]
        const point = data[name].find((d) => d.time === time)
        if (point) {
          if (dataType === "traffic") {
            return Math.max(point.bps_in || 0, point.bps_out || 0)
          } else {
            return Math.max(point.bytes_in || 0, point.bytes_out || 0)
          }
        }
        return null
      })

      return {
        label: name,
        data: labels.map((label, index) => ({
          x: label,
          y: datasetData[index]
        })),
        borderColor: COLOR[index % COLOR.length],
        backgroundColor: COLOR[index % COLOR.length],
        borderWidth: 2,
        fill: false,
        tension: 0.1,
        pointRadius: 0
      }
    })

    setChartData({
      labels,
      datasets
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
        text: dataType === 'traffic' ? 'Tráfico de Red Detallado' : 'Volumen de la Red Detallado',
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

  if (!data || Object.keys(data).length === 0) {
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

