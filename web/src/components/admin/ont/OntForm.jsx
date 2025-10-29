import { useState } from 'react'

export default function OntForm() {
  const [step, setStep] = useState(1) // 1: PON selection, 2: ONT selection, 3: Create
  const [formData, setFormData] = useState({
    ip: '',
    community: '',
    pon_idx: '',
    ont_idx: '',
    description: ''
  })
  const [ponPorts, setPonPorts] = useState([])
  const [ontSerials, setOntSerials] = useState([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  const [success, setSuccess] = useState(null)
  const [selectedOnt, setSelectedOnt] = useState(null)

  // Pagination states
  const [ponCurrentPage, setPonCurrentPage] = useState(1)
  const [ontCurrentPage, setOntCurrentPage] = useState(1)
  const itemsPerPage = 10

  // Pagination calculations
  const getCurrentPageItems = (items, currentPage) => {
    const startIndex = (currentPage - 1) * itemsPerPage
    const endIndex = startIndex + itemsPerPage
    return items.slice(startIndex, endIndex)
  }

  const getTotalPages = (items) => {
    return Math.ceil(items.length / itemsPerPage)
  }

  const BASE_URL = `${import.meta.env.PUBLIC_URL || ""}/api`
  const TOKEN = sessionStorage.getItem("access_token")?.replace("Bearer ", "") || ""

  const headers = {
    headers: {
      'Authorization': `Bearer ${TOKEN}`,
      'Content-Type': 'application/json'
    }
  }

  const handleChange = (e) => {
    const { name, value } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: value
    }))
  }

  // Step 1: Get PON ports
  const handleGetPonPorts = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError(null)

    try {
      const response = await fetch(`${BASE_URL}/prometheus/pons/ip/${formData.ip}`, headers)

      if (response.status === 401 || response.status === 403) {
        sessionStorage.removeItem("access_token")
        window.location.href = "/auth/login"
        return
      }

      if (!response.ok) {
        throw new Error(`Error: ${response.status}`)
      }

      const data = await response.json()
      setPonPorts(data)
      setStep(2)
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  // Step 2: Get ONT serials
  const handleGetOntSerials = async (selectedPon) => {
    setLoading(true)
    setError(null)

    try {
      const response = await fetch(`${BASE_URL}/ont/snmp/serial-despt`, {
        method: 'POST',
        body: JSON.stringify({
          ip: formData.ip,
          community: formData.community,
          pon_idx: selectedPon.idx
        }),
        ...headers
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || `Error: ${response.status}`)
      }

      const data = await response.json()
      setOntSerials(data)
      setFormData(prev => ({ ...prev, pon_idx: selectedPon.idx }))
      setStep(3)
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  // Step 3: Select ONT and prepare for creation
  const handleSelectOnt = (ont) => {
    setFormData(prev => ({ ...prev, ont_idx: ont.ont_idx }))
    setSelectedOnt(ont)
  }

  // Step 3: Create ONT
  const handleCreateOnt = async () => {
    if (!selectedOnt) {
      setError('Debe seleccionar una ONT primero')
      return
    }

    setLoading(true)
    setError(null)

    try {
      const response = await fetch(`${BASE_URL}/ont/`, {
        method: 'POST',
        body: JSON.stringify({
          ...formData,
          ont_idx: selectedOnt.ont_idx
        }),
        ...headers
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || `Error: ${response.status}`)
      }

      setSuccess('ONT creada exitosamente')
      setFormData({
        ip: '',
        community: '',
        pon_idx: '',
        ont_idx: '',
        description: ''
      })
      setPonPorts([])
      setOntSerials([])
      setSelectedOnt(null)
      setStep(1)

      // Refresh ONT list
      window.dispatchEvent(new CustomEvent('ontCreated'))
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const getStatusText = (status) => {
    switch (status) {
      case 1: return 'Activo'
      case 2: return 'Inactivo'
      default: return 'Error'
    }
  }

  const getStatusColor = (status) => {
    switch (status) {
      case 1: return 'bg-green-500'
      case 2: return 'bg-yellow-500'
      default: return 'bg-red-500'
    }
  }

  const resetForm = () => {
    setStep(1)
    setFormData({
      ip: '',
      community: '',
      pon_idx: '',
      ont_idx: '',
      description: ''
    })
    setPonPorts([])
    setOntSerials([])
    setSelectedOnt(null)
    setPonCurrentPage(1)
    setOntCurrentPage(1)
    setError(null)
    setSuccess(null)
  }

  return (
    <div>
      {error && (
        <div className="p-3 rounded bg-red-900/50 border border-red-500 text-red-200 mb-4">
          {error}
        </div>
      )}

      {success && (
        <div className="p-3 rounded bg-green-900/50 border border-green-500 text-green-200 mb-4">
          {success}
        </div>
      )}

      {/* Step 1: IP and Community */}
      {step === 1 && (
        <form onSubmit={handleGetPonPorts} className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="flex flex-col gap-1">
              <label htmlFor="ip" className="text-gray-200">
                IP del OLT *
              </label>
              <input
                id="ip"
                name="ip"
                type="text"
                value={formData.ip}
                onChange={handleChange}
                placeholder="192.168.1.1"
                required
                className="w-full border border-[hsl(217,33%,20%)] bg-[#0f1729] text-gray-200 p-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div className="flex flex-col gap-1">
              <label htmlFor="community" className="text-gray-200">
                Comunidad SNMP *
              </label>
              <input
                id="community"
                name="community"
                type="text"
                value={formData.community}
                onChange={handleChange}
                placeholder="public"
                required
                className="w-full border border-[hsl(217,33%,20%)] bg-[#0f1729] text-gray-200 p-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>

          <div className="flex justify-end">
            <button
              type="submit"
              disabled={loading}
              className="px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white rounded-lg font-medium transition-colors"
            >
              {loading ? 'Buscando...' : 'Buscar Puertos PON'}
            </button>
          </div>
        </form>
      )}

      {/* Step 2: PON Port Selection */}
      {step === 2 && (
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-semibold text-gray-200">Seleccionar Puerto PON</h3>
            <button
              onClick={resetForm}
              className="text-sm text-gray-400 hover:text-gray-200"
            >
              Cambiar IP/Comunidad
            </button>
          </div>

          {ponPorts.length === 0 ? (
            <div className="text-center py-4 text-gray-400">
              No se encontraron puertos PON
            </div>
          ) : (
            <>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-3 max-h-96 overflow-y-auto">
                {getCurrentPageItems(ponPorts, ponCurrentPage).map((port) => (
                  <div
                    key={port.idx}
                    className="p-3 rounded-lg border border-[hsl(217,33%,20%)] bg-[#0f1729] cursor-pointer hover:border-blue-500 transition-colors"
                    onClick={() => handleGetOntSerials(port)}
                  >
                    <div className="flex justify-between items-center">
                      <div>
                        <div className="font-medium text-gray-200">
                          GPON {port.shell}/{port.card}/{port.port}
                        </div>
                        <div className="text-sm text-gray-400">
                          IDX: {port.idx}
                        </div>
                      </div>
                      <div className="flex items-center gap-2">
                        <div className={`w-3 h-3 rounded-full ${getStatusColor(port.status)}`}></div>
                        <span className="text-sm text-gray-300">
                          {getStatusText(port.status)}
                        </span>
                      </div>
                    </div>
                  </div>
                ))}
              </div>

              {/* Pagination Controls */}
              {getTotalPages(ponPorts) > 1 && (
                <div className="flex justify-center items-center gap-4">
                  <button
                    onClick={() => setPonCurrentPage(prev => Math.max(1, prev - 1))}
                    disabled={ponCurrentPage === 1}
                    className="px-3 py-1 rounded border border-[hsl(217,33%,20%)] bg-[#0f1729] text-gray-300 disabled:text-gray-500 disabled:cursor-not-allowed"
                  >
                    Anterior
                  </button>
                  <span className="text-gray-300">
                    Página {ponCurrentPage} de {getTotalPages(ponPorts)}
                  </span>
                  <button
                    onClick={() => setPonCurrentPage(prev => Math.min(getTotalPages(ponPorts), prev + 1))}
                    disabled={ponCurrentPage === getTotalPages(ponPorts)}
                    className="px-3 py-1 rounded border border-[hsl(217,33%,20%)] bg-[#0f1729] text-gray-300 disabled:text-gray-500 disabled:cursor-not-allowed"
                  >
                    Siguiente
                  </button>
                </div>
              )}
            </>
          )}
        </div>
      )}

      {/* Step 3: ONT Selection */}
      {step === 3 && (
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-semibold text-gray-200">Seleccionar ONT</h3>
            <button
              onClick={() => setStep(2)}
              className="text-sm text-gray-400 hover:text-gray-200"
            >
              Volver a Puertos PON
            </button>
          </div>

          <div className="flex flex-col gap-1">
            <label htmlFor="description" className="text-gray-200">
              Descripción
            </label>
            <input
              id="description"
              name="description"
              type="text"
              value={formData.description}
              onChange={handleChange}
              placeholder="Descripción de la ONT"
              className="w-full border border-[hsl(217,33%,20%)] bg-[#0f1729] text-gray-200 p-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          {/* Selected ONT Preview */}
          {selectedOnt && (
            <div className="p-4 rounded-lg border-2 border-blue-500 bg-blue-900/20">
              <h4 className="font-semibold text-gray-200 mb-2">ONT Seleccionada:</h4>
              <div className="text-sm text-gray-300">
                <div><strong>Descripción:</strong> {selectedOnt.despt || selectedOnt.serial}</div>
                <div><strong>Serial:</strong> {selectedOnt.serial}</div>
                <div><strong>ONT ID:</strong> {selectedOnt.ont_idx}</div>
              </div>
            </div>
          )}

          {ontSerials.length === 0 ? (
            <div className="text-center py-4 text-gray-400">
              No se encontraron ONTs en este puerto PON
            </div>
          ) : (
            <>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-3 max-h-96 overflow-y-auto">
                {getCurrentPageItems(ontSerials, ontCurrentPage).map((ont) => (
                  <div
                    key={`${ont.pon_idx}-${ont.ont_idx}`}
                    className={`p-3 rounded-lg border bg-[#0f1729] cursor-pointer transition-colors ${
                      selectedOnt?.ont_idx === ont.ont_idx
                        ? 'border-blue-500 border-2'
                        : 'border-[hsl(217,33%,20%)] hover:border-blue-500'
                    }`}
                    onClick={() => handleSelectOnt(ont)}
                  >
                    <div className="flex flex-col gap-1">
                      <div className="font-medium text-gray-200 truncate">
                        {ont.despt || ont.serial}
                      </div>
                      <div className="text-sm text-gray-400 truncate">
                        Serial: {ont.serial}
                      </div>
                      <div className="text-sm text-gray-300">
                        ONT ID: {ont.ont_idx}
                      </div>
                    </div>
                  </div>
                ))}
              </div>

              {/* Pagination Controls */}
              {getTotalPages(ontSerials) > 1 && (
                <div className="flex justify-center items-center gap-4">
                  <button
                    onClick={() => setOntCurrentPage(prev => Math.max(1, prev - 1))}
                    disabled={ontCurrentPage === 1}
                    className="px-3 py-1 rounded border border-[hsl(217,33%,20%)] bg-[#0f1729] text-gray-300 disabled:text-gray-500 disabled:cursor-not-allowed"
                  >
                    Anterior
                  </button>
                  <span className="text-gray-300">
                    Página {ontCurrentPage} de {getTotalPages(ontSerials)}
                  </span>
                  <button
                    onClick={() => setOntCurrentPage(prev => Math.min(getTotalPages(ontSerials), prev + 1))}
                    disabled={ontCurrentPage === getTotalPages(ontSerials)}
                    className="px-3 py-1 rounded border border-[hsl(217,33%,20%)] bg-[#0f1729] text-gray-300 disabled:text-gray-500 disabled:cursor-not-allowed"
                  >
                    Siguiente
                  </button>
                </div>
              )}
            </>
          )}

          {/* Confirm Button */}
          {selectedOnt && (
            <div className="flex justify-end">
              <button
                onClick={handleCreateOnt}
                disabled={loading}
                className="px-4 py-2 bg-green-600 hover:bg-green-700 disabled:bg-green-400 text-white rounded-lg font-medium transition-colors"
              >
                {loading ? 'Creando...' : 'Confirmar y Crear ONT'}
              </button>
            </div>
          )}
        </div>
      )}
    </div>
  )
}

