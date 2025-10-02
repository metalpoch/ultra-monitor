import { useState, useEffect } from 'react'
import { useDebounce } from '../../hooks/useDebounce'
import { REGIONS, STATES_BY_REGION } from '../../constants/regions'
import dayjs from 'dayjs'

const BASE_URL_TRAFFIC = `${import.meta.env.PUBLIC_URL || ""}/api/traffic`
const TOKEN = sessionStorage.getItem("access_token")?.replace("Bearer ", "") || ""

const Form = () => {
  const [selectedLevel, setSelectedLevel] = useState('national')
  const [selectedRegion, setSelectedRegion] = useState('')
  const [selectedState, setSelectedState] = useState('')
  const [selectedOlt, setSelectedOlt] = useState('')
  const [futureDays, setFutureDays] = useState(7)
  const [confidence, setConfidence] = useState(0.95)
  const [dateRange, setDateRange] = useState({
    initDate: dayjs().subtract(30, 'day').toISOString(),
    finalDate: dayjs().toISOString()
  })
  const [olts, setOlts] = useState([])

  const debouncedFutureDays = useDebounce(futureDays, 300)
  const debouncedConfidence = useDebounce(confidence, 300)
  const debouncedDateRange = useDebounce(dateRange, 300)

  // Fetch OLTS when state is selected
  useEffect(() => {
    if (selectedState) {
      fetchOlts(selectedState)
      setSelectedOlt('')
      setOlts([])
    } else {
      setOlts([])
      setSelectedOlt('')
    }
  }, [selectedState])

  const fetchOlts = async (state) => {
    try {
      const response = await fetch(`${BASE_URL_TRAFFIC}/sysname/${encodeURIComponent(state)}`, {
        headers: { Authorization: `Bearer ${TOKEN}` }
      })

      if (response.status === 401 || response.status === 403) {
        sessionStorage.removeItem("access_token")
        window.location.href = "/"
        return
      }

      if (response.ok) {
        const data = await response.json()
        const oltList = Object.keys(data).sort()
        setOlts(oltList)
      }
    } catch (error) {
      console.error('Error fetching OLTS:', error)
    }
  }

  const handleSubmit = (e) => {
    e.preventDefault()

    const trendData = {
      futureDays: debouncedFutureDays,
      confidence: debouncedConfidence,
      initDate: debouncedDateRange.initDate,
      finalDate: debouncedDateRange.finalDate,
      selectedLevel,
      selectedRegion,
      selectedState,
      selectedOlt
    }

    console.log('Form submitted with data:', trendData) // Debug log

    // Dispatch custom event with trend data
    const event = new CustomEvent('trendFormSubmit', { detail: trendData })
    window.dispatchEvent(event)
  }

  // Auto-submit when component mounts to show default national trend
  useEffect(() => {
    const handleSubmit = () => {
      const trendData = {
        futureDays: debouncedFutureDays,
        confidence: debouncedConfidence,
        initDate: debouncedDateRange.initDate,
        finalDate: debouncedDateRange.finalDate,
        selectedLevel: 'national',
        selectedRegion: '',
        selectedState: '',
        selectedOlt: ''
      }

      const event = new CustomEvent('trendFormSubmit', { detail: trendData })
      window.dispatchEvent(event)
    }

    // Auto-submit after a short delay to ensure form is ready
    const timer = setTimeout(handleSubmit, 100)
    return () => clearTimeout(timer)
  }, [])

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4">
      <div>
        <label className="block text-sm font-medium text-gray-300 mb-2">
          Nivel de Visualización
        </label>
        <select
          value={selectedLevel}
          onChange={(e) => {
            setSelectedLevel(e.target.value)
            setSelectedRegion('')
            setSelectedState('')
            setSelectedOlt('')
            setOlts([])
          }}
          className="w-full px-3 py-2 bg-[#1e293b] border border-[#334155] rounded-md text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="national">Nacional</option>
          <option value="regional">Regional</option>
          <option value="state">Estatal</option>
          <option value="olt">OLT</option>
        </select>
      </div>

      {selectedLevel === 'regional' && (
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Región
          </label>
          <select
            value={selectedRegion}
            onChange={(e) => {
              setSelectedRegion(e.target.value)
              setSelectedState('')
              setSelectedOlt('')
              setOlts([])
            }}
            className="w-full px-3 py-2 bg-[#1e293b] border border-[#334155] rounded-md text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="">Seleccionar Región</option>
            {REGIONS.map(region => (
              <option key={region.value} value={region.value}>{region.label}</option>
            ))}
          </select>
        </div>
      )}

      {(selectedLevel === 'state' || (selectedLevel === 'regional' && selectedRegion)) && (
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Estado
          </label>
          <select
            value={selectedState}
            onChange={(e) => {
              setSelectedState(e.target.value)
              setSelectedOlt('')
              setOlts([])
            }}
            className="w-full px-3 py-2 bg-[#1e293b] border border-[#334155] rounded-md text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="">Seleccionar Estado</option>
            {selectedLevel === 'regional' && selectedRegion ? (
              STATES_BY_REGION[selectedRegion]?.map(state => (
                <option key={state} value={state}>{state}</option>
              ))
            ) : (
              Object.values(STATES_BY_REGION).flat().sort().map(state => (
                <option key={state} value={state}>{state}</option>
              ))
            )}
          </select>
        </div>
      )}

      {(selectedLevel === 'olt' || (selectedLevel === 'state' && selectedState)) && (
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            OLT
          </label>
          <select
            value={selectedOlt}
            onChange={(e) => setSelectedOlt(e.target.value)}
            disabled={!selectedState}
            className="w-full px-3 py-2 bg-[#1e293b] border border-[#334155] rounded-md text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <option value="">Seleccionar OLT</option>
            {olts.map(olt => (
              <option key={olt} value={olt}>{olt}</option>
            ))}
          </select>
        </div>
      )}

      <div>
        <label className="block text-sm font-medium text-gray-300 mb-2">
          Rango de Fechas
        </label>
        <div className="flex flex-col gap-2">
          <div>
            <label className="block text-xs text-gray-400 mb-1">Fecha Inicial</label>
            <input
              type="date"
              value={dayjs(dateRange.initDate).format('YYYY-MM-DD')}
              onChange={(e) => setDateRange(prev => ({
                ...prev,
                initDate: dayjs(e.target.value).toISOString()
              }))}
              className="w-full px-3 py-2 bg-[#1e293b] border border-[#334155] rounded-md text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label className="block text-xs text-gray-400 mb-1">Fecha Final</label>
            <input
              type="date"
              value={dayjs(dateRange.finalDate).format('YYYY-MM-DD')}
              onChange={(e) => setDateRange(prev => ({
                ...prev,
                finalDate: dayjs(e.target.value).toISOString()
              }))}
              className="w-full px-3 py-2 bg-[#1e293b] border border-[#334155] rounded-md text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-300 mb-2">
          Días de Predicción
        </label>
        <input
          type="number"
          min="1"
          max="30"
          value={futureDays}
          onChange={(e) => setFutureDays(parseInt(e.target.value))}
          className="w-full px-3 py-2 bg-[#1e293b] border border-[#334155] rounded-md text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-300 mb-2">
          Nivel de Confianza
        </label>
        <select
          value={confidence}
          onChange={(e) => setConfidence(parseFloat(e.target.value))}
          className="w-full px-3 py-2 bg-[#1e293b] border border-[#334155] rounded-md text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value={0.90}>90%</option>
          <option value={0.95}>95%</option>
          <option value={0.99}>99%</option>
        </select>
      </div>

      <button
        type="submit"
        className="w-full px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-md transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-[#121b31]"
      >
        Generar Tendencia
      </button>
    </form>
  )
}

export default Form

