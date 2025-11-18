import { useState, useEffect } from "react";
import useFetch from "../../hooks/useFetch";
import DatalistField from "../ui/DatalistField";

const BASE_URL = `${import.meta.env.PUBLIC_URL || ""}/api`;
const TOKEN = sessionStorage.getItem("access_token")?.replace("Bearer ", "") || "";

export default function InterfacesManagement() {
  // Fetch data from both endpoints
  const {
    data: bandwidthData,
    loading: bandwidthLoading,
    error: bandwidthError,
    refetch: refetchBandwidth
  } = useFetch(`${BASE_URL}/interfaces-bandwidth/`, {
    headers: { Authorization: `Bearer ${TOKEN}` }
  });

  const {
    data: oltData,
    loading: oltLoading,
    error: oltError,
    refetch: refetchOlt
  } = useFetch(`${BASE_URL}/interfaces-olt/`, {
    headers: { Authorization: `Bearer ${TOKEN}` }
  });

  const [olts, setOlts] = useState([]);
  const [interfaces, setInterfaces] = useState([]);
  const [selectedOlt, setSelectedOlt] = useState(null);
  const [suggestedVerbose, setSuggestedVerbose] = useState("");
  const [customVerbose, setCustomVerbose] = useState("");
  const [updateLoading, setUpdateLoading] = useState(false);
  const [updateMessage, setUpdateMessage] = useState({ type: "", text: "" });
  const [availableVerboseNames, setAvailableVerboseNames] = useState([]);
  const [validVerboseNames, setValidVerboseNames] = useState([]);

  // Filter states
  const [filters, setFilters] = useState({
    region: "",
    state: "",
    sysname: "",
    ip: "",
    olt_verbose: ""
  });

  // Sorting state
  const [sortConfig, setSortConfig] = useState({
    key: null,
    direction: 'asc' // 'asc' or 'desc'
  });

  // Handle sorting
  const handleSort = (key) => {
    let direction = 'asc';
    if (sortConfig.key === key && sortConfig.direction === 'asc') {
      direction = 'desc';
    }
    setSortConfig({ key, direction });
  };

  // Filter OLTs based on filter criteria
  const filteredOlts = olts.filter(olt => {
    return (
      (filters.region === "" || olt.region?.toLowerCase().includes(filters.region.toLowerCase())) &&
      (filters.state === "" || olt.state?.toLowerCase().includes(filters.state.toLowerCase())) &&
      (filters.sysname === "" || olt.sysname?.toLowerCase().includes(filters.sysname.toLowerCase())) &&
      (filters.ip === "" || olt.ip?.includes(filters.ip)) &&
      (filters.olt_verbose === "" ||
        (olt.olt_verbose?.toLowerCase().includes(filters.olt_verbose.toLowerCase()) ||
         (filters.olt_verbose.toLowerCase() === "sin asignar" && !olt.olt_verbose)))
    );
  });

  // Sort filtered OLTs
  const sortedOlts = [...filteredOlts].sort((a, b) => {
    if (!sortConfig.key) return 0;

    const aValue = a[sortConfig.key] || "";
    const bValue = b[sortConfig.key] || "";

    if (aValue < bValue) {
      return sortConfig.direction === 'asc' ? -1 : 1;
    }
    if (aValue > bValue) {
      return sortConfig.direction === 'asc' ? 1 : -1;
    }
    return 0;
  });

  // Handle filter changes
  const handleFilterChange = (field, value) => {
    setFilters(prev => ({
      ...prev,
      [field]: value
    }));
  };

  // Clear all filters
  const clearFilters = () => {
    setFilters({
      region: "",
      state: "",
      sysname: "",
      ip: "",
      olt_verbose: ""
    });
  };

  // Process data when fetched
  useEffect(() => {
    if (oltData) {
      setOlts(oltData);
    }
  }, [oltData]);

  useEffect(() => {
    if (bandwidthData) {
      setInterfaces(bandwidthData);
      // Extract unique olt_verbose values with interface information for datalist
      const verboseNameMap = new Map();

      bandwidthData.forEach(intf => {
        if (intf.olt_verbose) {
          if (!verboseNameMap.has(intf.olt_verbose)) {
            verboseNameMap.set(intf.olt_verbose, []);
          }
          verboseNameMap.get(intf.olt_verbose).push(intf.interface);
        }
      });

      // Create datalist options with interface information
      const datalistOptions = Array.from(verboseNameMap.entries()).map(([verboseName, interfaces]) => {
        const interfaceCount = interfaces.length;
        const sampleInterface = interfaces[0];
        const label = interfaceCount > 1
          ? `${verboseName} (${interfaceCount} interfaces, ej: ${sampleInterface})`
          : `${verboseName} (${sampleInterface})`;

        return {
          value: verboseName,
          label: label
        };
      }).sort((a, b) => a.value.localeCompare(b.value));

      setAvailableVerboseNames(datalistOptions);
      // Extract valid values for validation
      setValidVerboseNames(Array.from(verboseNameMap.keys()));
    }
  }, [bandwidthData]);

  // Handle OLT selection
  const handleOltSelect = (olt) => {
    setSelectedOlt(olt);
    setCustomVerbose(olt.olt_verbose || "");

    // Find matching interface and suggest olt_verbose
    const matchingInterface = interfaces.find(intf =>
      intf.interface.includes(olt.sysname) ||
      intf.interface.includes(olt.ip.replace(/\./g, "-"))
    );

    if (matchingInterface) {
      setSuggestedVerbose(matchingInterface.olt_verbose);
      if (!olt.olt_verbose) {
        setCustomVerbose(matchingInterface.olt_verbose);
      }
    } else {
      setSuggestedVerbose("");
    }
  };

  // Handle verbose name update
  const handleUpdateVerbose = async () => {
    if (!selectedOlt || !customVerbose.trim()) {
      setUpdateMessage({ type: "error", text: "Selecciona un OLT y proporciona un nombre descriptivo" });
      return;
    }

    // Validate that the olt_verbose exists in the datalist
    if (!validVerboseNames.includes(customVerbose.trim())) {
      setUpdateMessage({
        type: "error",
        text: "El nombre descriptivo debe coincidir con uno de los valores disponibles en la lista"
      });
      return;
    }

    setUpdateLoading(true);
    setUpdateMessage({ type: "", text: "" });

    try {
      const response = await fetch(`${BASE_URL}/interfaces-olt/`, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${TOKEN}`
        },
        body: JSON.stringify({
          olt_ip: selectedOlt.ip,
          olt_verbose: customVerbose.trim()
        })
      });

      if (response.ok) {
        setUpdateMessage({ type: "success", text: "OLT actualizado correctamente" });
        // Update local state
        setOlts(prev => prev.map(olt =>
          olt.ip === selectedOlt.ip
            ? { ...olt, olt_verbose: customVerbose.trim() }
            : olt
        ));
        setSelectedOlt(prev => ({ ...prev, olt_verbose: customVerbose.trim() }));

        // Clear message after 3 seconds
        setTimeout(() => setUpdateMessage({ type: "", text: "" }), 3000);
      } else {
        const errorData = await response.json();
        setUpdateMessage({
          type: "error",
          text: errorData.error || "Error al actualizar el OLT"
        });
      }
    } catch (err) {
      setUpdateMessage({ type: "error", text: "Error de conexión" });
    } finally {
      setUpdateLoading(false);
    }
  };

  // Handle authentication errors
  useEffect(() => {
    if ((oltData?.status === 401 || oltData?.status === 403) ||
        (bandwidthData?.status === 401 || bandwidthData?.status === 403)) {
      sessionStorage.removeItem("access_token");
      window.location.href = "/";
    }
  }, [oltData, bandwidthData]);

  const isLoading = bandwidthLoading || oltLoading;
  const hasError = bandwidthError || oltError;

  if (isLoading) {
    return (
      <div className="flex justify-center items-center py-8">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (hasError) {
    return (
      <div className="bg-red-900/20 border border-red-500 text-red-200 p-4 rounded-lg">
        <p>Error al cargar los datos. Por favor, intenta nuevamente.</p>
        <button
          onClick={() => { refetchBandwidth(); refetchOlt(); }}
          className="mt-2 px-4 py-2 bg-red-600 hover:bg-red-700 rounded text-white"
        >
          Reintentar
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* OLTs List - Takes 2/3 of the space */}
        <div className="lg:col-span-2 bg-[#0f172a] rounded-lg border border-gray-700 overflow-hidden">
          <div className="p-4 border-b border-gray-700">
            <div className="flex justify-between items-center">
              <div>
                <h2 className="text-lg font-semibold text-gray-200">Lista de OLTs</h2>
                <p className="text-sm text-gray-400 mt-1">
                  Selecciona un OLT para asignarle un nombre descriptivo
                </p>
              </div>
              {/* Filter Status and Clear Button */}
              <div className="flex items-center gap-3">
                {Object.values(filters).some(filter => filter !== "") && (
                  <div className="text-sm text-blue-300">
                    Filtros activos ({Object.values(filters).filter(f => f !== "").length})
                  </div>
                )}
                <button
                  onClick={clearFilters}
                  disabled={Object.values(filters).every(filter => filter === "")}
                  className="px-3 py-1 text-sm bg-gray-700 hover:bg-gray-600 disabled:bg-gray-800 disabled:text-gray-500 disabled:cursor-not-allowed rounded text-gray-200 transition-colors"
                >
                  Limpiar Filtros
                </button>
              </div>
            </div>
          </div>
          <div className="max-h-96 overflow-y-auto">
            <table className="w-full">
              <thead className="sticky top-0 bg-[#121b31] text-left">
                <tr>
                  <th
                    className="p-3 text-sm font-medium text-gray-300 cursor-pointer hover:bg-gray-700/50 transition-colors"
                    onClick={() => handleSort("region")}
                  >
                    <div className="flex items-center gap-1">
                      Región
                      {sortConfig.key === "region" && (
                        sortConfig.direction === 'asc' ? '↑' : '↓'
                      )}
                    </div>
                  </th>
                  <th
                    className="p-3 text-sm font-medium text-gray-300 cursor-pointer hover:bg-gray-700/50 transition-colors"
                    onClick={() => handleSort("state")}
                  >
                    <div className="flex items-center gap-1">
                      Estado
                      {sortConfig.key === "state" && (
                        sortConfig.direction === 'asc' ? '↑' : '↓'
                      )}
                    </div>
                  </th>
                  <th
                    className="p-3 text-sm font-medium text-gray-300 cursor-pointer hover:bg-gray-700/50 transition-colors"
                    onClick={() => handleSort("sysname")}
                  >
                    <div className="flex items-center gap-1">
                      Sysname
                      {sortConfig.key === "sysname" && (
                        sortConfig.direction === 'asc' ? '↑' : '↓'
                      )}
                    </div>
                  </th>
                  <th
                    className="p-3 text-sm font-medium text-gray-300 cursor-pointer hover:bg-gray-700/50 transition-colors"
                    onClick={() => handleSort("ip")}
                  >
                    <div className="flex items-center gap-1">
                      IP
                      {sortConfig.key === "ip" && (
                        sortConfig.direction === 'asc' ? '↑' : '↓'
                      )}
                    </div>
                  </th>
                  <th
                    className="p-3 text-sm font-medium text-gray-300 cursor-pointer hover:bg-gray-700/50 transition-colors"
                    onClick={() => handleSort("olt_verbose")}
                  >
                    <div className="flex items-center gap-1">
                      OLT Verbose
                      {sortConfig.key === "olt_verbose" && (
                        sortConfig.direction === 'asc' ? '↑' : '↓'
                      )}
                    </div>
                  </th>
                </tr>
                {/* Filter Row */}
                <tr className="bg-[#0f172a]">
                  <th className="p-2">
                    <input
                      type="text"
                      placeholder="Filtrar región..."
                      value={filters.region}
                      onChange={(e) => handleFilterChange("region", e.target.value)}
                      className="w-full px-2 py-1 text-sm bg-gray-800 border border-gray-600 rounded text-gray-200 placeholder-gray-500 focus:outline-none focus:border-blue-500"
                    />
                  </th>
                  <th className="p-2">
                    <input
                      type="text"
                      placeholder="Filtrar estado..."
                      value={filters.state}
                      onChange={(e) => handleFilterChange("state", e.target.value)}
                      className="w-full px-2 py-1 text-sm bg-gray-800 border border-gray-600 rounded text-gray-200 placeholder-gray-500 focus:outline-none focus:border-blue-500"
                    />
                  </th>
                  <th className="p-2">
                    <input
                      type="text"
                      placeholder="Filtrar sysname..."
                      value={filters.sysname}
                      onChange={(e) => handleFilterChange("sysname", e.target.value)}
                      className="w-full px-2 py-1 text-sm bg-gray-800 border border-gray-600 rounded text-gray-200 placeholder-gray-500 focus:outline-none focus:border-blue-500"
                    />
                  </th>
                  <th className="p-2">
                    <input
                      type="text"
                      placeholder="Filtrar IP..."
                      value={filters.ip}
                      onChange={(e) => handleFilterChange("ip", e.target.value)}
                      className="w-full px-2 py-1 text-sm bg-gray-800 border border-gray-600 rounded text-gray-200 placeholder-gray-500 focus:outline-none focus:border-blue-500"
                    />
                  </th>
                  <th className="p-2">
                    <input
                      type="text"
                      placeholder="Filtrar OLT verbose..."
                      value={filters.olt_verbose}
                      onChange={(e) => handleFilterChange("olt_verbose", e.target.value)}
                      className="w-full px-2 py-1 text-sm bg-gray-800 border border-gray-600 rounded text-gray-200 placeholder-gray-500 focus:outline-none focus:border-blue-500"
                    />
                  </th>
                </tr>
              </thead>
              <tbody>
                {sortedOlts.map((olt, index) => (
                  <tr
                    key={`${olt.ip}-${index}`}
                    className={`cursor-pointer hover:bg-gray-800/50 transition-colors ${
                      selectedOlt?.ip === olt.ip ? 'bg-blue-900/20 border-l-4 border-blue-500' : ''
                    } ${!olt.olt_verbose ? 'text-yellow-300' : 'text-gray-300'}`}
                    onClick={() => handleOltSelect(olt)}
                  >
                    <td className="p-3 text-sm">{olt.region}</td>
                    <td className="p-3 text-sm">{olt.state}</td>
                    <td className="p-3 text-sm font-mono">{olt.sysname}</td>
                    <td className="p-3 text-sm font-mono">{olt.ip}</td>
                    <td className="p-3 text-sm font-mono">
                      {olt.olt_verbose || (
                        <span className="text-yellow-400 italic">Sin asignar</span>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Update Form - Takes 1/3 of the space */}
        <div className="lg:col-span-1">
          {selectedOlt ? (
            <div className="bg-[#0f172a] rounded-lg border border-gray-700 p-6 sticky top-6">
              <h3 className="text-lg font-semibold text-gray-200 mb-4">
                Actualizar OLT
              </h3>

              {/* Selected OLT Info */}
              <div className="mb-4">
                <div className="bg-gray-800/50 rounded p-3 space-y-2 text-sm">
                  <div className="flex justify-between">
                    <span className="text-gray-400">Sysname:</span>
                    <span className="text-gray-200 font-mono">{selectedOlt.sysname}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-400">IP:</span>
                    <span className="text-gray-200 font-mono">{selectedOlt.ip}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-400">OLT Verbose:</span>
                    <span className="text-gray-200 font-mono">
                      {selectedOlt.olt_verbose || (
                        <span className="text-yellow-400 italic">Sin asignar</span>
                      )}
                    </span>
                  </div>
                </div>
              </div>

              {/* Suggested Interface */}
              {suggestedVerbose && (
                <div className="bg-blue-900/20 border border-blue-500 rounded p-3 mb-4">
                  <p className="text-sm text-blue-300">
                    <strong>Sugerencia encontrada:</strong> {suggestedVerbose}
                  </p>
                  <p className="text-xs text-blue-400 mt-1">
                    Basado en las interfaces de ancho de banda
                  </p>
                </div>
              )}

              {/* Update Form */}
              <div className="space-y-4">
                <DatalistField
                  id="oltVerbose"
                  label="Nombre Descriptivo (olt_verbose)"
                  options={availableVerboseNames}
                  value={customVerbose}
                  onChange={(e) => setCustomVerbose(e.target.value)}
                  placeholder="Ej: OLT-HW-SANTA-TERESA-DEL-TUY"
                />
                <p className="text-xs text-gray-400 mt-1">
                  Asigna un nombre descriptivo para identificar este OLT. Puedes seleccionar de la lista o escribir uno nuevo.
                </p>

                {/* Action Buttons */}
                <div className="flex gap-3">
                  <button
                    onClick={handleUpdateVerbose}
                    disabled={updateLoading || !customVerbose.trim()}
                    className="px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 disabled:cursor-not-allowed rounded text-white font-medium transition-colors"
                  >
                    {updateLoading ? (
                      <span className="flex items-center gap-2">
                        <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                        Actualizando...
                      </span>
                    ) : (
                      "Actualizar OLT"
                    )}
                  </button>

                  {suggestedVerbose && suggestedVerbose !== customVerbose && (
                    <button
                      onClick={() => setCustomVerbose(suggestedVerbose)}
                      className="px-4 py-2 bg-green-600 hover:bg-green-700 rounded text-white font-medium transition-colors"
                    >
                      Usar Sugerencia
                    </button>
                  )}
                </div>

                {/* Update Message */}
                {updateMessage.text && (
                  <div className={`p-3 rounded text-sm ${
                    updateMessage.type === 'success'
                      ? 'bg-green-900/20 border border-green-500 text-green-300'
                      : 'bg-red-900/20 border border-red-500 text-red-300'
                  }`}>
                    {updateMessage.text}
                  </div>
                )}
              </div>
            </div>
          ) : (
            <div className="bg-[#0f172a] rounded-lg border border-gray-700 p-6 text-center">
              <div className="text-gray-400">
                <svg className="w-12 h-12 mx-auto mb-3 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                <p className="text-sm">Selecciona un OLT de la lista para comenzar a actualizar su información.</p>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Empty State */}
      {!selectedOlt && olts.length > 0 && (
        <div className="bg-yellow-900/20 border border-yellow-500 text-yellow-300 p-4 rounded-lg">
          <p>Selecciona un OLT de la lista para comenzar a actualizar su información.</p>
        </div>
      )}

      {/* No Data State */}
      {olts.length === 0 && !isLoading && (
        <div className="bg-gray-800/50 border border-gray-600 text-gray-300 p-4 rounded-lg text-center">
          <p>No se encontraron OLTs para mostrar.</p>
        </div>
      )}
    </div>
  );
}
