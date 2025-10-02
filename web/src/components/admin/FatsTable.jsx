import { useState, useEffect } from "react";
import useFetch from "../../hooks/useFetch";
import DatalistField from "../ui/DatalistField";
import SelectField from "../ui/SelectField";

const BASE_URL = `${import.meta.env.PUBLIC_URL || ""}/api/fat/`;
const TOKEN = sessionStorage.getItem("access_token")?.replace("Bearer ", "") || "";

export default function FatsTable() {
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(10);
  const [searchField, setSearchField] = useState("");
  const [searchValue, setSearchValue] = useState("");
  const [datalistOptions, setDatalistOptions] = useState([]);

  const buildUrl = () => {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString()
    });

    if (searchValue.trim()) {
      params.append("field", searchField);
      params.append("value", searchValue.trim());
    }

    return `${BASE_URL}?${params.toString()}`;
  };

  const [url, setUrl] = useState(buildUrl());

  const { data, loading, error, status } = useFetch(url, {
    headers: { Authorization: `Bearer ${TOKEN}` }
  });

  const [isSearching, setIsSearching] = useState(false);
  const [currentData, setCurrentData] = useState([]);

  // Update current data when new data is loaded
  useEffect(() => {
    if (!loading && data) {
      setCurrentData(data);
    }
  }, [data, loading]);

  // Trigger search when manually requested or when pagination/search parameters change
  useEffect(() => {
    setIsSearching(true);
    setUrl(buildUrl());
  }, [page, limit, searchValue]);

  // Reset searching state when data is loaded
  useEffect(() => {
    if (!loading) {
      setIsSearching(false);
    }
  }, [loading]);

  // Reset to page 1 when search field changes
  useEffect(() => {
    setPage(1);
    setSearchValue("")
  }, [searchField]);

  // Fetch datalist options when search field changes
  useEffect(() => {
    const fetchOptions = async () => {
      if (searchField) {
        try {
          const fetchOptionsUrl = `${BASE_URL}options/${searchField}`;
          const response = await fetch(fetchOptionsUrl, {
            headers: { Authorization: `Bearer ${TOKEN}` }
          });
          if (response.ok) {
            const data = await response.json();
            setDatalistOptions(data);
          } else {
            console.error("Error response:", response.status, response.statusText);
          }
        } catch (err) {
          console.error("Error fetching options:", err);
        }
      }
    };

    fetchOptions();
  }, [searchField]);

  if (status === 401 || status === 403) {
    sessionStorage.removeItem("access_token");
    window.location.href = "/";
  }

  return (
    <div className="space-y-4">
      {/* Search Controls */}
      <form className="flex gap-4 items-center">
        <SelectField
          id="fat-field"
    label="Buscar por"
          onChange={(e) => setSearchField(e.target.value)}
          value={searchField}
          options={[
            { value: "", label: "Seleccionar", disabled: true, hidden: true },
            { value: "ip", label: "IP" },
            { value: "region", label: "Región" },
            { value: "state", label: "Estado" },
            { value: "municipality", label: "Municipio" },
            { value: "county", label: "Parroquia" },
            { value: "odn", label: "ODN" },
          ]}

        />
        {searchField &&
          <DatalistField
            id="fat-search"
    label="Filtrar por..."
            value={searchValue}
            onChange={(e) => setSearchValue(e.target.value)}
            options={datalistOptions}
            className="flex-1"
          />
        }
      </form>

      {/* Pagination Controls */}
      <div className="flex justify-between items-center">
        <div className="flex items-center gap-4">
          <span>Mostrar:</span>
          <select
            value={limit}
            onChange={(e) => setLimit(Number(e.target.value))}
            className="px-3 py-2 bg-slate-800 border border-slate-600 rounded-md"
          >
            <option value={10}>10</option>
            <option value={25}>25</option>
            <option value={50}>50</option>
            <option value={100}>100</option>
          </select>
          <span>registros por página</span>
        </div>

        <div className="flex gap-2">
          <button
            onClick={() => setPage(Math.max(1, page - 1))}
            disabled={page === 1}
            className="px-3 py-2 bg-slate-700 hover:bg-slate-600 disabled:opacity-50 rounded-md"
          >
            Anterior
          </button>
          <span className="px-3 py-2">Página {page}</span>
          <button
            onClick={() => setPage(page + 1)}
            className="px-3 py-2 bg-slate-700 hover:bg-slate-600 rounded-md"
          >
            Siguiente
          </button>
        </div>
      </div>

      {/* Error State */}
      {error ? (
        <div className="text-red-500 text-center py-4">Error: {error.message}</div>
      ) : null}

      {/* Loading State */}
      {(loading || isSearching) && (
        <div className="mt-4 flex justify-center">
          <div className="mx-auto py-20 text-slate-400">Cargando datos...</div>
        </div>
      )}

      {/* Table - Only show when not loading and has data */}
      {!loading && !isSearching && currentData && currentData.length > 0 && (
        <div className="min-w-full overflow-x-auto">
          <table className="min-w-full text-sm table-auto border-collapse">
            <thead className="sticky top-0 bg-[#121b31] pb-1 text-left">
              <tr>
                <th>IP</th>
                <th className="px-2">Región</th>
                <th>Estado</th>
                <th className="px-2">Municipio</th>
                <th>Parroquia</th>
                <th className="px-2">ODN</th>
                <th>FAT</th>
                <th className="px-2">Puerto</th>
                <th>Activos</th>
                <th className="px-2">Offline</th>
                <th>Corte</th>
                <th className="px-2">En Progreso</th>
                <th>Fecha Creación</th>
              </tr>
            </thead>
            <tbody>
              {currentData.map((fat) => (
                <tr key={fat.id} className="duration-150 hover:bg-slate-800">
                  <td>{fat.ip}</td>
                  <td className="px-2">{fat.region}</td>
                  <td>{fat.state}</td>
                  <td className="px-2">{fat.municipality}</td>
                  <td>{fat.county}</td>
                  <td className="px-2">{fat.odn}</td>
                  <td>{fat.fat}</td>
                  <td className="px-2">GPON {fat.shell}/{fat.card}/{fat.port}</td>
                  <td>{fat.actives}</td>
                  <td className="px-2">{fat.provisioned_offline}</td>
                  <td>{fat.cut_off}</td>
                  <td className="px-2">{fat.in_progress}</td>
                  <td>{new Date(fat.created_at).toLocaleDateString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Empty State - Show when not loading but no data */}
      {!loading && !isSearching && (!currentData || currentData.length === 0) && (
        <div className="text-center py-20 text-slate-400">
          No se encontraron resultados
        </div>
      )}
    </div>
  );
}
