import { useState, useEffect } from 'react'
import { useDebounce } from '../../hooks/useDebounce'
import dayjs from 'dayjs'

const BASE_URL_TRAFFIC = `${import.meta.env.PUBLIC_URL || ""}/api/traffic`
const TOKEN = sessionStorage.getItem("access_token")?.replace("Bearer ", "") || ""

const Form = () => {
  const [selectedRegion, setSelectedRegion] = useState('')
  const [selectedState, setSelectedState] = useState('')
  const [selectedOlt, setSelectedOlt] = useState('')
  const [futureDays, setFutureDays] = useState(7)
  const [confidence, setConfidence] = useState(0.95)
  const [dateRange, setDateRange] = useState({
    initDate: dayjs().subtract(1, 'year').format('YYYY-MM-DDTHH:mm:ss-04:00'),
    finalDate: dayjs().format('YYYY-MM-DDTHH:mm:ss-04:00')
  })
  const [olts, setOlts] = useState([])
  const [regions, setRegions] = useState([])
  const [states, setStates] = useState([])
  const [infoAllOlt, setInfoAllOlt] = useState([])

  const debouncedFutureDays = useDebounce(futureDays, 300)
  const debouncedConfidence = useDebounce(confidence, 300)
  const debouncedDateRange = useDebounce(dateRange, 300)

  // Fetch regions and states when component mounts
  useEffect(() => {
    fetchRegionsAndStates()
  }, [])

  // Fetch OLTS when state is selected
  useEffect(() => {
    if (selectedState) {
      fetchOlts(selectedState)
    } else {
      setOlts([])
      setSelectedOlt('')
    }
  }, [selectedState])

  const fetchOlts = async (state) => {
    try {
      const response = await fetch(`${BASE_URL_TRAFFIC}/hierarchy?initDate=${dateRange.initDate}&finalDate=${dateRange.finalDate}`, {
        headers: { Authorization: `Bearer ${TOKEN}` }
      })

      if (response.status === 401 || response.status === 403) {
        sessionStorage.removeItem("access_token")
        window.location.href = "/"
        return
      }

      if (response.ok) {
        const data = await response.json()
        // Get OLTS for the selected state from the hierarchy
        const oltList = (data.olts[state] || [])
          .map(olt => ({
            ip: olt.ip,
            sysName: olt.sys_name,
            displayText: `${olt.sys_name} (${olt.ip})`
          }))
          .sort((a, b) => a.sysName.localeCompare(b.sysName))

        setOlts(oltList)
      }
    } catch (error) {
      console.error('Error fetching OLTS:', error)
    }
  }

  const fetchRegionsAndStates = async () => {
    try {
      const response = await fetch(`${BASE_URL_TRAFFIC}/hierarchy?initDate=${dateRange.initDate}&finalDate=${dateRange.finalDate}`, {
        headers: { Authorization: `Bearer ${TOKEN}` }
      })

      if (response.status === 401 || response.status === 403) {
        sessionStorage.removeItem("access_token")
        window.location.href = "/"
        return
      }

      if (response.ok) {
        const data = await response.json()
        // Get regions from the hierarchy
        const uniqueRegions = data.regions
          .filter(region => region)
          .sort()
          .map(region => ({ value: region, label: region }))

        // Get all states from the hierarchy
        const allStates = []
        Object.values(data.states).forEach(stateList => {
          allStates.push(...stateList)
        })
        const uniqueStates = [...new Set(allStates)]
          .filter(state => state)
          .sort()
          .map(state => ({ value: state, label: state }))

        setRegions(uniqueRegions)
        setStates(uniqueStates)
        setInfoAllOlt(data)
      }
    } catch (error) {
      console.error('Error fetching regions and states:', error)
    }
  }

  const handleSubmit = (e) => {
    e.preventDefault()

    // Determine selected level based on form selections
    let selectedLevel = 'national'
    if (selectedOlt) {
      selectedLevel = 'olt'
    } else if (selectedState) {
      selectedLevel = 'state'
    } else if (selectedRegion) {
      selectedLevel = 'regional'
    }

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
      {/* Región */}
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
          {regions.map(region => (
            <option key={region.value} value={region.value}>{region.label}</option>
          ))}
        </select>
      </div>

      {/* Estado: se habilita solo si región está seleccionada */}
      {selectedRegion && (
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
            {states
              .filter(state => {
                // Filter states by selected region using the hierarchy data
                const statesInRegion = infoAllOlt?.states[selectedRegion] || []
                return statesInRegion.includes(state.value)
              })
              .map(state => (
                <option key={state.value} value={state.value}>{state.label}</option>
              ))}
          </select>
        </div>
      )}

      {/* OLT: se habilita solo si estado está seleccionado */}
      {selectedState && (
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            OLT
          </label>
          <select
            value={selectedOlt}
            onChange={(e) => setSelectedOlt(e.target.value)}
            className="w-full px-3 py-2 bg-[#1e293b] border border-[#334155] rounded-md text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="">Seleccionar OLT</option>
            {olts.map(olt => (
              <option key={olt.ip} value={olt.ip}>{olt.displayText}</option>
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
                initDate: dayjs(e.target.value).format('YYYY-MM-DDTHH:mm:ss-04:00')
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
                finalDate: dayjs(e.target.value).format('YYYY-MM-DDTHH:mm:ss-04:00')
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

