import { useState, useEffect } from 'react'

const Metrics = () => {
  const [trendData, setTrendData] = useState(null)

  useEffect(() => {
    const handleTrendFormSubmit = (event) => {
      // This will be populated when the chart fetches data
      // We'll update this component when the chart data is available
    }

    const handleTrendDataUpdate = (event) => {
      setTrendData(event.detail)
    }

    window.addEventListener('trendFormSubmit', handleTrendFormSubmit)
    window.addEventListener('trendDataUpdate', handleTrendDataUpdate)

    return () => {
      window.removeEventListener('trendFormSubmit', handleTrendFormSubmit)
      window.removeEventListener('trendDataUpdate', handleTrendDataUpdate)
    }
  }, [])

  const formatNumber = (num) => {
    if (num === null || num === undefined || isNaN(num)) return 'N/A'
    return num.toFixed(4)
  }

  const formatBps = (bps) => {
    if (bps === null || bps === undefined || isNaN(bps) || !isFinite(bps)) {
      return 'N/A'
    }
    if (bps === 0) return '0 bps'

    const k = 1000
    const sizes = ['bps', 'Kbps', 'Mbps', 'Gbps', 'Tbps']

    try {
      // Handle negative values by using absolute value for calculation
      const absBps = Math.abs(bps)
      const i = Math.floor(Math.log(absBps) / Math.log(k))

      // Ensure index is within bounds
      const index = Math.min(Math.max(i, 0), sizes.length - 1)
      const value = absBps / Math.pow(k, index)

      const sign = bps < 0 ? '-' : ''
      return sign + parseFloat(value.toFixed(2)) + ' ' + sizes[index]
    } catch (error) {
      return 'N/A'
    }
  }

  const getTrendDirection = (metrics) => {
    if (!metrics) return 'neutral'
    if (metrics.is_increasing) return 'increasing'
    if (metrics.is_decreasing) return 'decreasing'
    return 'stable'
  }

  const getTrendIcon = (direction) => {
    switch (direction) {
      case 'increasing':
        return (
          <svg className="w-5 h-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 10l7-7m0 0l7 7m-7-7v18" />
          </svg>
        )
      case 'decreasing':
        return (
          <svg className="w-5 h-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 14l-7 7m0 0l-7-7m7 7V3" />
          </svg>
        )
      default:
        return (
          <svg className="w-5 h-5 text-yellow-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 12h14" />
          </svg>
        )
    }
  }

  const getTrendColor = (direction) => {
    switch (direction) {
      case 'increasing':
        return 'text-green-400'
      case 'decreasing':
        return 'text-red-400'
      default:
        return 'text-yellow-400'
    }
  }

  const getTrendText = (direction) => {
    switch (direction) {
      case 'increasing':
        return 'Tendencia Creciente'
      case 'decreasing':
        return 'Tendencia Decreciente'
      default:
        return 'Tendencia Estable'
    }
  }

  const getRSquaredInterpretation = (rSquared) => {
    if (rSquared === null || rSquared === undefined || isNaN(rSquared)) return 'Datos insuficientes'
    if (rSquared >= 0.9) return 'Excelente ajuste'
    if (rSquared >= 0.7) return 'Buen ajuste'
    if (rSquared >= 0.5) return 'Ajuste moderado'
    if (rSquared >= 0.3) return 'Ajuste débil'
    if (rSquared >= 0.1) return 'Ajuste muy débil'
    return 'Sin correlación significativa'
  }

  if (!trendData) {
    return (
      <div className="w-full p-6 bg-[#1e293b] border border-[#334155] rounded-lg">
        <h3 className="text-lg font-semibold text-white mb-4">Métricas de Tendencia</h3>
        <div className="text-gray-400 text-center py-8">
          Los métricos se mostrarán aquí después de generar una tendencia
        </div>
      </div>
    )
  }

  const { metrics } = trendData
  const trendDirection = getTrendDirection(metrics)

  return (
    <div className="w-full p-6 bg-[#1e293b] border border-[#334155] rounded-lg">
      <h3 className="text-lg font-semibold text-white mb-4">Métricas de Tendencia</h3>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Trend Direction */}
        <div className="bg-[#121b31] p-4 rounded-lg border border-[#334155]">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm text-gray-400">Dirección</span>
            {getTrendIcon(trendDirection)}
          </div>
          <div className={`text-lg font-semibold ${getTrendColor(trendDirection)}`}>
            {getTrendText(trendDirection)}
          </div>
        </div>

        {/* Slope */}
        <div className="bg-[#121b31] p-4 rounded-lg border border-[#334155]">
          <div className="text-sm text-gray-400 mb-2">Pendiente</div>
          <div className="text-lg font-semibold text-white">
            {formatBps(metrics.slope)}
          </div>
          <div className="text-xs text-gray-500 mt-1">
            Cambio por día
          </div>
        </div>

        {/* R-Squared */}
        <div className="bg-[#121b31] p-4 rounded-lg border border-[#334155]">
          <div className="text-sm text-gray-400 mb-2">R²</div>
          <div className="text-lg font-semibold text-white">
            {formatNumber(metrics.r_squared)}
          </div>
          <div className="text-xs text-gray-500 mt-1">
            {getRSquaredInterpretation(metrics.r_squared)}
          </div>
        </div>

        {/* Intercept */}
        <div className="bg-[#121b31] p-4 rounded-lg border border-[#334155]">
          <div className="text-sm text-gray-400 mb-2">Intercepto</div>
          <div className="text-lg font-semibold text-white">
            {formatBps(metrics.intercept)}
          </div>
          <div className="text-xs text-gray-500 mt-1">
            Valor inicial
          </div>
        </div>
      </div>

      {/* Additional Information */}
      <div className="mt-6 grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="bg-[#121b31] p-4 rounded-lg border border-[#334155]">
          <div className="text-sm text-gray-400 mb-2">Interpretación</div>
          <div className="text-white text-sm">
            {trendDirection === 'increasing' && (
              <span>El tráfico muestra una tendencia creciente con una pendiente positiva de {formatBps(metrics.slope)} por día.</span>
            )}
            {trendDirection === 'decreasing' && (
              <span>El tráfico muestra una tendencia decreciente con una pendiente negativa de {formatBps(metrics.slope)} por día.</span>
            )}
            {trendDirection === 'stable' && (
              <span>El tráfico se mantiene relativamente estable sin cambios significativos en la tendencia.</span>
            )}
          </div>
        </div>

        <div className="bg-[#121b31] p-4 rounded-lg border border-[#334155]">
          <div className="text-sm text-gray-400 mb-2">Calidad del Modelo</div>
          <div className="text-white text-sm">
            El valor R² de {formatNumber(metrics.r_squared)} indica que el modelo explica el
            {(metrics.r_squared * 100).toFixed(1)}% de la variabilidad en los datos.
            {getRSquaredInterpretation(metrics.r_squared)} para predicciones.
          </div>
        </div>
      </div>
    </div>
  )
}

export default Metrics

