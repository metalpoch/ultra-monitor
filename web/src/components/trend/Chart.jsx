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

const BASE_URL_TRAFFIC = `${import.meta.env.PUBLIC_URL || ""}/api/traffic`
const TOKEN = sessionStorage.getItem("access_token")?.replace("Bearer ", "") || ""

const Chart = () => {
  const [trendData, setTrendData] = useState(null)
  const [historicalData, setHistoricalData] = useState(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)

  // Load national trend by default on component mount
  useEffect(() => {
    loadDefaultTrend()
  }, [])

  const loadDefaultTrend = async () => {
    setLoading(true)
    setError(null)

    try {
      const defaultDateRange = {
        initDate: dayjs().subtract(30, 'day').toISOString(),
        finalDate: dayjs().toISOString()
      }

      const params = new URLSearchParams({
        futureDays: '7',
        confidence: '0.95',
        initDate: defaultDateRange.initDate,
        finalDate: defaultDateRange.finalDate
      })

      // Fetch trend data
      const trendResponse = await fetch(`${BASE_URL_TRAFFIC}/trend/national?${params}`, {
        headers: { Authorization: `Bearer ${TOKEN}` }
      })

      if (!trendResponse.ok) {
        throw new Error(`Error: ${trendResponse.status} ${trendResponse.statusText}`)
      }

      const trendData = await trendResponse.json()
      setTrendData(trendData)

      // Fetch historical traffic data
      const historicalParams = new URLSearchParams({
        initDate: defaultDateRange.initDate,
        finalDate: defaultDateRange.finalDate
      })

      const historicalResponse = await fetch(`${BASE_URL_TRAFFIC}/total?${historicalParams}`, {
        headers: { Authorization: `Bearer ${TOKEN}` }
      })

      if (!historicalResponse.ok) {
        throw new Error(`Error fetching historical data: ${historicalResponse.status} ${historicalResponse.statusText}`)
      }

      const historicalData = await historicalResponse.json()
      setHistoricalData(historicalData)

      // Dispatch event to update metrics component
      const event = new CustomEvent('trendDataUpdate', { detail: trendData })
      window.dispatchEvent(event)
    } catch (err) {
      setError(err.message)
      console.error('Error fetching data:', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    const handleTrendFormSubmit = async (event) => {
      const {
        futureDays,
        confidence,
        initDate,
        finalDate,
        selectedLevel,
        selectedRegion,
        selectedState,
        selectedOlt
      } = event.detail

      setLoading(true)
      setError(null)

      try {
        let endpoint = '/api/traffic/trend/national'
        let params = new URLSearchParams({
          futureDays: futureDays.toString(),
          confidence: confidence.toString(),
          initDate: initDate,
          finalDate: finalDate
        })

        // Determine endpoint based on selected level
        switch (selectedLevel) {
          case 'national':
            endpoint = `${BASE_URL_TRAFFIC}/trend/national`
            break
          case 'regional':
            endpoint = `${BASE_URL_TRAFFIC}/trend/region/${encodeURIComponent(selectedRegion)}`
            break
          case 'state':
            endpoint = `${BASE_URL_TRAFFIC}/trend/state/${encodeURIComponent(selectedState)}`
            break
          case 'olt':
            endpoint = `${BASE_URL_TRAFFIC}/trend/olt/${encodeURIComponent(selectedOlt)}`
            break
          default:
            endpoint = `${BASE_URL_TRAFFIC}/trend/national`
        }

              // Fetch trend data
              const trendResponse = await fetch(`${endpoint}?${params}`, {
          headers: { Authorization: `Bearer ${TOKEN}` }
        })

              if (!trendResponse.ok) {
                throw new Error(`Error: ${trendResponse.status} ${trendResponse.statusText}`)
        }

              const trendData = await trendResponse.json()
              setTrendData(trendData)
        
              // Fetch historical traffic data based on selected level
              const historicalParams = new URLSearchParams({
                initDate: initDate,
                finalDate: finalDate
              })
        
              let historicalEndpoint = '/api/traffic/total'
              switch (selectedLevel) {
                case 'national':
                  historicalEndpoint = '/api/traffic/total'
                  break
                case 'regional':
                  historicalEndpoint = `/api/traffic/region/${encodeURIComponent(selectedRegion)}`
                  break
                case 'state':
                  historicalEndpoint = `/api/traffic/state/${encodeURIComponent(selectedState)}`
                  break
                case 'olt':
                  historicalEndpoint = `/api/traffic/instances?ip=${encodeURIComponent(selectedOlt)}`
                  break
                default:
                  historicalEndpoint = '/api/traffic/total'
              }
        
              const historicalResponse = await fetch(`${BASE_URL_TRAFFIC}${historicalEndpoint.replace('/api/traffic', '')}?${historicalParams}`, {
                headers: { Authorization: `Bearer ${TOKEN}` }
              })

              if (!historicalResponse.ok) {
                setHistoricalData(null)
              } else {
                const historicalData = await historicalResponse.json()
                setHistoricalData(historicalData)
              }

        // Dispatch event to update metrics component
              const event = new CustomEvent('trendDataUpdate', { detail: trendData })
        window.dispatchEvent(event)
      } catch (err) {
        setError(err.message)
        console.error('Error fetching trend data:', err)
      } finally {
        setLoading(false)
      }
    }

    window.addEventListener('trendFormSubmit', handleTrendFormSubmit)

    return () => {
      window.removeEventListener('trendFormSubmit', handleTrendFormSubmit)
    }
  }, [])

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

  const getChartData = () => {
    if (!trendData) return null


    // Prepare historical data - handle timezone issues
    const historicalLabels = historicalData ? historicalData.map(h => {
      const date = new Date(h.time)
      // Adjust for timezone offset to get correct local date
      return new Date(date.getTime() + date.getTimezoneOffset() * 60000)
    }) : []
    const historicalBpsIn = historicalData ? historicalData.map(h => h.bps_in) : []
    const historicalBpsOut = historicalData ? historicalData.map(h => h.bps_out) : []

    // Get the last historical date to adjust prediction dates
    const lastHistoricalDate = historicalData && historicalData.length > 0
      ? new Date(Math.max(...historicalData.map(h => {
          const date = new Date(h.time)
          return new Date(date.getTime() + date.getTimezoneOffset() * 60000)
        })))
      : null


    // Adjust prediction dates to start from the day after the last historical data
    const adjustedPredictionLabels = trendData.predictions.map((p, index) => {
      if (lastHistoricalDate) {
        // Start predictions from the day after the last historical data
        const predictionStartDate = new Date(lastHistoricalDate)
        predictionStartDate.setDate(predictionStartDate.getDate() + 1 + index)
        return predictionStartDate
      }
      return new Date(p.date)
    })

    const predictedData = trendData.predictions.map(p => p.predicted_bps)
    const lowerBounds = trendData.predictions.map(p => p.lower_bound || null)
    const upperBounds = trendData.predictions.map(p => p.upper_bound || null)

    return {
      labels: [...historicalLabels, ...adjustedPredictionLabels],
      datasets: [
        // Historical BPS In
        ...(historicalData ? [{
          label: 'BPS Entrada (Histórico)',
          data: historicalLabels.map((label, index) => ({
            x: label,
            y: historicalBpsIn[index]
          })),
          borderColor: 'rgb(34, 197, 94)',
          backgroundColor: 'rgba(34, 197, 94, 0.1)',
          borderWidth: 2,
          fill: false,
          tension: 0.1,
          pointRadius: 0
        }] : []),

        // Historical BPS Out
        ...(historicalData ? [{
          label: 'BPS Salida (Histórico)',
          data: historicalLabels.map((label, index) => ({
            x: label,
            y: historicalBpsOut[index]
          })),
          borderColor: 'rgb(239, 68, 68)',
          backgroundColor: 'rgba(239, 68, 68, 0.1)',
          borderWidth: 2,
          fill: false,
          tension: 0.1,
          pointRadius: 0
        }] : []),

        // Historical Total BPS
        ...(historicalData ? [{
          label: 'Total BPS (Histórico)',
          data: historicalLabels.map((label, index) => ({
            x: label,
            y: historicalBpsIn[index] + historicalBpsOut[index]
          })),
          borderColor: 'rgb(59, 130, 246)',
          backgroundColor: 'rgba(59, 130, 246, 0.1)',
          borderWidth: 2,
          fill: false,
          tension: 0.1,
          pointRadius: 0
        }] : []),

        // Prediction Total BPS
        {
          label: 'Total BPS (Predicción)',
          data: adjustedPredictionLabels.map((label, index) => ({
            x: label,
            y: predictedData[index]
          })),
          borderColor: 'rgb(147, 51, 234)',
          backgroundColor: 'rgba(147, 51, 234, 0.1)',
          borderWidth: 2,
          borderDash: [5, 5],
          fill: false,
          tension: 0.1,
          pointRadius: 0
        },

        // Confidence bounds (only for prediction period)
        ...(lowerBounds[0] !== null && upperBounds[0] !== null ? [
          {
            label: 'Límite Inferior',
            data: adjustedPredictionLabels.map((label, index) => ({
              x: label,
              y: lowerBounds[index]
            })),
            borderColor: 'rgb(156, 163, 175)', // gray-400
            backgroundColor: 'rgba(156, 163, 175, 0.1)',
            borderWidth: 1,
            borderDash: [5, 5],
            fill: false,
            pointRadius: 0
          },
          {
            label: 'Límite Superior',
            data: adjustedPredictionLabels.map((label, index) => ({
              x: label,
              y: upperBounds[index]
            })),
            borderColor: 'rgb(156, 163, 175)', // gray-400
            backgroundColor: 'rgba(156, 163, 175, 0.1)',
            borderWidth: 1,
            borderDash: [5, 5],
            fill: false,
            pointRadius: 0
          }
        ] : [])
      ]
    }
  }

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
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
        text: 'Predicción de Tendencia de Tráfico',
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
              label += formatBps(context.parsed.y)
            }
            return label
          }
        }
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
          color: '#94a3b8'
        }
      },
      y: {
        grid: {
          display: false
        },
        ticks: {
          color: '#94a3b8',
          callback: function(value) {
            return formatBps(value)
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

  if (loading) {
    return (
      <div className="h-96 flex items-center justify-center">
        <div className="text-white text-lg">Cargando datos de tendencia...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="h-96 flex items-center justify-center">
        <div className="text-red-400 text-lg">Error: {error}</div>
      </div>
    )
  }

  if (!trendData) {
    return (
      <div className="h-96 flex items-center justify-center">
        <div className="text-gray-400 text-lg">
          Selecciona los parámetros y haz clic en "Generar Tendencia" para ver la predicción
        </div>
      </div>
    )
  }

  return (
    <div className="w-full">
      <div className="h-96">
        <Line data={getChartData()} options={chartOptions} />
      </div>
    </div>
  )
}

export default Chart

