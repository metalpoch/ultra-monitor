import { useState, useEffect } from 'react'
import useFetch from '../../../hooks/useFetch'

export default function OntList() {
  const [onts, setOnts] = useState([])
  const [selectedOnt, setSelectedOnt] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [refreshTrigger, setRefreshTrigger] = useState(0)

  const BASE_URL = `${import.meta.env.PUBLIC_URL || ""}/api/ont`
  const TOKEN = sessionStorage.getItem("access_token")?.replace("Bearer ", "") || ""

  const headers = {
    headers: {
      Authorization: `Bearer ${TOKEN}`
    }
  }

  // Fetch all ONTs
  const fetchOnts = async () => {
    try {
      setLoading(true)
      const response = await fetch(`${BASE_URL}/`, headers)

      if (response.status === 401 || response.status === 403) {
        sessionStorage.removeItem("access_token")
        window.location.href = "/auth/login"
        return
      }

      if (!response.ok) {
        throw new Error(`Error: ${response.status}`)
      }

      const data = await response.json()
      setOnts(data)
      setError(null)
    } catch (err) {
      setError(err.message)
      console.error('Error fetching ONTs:', err)
    } finally {
      setLoading(false)
    }
  }

  // Listen for ONT creation events
  useEffect(() => {
    const handleOntCreated = () => {
      console.log('ONT created, refreshing list...')
      fetchOnts()
    }

    window.addEventListener('ontCreated', handleOntCreated)
    return () => {
      window.removeEventListener('ontCreated', handleOntCreated)
    }
  }, [])

  useEffect(() => {
    fetchOnts()
  }, [refreshTrigger])

  // Handle ONT selection
  const handleSelectOnt = (ont) => {
    setSelectedOnt(ont)
    // Emit event for traffic chart
    window.dispatchEvent(new CustomEvent('ontSelected', { detail: ont }))
  }

  // Handle enable/disable ONT
  const handleToggleOnt = async (ontId, enable) => {
    try {
      const endpoint = enable ? 'enable' : 'disable'
      const response = await fetch(`${BASE_URL}/${ontId}/${endpoint}`, {
        method: 'PATCH',
        ...headers
      })

      if (!response.ok) {
        throw new Error(`Error: ${response.status}`)
      }

      // Refresh ONT list
      fetchOnts()
    } catch (err) {
      console.error('Error toggling ONT:', err)
      alert('Error al cambiar el estado de la ONT')
    }
  }

  // Handle delete ONT
  const handleDeleteOnt = async (ontId) => {
    if (!confirm('¿Está seguro de que desea eliminar esta ONT?')) {
      return
    }

    try {
      const response = await fetch(`${BASE_URL}/${ontId}`, {
        method: 'DELETE',
        ...headers
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || `Error: ${response.status}`)
      }

      // Refresh ONT list and clear selection if deleted
      fetchOnts()
      if (selectedOnt && selectedOnt.id === ontId) {
        setSelectedOnt(null)
        window.dispatchEvent(new CustomEvent('ontSelected', { detail: null }))
      }
    } catch (err) {
      console.error('Error deleting ONT:', err)
      if (err.message.includes('enabled') || err.message.includes('habilitada')) {
        alert('Error: La ONT debe estar deshabilitada antes de poder eliminarla')
      } else {
        alert(`Error al eliminar la ONT: ${err.message}`)
      }
    }
  }

  const getStatusColor = (enabled, status) => {
    if (!enabled) return 'bg-gray-500'
    return status ? 'bg-green-500' : 'bg-red-500'
  }

  const getStatusText = (enabled, status) => {
    if (!enabled) return 'Deshabilitada'
    return status ? 'Activa' : 'Inactiva'
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center py-8">
        <div className="text-gray-400">Cargando ONTs...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex items-center justify-center py-8">
        <div className="text-red-400">Error al cargar ONTs: {error}</div>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {onts.length === 0 ? (
        <div className="text-center py-8 text-gray-400">
          No hay ONTs registradas
        </div>
      ) : (
        onts.map((ont) => (
          <div
            key={ont.id}
            className={`ont-card p-4 rounded-lg border-2 border-[hsl(217,33%,20%)] bg-[#0f1729] cursor-pointer transition-all duration-200 ${
              selectedOnt?.id === ont.id ? 'selected border-blue-500' : ''
            }`}
            onClick={() => handleSelectOnt(ont)}
          >
            <div className="flex justify-between items-start mb-3">
              <div>
                <h3 className="text-lg font-semibold text-gray-200">
                  {ont.serial} - {ont.despt}
                </h3>
                {ont.description && (
                  <p className="text-sm text-blue-300 font-medium mt-1 bg-blue-900/30 px-2 py-1 rounded">
                    📝 {ont.description}
                  </p>
                )}
              </div>
              <div className="flex items-center gap-2">
                <div className={`w-3 h-3 rounded-full ${getStatusColor(ont.enabled, ont.status)}`}></div>
                <span className="text-sm text-gray-300">
                  {getStatusText(ont.enabled, ont.status)}
                </span>
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4 text-sm text-gray-300">
              <div>
                <span className="text-gray-400">IP:</span> {ont.ip}
              </div>
              <div>
                <span className="text-gray-400">ONT ID:</span> {ont.ont_idx}
              </div>
              <div>
                <span className="text-gray-400">Perfil:</span> {ont.line_prof}
              </div>
              <div>
                <span className="text-gray-400">Distancia:</span> {ont.olt_distance}m
              </div>
              <div>
                <span className="text-gray-400">Estado:</span>
                <span className={ont.status ? 'text-green-400' : 'text-red-400'}>
                  {ont.status ? 'Activo' : 'Inactivo'}
                </span>
              </div>
              <div>
                <span className="text-gray-400">Última verificación:</span>
                {new Date(ont.last_check).toLocaleString('es-ES')}
              </div>
            </div>

            <div className="flex justify-end gap-2 mt-4">
              <button
                onClick={(e) => {
                  e.stopPropagation()
                  handleToggleOnt(ont.id, !ont.enabled)
                }}
                className={`px-3 py-1 rounded text-sm font-medium ${
                  ont.enabled
                    ? 'bg-yellow-600 hover:bg-yellow-700 text-white'
                    : 'bg-green-600 hover:bg-green-700 text-white'
                }`}
              >
                {ont.enabled ? 'Deshabilitar' : 'Habilitar'}
              </button>

              <button
                onClick={(e) => {
                  e.stopPropagation()
                  handleDeleteOnt(ont.id)
                }}
                className="px-3 py-1 rounded text-sm font-medium bg-red-600 hover:bg-red-700 text-white"
              >
                Eliminar
              </button>
            </div>
          </div>
        ))
      )}
    </div>
  )
}

