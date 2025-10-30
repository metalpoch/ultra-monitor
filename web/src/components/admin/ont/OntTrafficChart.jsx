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

export default function OntTrafficChart() {
  const [selectedOnt, setSelectedOnt] = useState(null)
  const [trafficData, setTrafficData] = useState([])
  const [loading, setLoading] = useState(false)
  const [dataType, setDataType] = useState('traffic') // 'traffic', 'volume', 'signal', or 'temperature'
  const [timeRange, setTimeRange] = useState('week') // 'day', 'week', 'month'

  const BASE_URL = `${import.meta.env.PUBLIC_URL || ""}/api/traffic`
  const TOKEN = sessionStorage.getItem("access_token")?.replace("Bearer ", "") || ""

  const headers = {
    headers: {
      Authorization: `Bearer ${TOKEN}`
    }
  }

  // Listen for ONT selection
  useEffect(() => {
    const handleOntSelected = (event) => {
      setSelectedOnt(event.detail)
    }

    window.addEventListener('ontSelected', handleOntSelected)
    return () => {
      window.removeEventListener('ontSelected', handleOntSelected)
    }
  }, [])

  // Fetch traffic data when ONT is selected
  useEffect(() => {
    if (selectedOnt) {
      fetchTrafficData()
    } else {
      setTrafficData([])
    }
  }, [selectedOnt, timeRange])

  const getDateRange = () => {
    const now = new Date()
    const initDate = new Date()

    switch (timeRange) {
      case 'day':
        initDate.setDate(now.getDate() - 1)
        break
      case 'week':
        initDate.setDate(now.getDate() - 7)
        break
      case 'month':
        initDate.setMonth(now.getMonth() - 1)
        break
      default:
        initDate.setDate(now.getDate() - 7)
    }

    return {
      initDate: initDate.toISOString(),
      finalDate: now.toISOString()
    }
  }

  const fetchTrafficData = async () => {
    if (!selectedOnt) return

    setLoading(true)
    try {
      const { initDate, finalDate } = getDateRange()
      const response = await fetch(
        `${BASE_URL}/ont/${selectedOnt.id}?initDate=${initDate}&finalDate=${finalDate}`,
        headers
      )

      if (response.status === 401 || response.status === 403) {
        sessionStorage.removeItem("access_token")
        window.location.href = "/auth/login"
        return
      }

      if (!response.ok) {
        throw new Error(`Error: ${response.status}`)
      }

      const data = await response.json()
      setTrafficData(data)
    } catch (err) {
      console.error('Error fetching traffic data:', err)
      setTrafficData([])
    } finally {
      setLoading(false)
    }
  }

  const prepareChartData = () => {
    if (!trafficData || trafficData.length === 0) {
      return null
    }

    // Use the original date strings - Chart.js will handle timezone conversion
    const labels = trafficData.map(item => new Date(item.time))

    const bpsIn = trafficData.map(item => item.bps_in || 0)
    const bpsOut = trafficData.map(item => item.bps_out || 0)
    const bytesIn = trafficData.map(item => item.bytes_in || 0)
    const bytesOut = trafficData.map(item => item.bytes_out || 0)
    const rx = trafficData.map(item => item.rx || 0)
    const tx = trafficData.map(item => item.tx || 0)
    const temperature = trafficData.map(item => item.temperature || 0)

    switch (dataType) {
      case 'traffic':
        return {
          labels,
          datasets: [
            {
              label: 'Entrante',
              data: labels.map((label, index) => ({
                x: label,
                y: bpsIn[index]
              })),
              borderColor: COLOR[9],
              backgroundColor: COLOR[9],
              borderWidth: 2,
              fill: false,
              tension: 0.1,
              pointRadius: 0
            },
            {
              label: 'Saliente',
              data: labels.map((label, index) => ({
                x: label,
                y: bpsOut[index]
              })),
              borderColor: COLOR[5],
              backgroundColor: COLOR[5],
              borderWidth: 2,
              fill: false,
              tension: 0.1,
              pointRadius: 0
            }
          ]
        }

      case 'volume':
        return {
          labels,
          datasets: [
            {
              label: 'Entrante',
              data: labels.map((label, index) => ({
                x: label,
                y: bytesIn[index]
              })),
              borderColor: COLOR[1],
              backgroundColor: COLOR[1],
              borderWidth: 2,
              fill: false,
              tension: 0.1,
              pointRadius: 0
            },
            {
              label: 'Saliente',
              data: labels.map((label, index) => ({
                x: label,
                y: bytesOut[index]
              })),
              borderColor: COLOR[3],
              backgroundColor: COLOR[3],
              borderWidth: 2,
              fill: false,
              tension: 0.1,
              pointRadius: 0
            }
          ]
        }

      case 'signal':
        return {
          labels,
          datasets: [
            {
              label: 'RX (dBm)',
              data: labels.map((label, index) => ({
                x: label,
                y: rx[index]
              })),
              borderColor: COLOR[2],
              backgroundColor: COLOR[2],
              borderWidth: 2,
              fill: false,
              tension: 0.1,
              pointRadius: 0
            },
            {
              label: 'TX (dBm)',
              data: labels.map((label, index) => ({
                x: label,
                y: tx[index]
              })),
              borderColor: COLOR[6],
              backgroundColor: COLOR[6],
              borderWidth: 2,
              fill: false,
              tension: 0.1,
              pointRadius: 0
            }
          ]
        }

      case 'temperature':
        return {
          labels,
          datasets: [
            {
              label: 'Temperatura (°C)',
              data: labels.map((label, index) => ({
                x: label,
                y: temperature[index]
              })),
              borderColor: COLOR[8],
              backgroundColor: COLOR[8],
              borderWidth: 2,
              fill: false,
              tension: 0.1,
              pointRadius: 0
            }
          ]
        }

      default:
        return null
    }
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
        text: selectedOnt
          ? `Tráfico de ONT: ${selectedOnt.serial} - ${selectedOnt.despt}`
          : 'Seleccione una ONT para ver el tráfico',
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
              if (dataType === 'traffic') {
                label += formatBps(context.parsed.y)
              } else if (dataType === 'volume') {
                label += formatBytes(context.parsed.y)
              } else if (dataType === 'signal') {
                label += context.parsed.y.toFixed(2) + ' dBm'
              } else if (dataType === 'temperature') {
                label += context.parsed.y.toFixed(1) + ' °C'
              }
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
          tooltipFormat: 'PPpp',
        },
        adapters: {
          date: {
            locale: es
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
            if (dataType === 'traffic') {
              return formatBps(value)
            } else if (dataType === 'volume') {
              return formatBytes(value)
            } else if (dataType === 'signal') {
              return value.toFixed(2) + ' dBm'
            } else if (dataType === 'temperature') {
              return value.toFixed(1) + ' °C'
            }
            return value
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

  const chartData = prepareChartData()

  if (!selectedOnt) {
    return (
      <div className="h-96 flex flex-col items-center justify-center text-gray-400">
        <div className="text-lg mb-2">Seleccione una ONT</div>
        <div className="text-sm">Haga clic en una ONT de la lista para ver su tráfico</div>
      </div>
    )
  }

  if (loading) {
    return (
      <div className="h-96 flex items-center justify-center">
        <div className="text-gray-400 text-lg">Cargando datos de tráfico...</div>
      </div>
    )
  }

  return (
    <div className="w-full h-full flex flex-col">
      {/* Controls */}
      <div className="flex flex-wrap gap-4 mb-4">
        <div className="flex items-center gap-2">
          <label className="text-gray-200 text-sm">Tipo de datos:</label>
          <select
            value={dataType}
            onChange={(e) => setDataType(e.target.value)}
            className="border border-[hsl(217,33%,20%)] bg-[#0f1729] text-gray-200 p-1 rounded text-sm"
          >
            <option value="traffic">Tráfico (bps)</option>
            <option value="volume">Volumen (bytes)</option>
            <option value="signal">Señal (RX/TX)</option>
            <option value="temperature">Temperatura</option>
          </select>
        </div>

        <div className="flex items-center gap-2">
          <label className="text-gray-200 text-sm">Rango de tiempo:</label>
          <select
            value={timeRange}
            onChange={(e) => setTimeRange(e.target.value)}
            className="border border-[hsl(217,33%,20%)] bg-[#0f1729] text-gray-200 p-1 rounded text-sm"
          >
            <option value="day">Último día</option>
            <option value="week">Última semana</option>
            <option value="month">Último mes</option>
          </select>
        </div>
      </div>

      {/* Chart */}
      <div className="flex-1">
        {!chartData ? (
          <div className="h-96 flex items-center justify-center">
            <div className="text-gray-400 text-lg">No hay datos de tráfico disponibles</div>
          </div>
        ) : (
          <div className="h-96 w-full">
            <Line data={chartData} options={chartOptions} />
          </div>
        )}
      </div>
    </div>
  )
}

