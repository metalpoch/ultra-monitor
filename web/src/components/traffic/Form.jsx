import { useEffect, useState } from "react";
import dayjs from "dayjs";
import SelectField from "../ui/SelectField";
import DateField from "../ui/DateField";
import DatalistField from "../ui/DatalistField";
import RadioGroup from "../ui/RadioGroup";
import useFetch from "../../hooks/useFetch";
import { removeAccentsAndToUpper } from "../../utils/formater";
import { generateTablePDF, extractTableData } from "../../utils/pdfTableGenerator";
import {
  initDate,
  endDate,
  region,
  state,
  ip,
  oltName,
  municipality,
  county,
  odn,
  gpon,
} from "../../stores/traffic";
import { pdfHeaderConfig } from "../../stores/pdfHeader";
import { useStore } from "@nanostores/react";
import { isIpv4 } from "../../utils/validator";

const BASE_URL_TRAFFIC = `${import.meta.env.PUBLIC_URL || ""}/api/traffic`;
const BASE_URL_FATS = `${import.meta.env.PUBLIC_URL || ""}/api/fat`;
const TOKEN = sessionStorage.getItem("access_token")?.replace("Bearer ", "") || ""
endDate.set(dayjs().toJSON());
initDate.set(dayjs().subtract(1, "week").startOf("day").toJSON());

export default function Form() {
  const [urlFatState, setUrlFatState] = useState(undefined);
  const [urlOlt, setUrlOlt] = useState(undefined);
  const [regions, setRegions] = useState([])
  const [states, setStates] = useState([])
  const [selectionMethod, setSelectionMethod] = useState("");
  const [isDownloading, setIsDownloading] = useState(false);

  const $initDate = useStore(initDate);
  const $endDate = useStore(endDate);
  const $region = useStore(region);
  const $state = useStore(state);
  const $municipality = useStore(municipality);
  const $county = useStore(county);
  const $ip = useStore(ip);
  const $oltName = useStore(oltName);
  const $gpon = useStore(gpon);
  const $odn = useStore(odn);

  // Reset all form states when component mounts
  useEffect(() => {
    endDate.set(dayjs().toISOString());
    initDate.set(dayjs().subtract(1, "week").startOf("day").toJSON());
    region.set("");
    state.set("");
    municipality.set("");
    county.set("");
    ip.set("");
    oltName.set("");
    gpon.set("");
    odn.set("");
    setSelectionMethod("");
  }, []);

  const headers = { headers: { Authorization: `Bearer ${TOKEN}` } }

  // Load all OLT from prometheus data
  const { data, status, loading, error } = useFetch(`${BASE_URL_TRAFFIC}/hierarchy?initDate=${$initDate}&finalDate=${$endDate}`, headers);

  // Load specific OLT info when URL is set
  const { data: infoOlt } = useFetch(urlOlt, headers);

  // Load specific OLT info when URL is set
  const { data: fatsByState } = useFetch(urlFatState, headers);

  // Get regions array and ensure OLT data is properly loaded
  useEffect(() => {
    if (data) setRegions(data.regions)
  }, [data]);

  const handleDateChange = ({ init, end }) => {
    const hour = dayjs().hour()
    const minute = dayjs().minute()
    if (init) initDate.set(dayjs(init).toJSON());
    if (end) endDate.set(dayjs(end).hour(hour).minute(minute).toJSON());
  };

  const handleChangeRegion = ({ target }) => {
    setStates(data.states[target.value])
    region.set(target.value);
    state.set("");
    municipality.set("");
    county.set("");
    ip.set("");
    oltName.set("");
    gpon.set("");
    odn.set("");
    setSelectionMethod("");
  };

  const handleChangeState = ({ target }) => {
    const formatedState = removeAccentsAndToUpper(target.value);
    setUrlFatState(
      `${BASE_URL_FATS}/location/${formatedState}?page=1&limit=65535`
    );

    state.set(target.value);
    municipality.set("");
    county.set("");
    ip.set("");
    oltName.set("");
    gpon.set("");
    odn.set("");
    setSelectionMethod("");
  };

  const handleChangeMethod = (method) => {
    setSelectionMethod(method);
    municipality.set("");
    county.set("");
    ip.set("");
    oltName.set("");
    gpon.set("");
    odn.set("");
  };

  const handleChangeOlt = ({ target }) => {
    if (isIpv4(target.value)) {
      setUrlOlt(`${BASE_URL_TRAFFIC}/info/instance/${target.value}`);
    }

    const selectedObject = data.olts[$state].find(({ ip }) => ip === target.value);

    ip.set(target.value);
    oltName.set(selectedObject.sys_name);
    gpon.set("");
    municipality.set("");
    county.set("");
    odn.set("");
  };

  const handleChangeMunicipality = ({ target }) => {
    municipality.set(target.value);
    county.set("");
    odn.set("");
    ip.set("");
    gpon.set("");
  };

  const handleChangeCounty = ({ target }) => {
    county.set(target.value);
    ip.set("");
    oltName.set("");
    gpon.set("");
    odn.set("");
  };

  const handleChangeOdn = ({ target }) => {
    odn.set(target.value);
  };

  const handleChangeGpon = ({ target }) => {
    gpon.set(target.value);
  };

  const handleDownloadPDF = async () => {
    if (isDownloading) return;

    setIsDownloading(true);

    try {
      // Find the table element
      const tableElement = document.querySelector('table');
      if (!tableElement) {
        alert('No hay datos de tabla disponibles para exportar');
        return;
      }

      // Extract table data only (no headers from DOM)
      const { data } = extractTableData(tableElement);

      if (data.length === 0) {
        alert('No hay datos disponibles para exportar');
        return;
      }

      // Get headers from the shared store
      const headerConfig = pdfHeaderConfig.get();
      const headers = headerConfig.headers;

      if (headers.length === 0) {
        alert('No hay configuración de encabezados disponible');
        return;
      }

      // Get current filter values
      const filters = {
        region: $region,
        state: $state,
        ip: $ip,
        gpon: $gpon,
        initDate: $initDate,
        endDate: $endDate
      };

      // Generate PDF
      const doc = generateTablePDF(data, headers, filters);

      // Generate filename
      const timestamp = dayjs().format('YYYY-MM-DD_HH-mm-ss');
      let filename = `traffic_report_${timestamp}.pdf`;

      if ($region) filename = `traffic_${$region}_${timestamp}.pdf`;
      if ($state) filename = `traffic_${$state}_${timestamp}.pdf`;
      if ($ip) filename = `traffic_${$ip.replace(/\./g, '_')}_${timestamp}.pdf`;

      // Save PDF
      doc.save(filename);

    } catch (error) {
      console.error('Error generating PDF:', error);
      alert('Error al generar el PDF. Por favor, intente nuevamente.');
    } finally {
      setIsDownloading(false);
    }
  };

  if (status === 401 || status === 403) {
    sessionStorage.removeItem("access_token");
    window.location.href = "/";
  }

  return (
    <form>
      <DateField
        id="initDate"
        label="Fecha inicial *"
        value={$initDate}
        max={$endDate ? $endDate.split("T")[0] : undefined}
        onChange={(init) => handleDateChange({ init })}
      />

      <DateField
        id="endDate"
        label="Fecha final *"
        endOfDay={true}
        value={$endDate}
        min={$initDate ? $initDate.split("T")[0] : undefined}
        onChange={(end) => handleDateChange({ end })}
      />

      <SelectField
        id="region"
        label="Región *"
        options={[
          {
            value: "",
            label: "Seleccionar región",
            disabled: true,
            hidden: true,
          },
          ...regions.map(region => ({ value: region, label: region })),
        ]}
        value={$region}
        onChange={handleChangeRegion}
      />

      {$initDate && $endDate && $region && (
        <SelectField
          id="state"
          label="Estado *"
          options={[
            {
              value: "",
              label: "Seleccionar estado",
              disabled: true,
              hidden: true,
            },
            ...states.map(state => ({ value: state, label: state })), ,
          ]}
          value={$state}
          onChange={handleChangeState}
        />
      )}

      {$state && (
        <RadioGroup
          id="selectionMethod"
          label="Método de busqueda"
          options={[
            { value: "olt", label: "OLT" },
            { value: "municipality", label: "Ubicación" },
          ]}
          value={selectionMethod}
          onChange={handleChangeMethod}
        />
      )}

      {$state && selectionMethod === "olt" && data && data.olts && data.olts[$state] && (
        <>
          <DatalistField
            id="olt"
            label="OLT"
            options={data.olts[$state].map(({ ip: value, sys_name: label }) => ({ value, label }))}
            value={$ip}
            onChange={handleChangeOlt}
            placeholder="Ingrese el OLT"
          />
          {$ip && infoOlt && infoOlt.length > 0 && (
            <SelectField
              id="gpon"
              label="GPON *"
              options={[
                {
                  value: "",
                  label: "Seleccionar Puerto GPON",
                  disabled: true,
                  hidden: true,
                },
                ...infoOlt.map(({ if_name, if_index }) => ({
                  value: if_index,
                  label: if_name,
                })),
              ]}
              value={$gpon}
              onChange={handleChangeGpon}
            />
          )}
        </>
      )}

      {$state && selectionMethod === "municipality" && fatsByState && (
        <>
          <SelectField
            id="municipality"
            label="Municipio *"
            options={[
              {
                value: "",
                label: "Seleccionar Municipio",
                disabled: true,
                hidden: true,
              },
              ...[...new Set(fatsByState.map((f) => f.municipality))].map(
                (m) => ({ value: m, label: m })
              ),
            ]}
            value={$municipality}
            onChange={handleChangeMunicipality}
          />

          {$municipality && (
            <SelectField
              id="county"
              label="Parroquia *"
              options={[
                {
                  value: "",
                  label: "Seleccionar Parroquia",
                  disabled: true,
                  hidden: true,
                },
                ...[
                  ...new Set(
                    fatsByState
                      .filter((f) => f.municipality === $municipality)
                      .map(({ county }) => county)
                  ),
                ].map((c) => ({ value: c, label: c })),
              ]}
              value={$county}
              onChange={handleChangeCounty}
            />
          )}

          {$county && (
            <SelectField
              id="odn"
              label="ODN *"
              options={[
                {
                  value: "",
                  label: "Seleccionar ODN",
                  disabled: true,
                  hidden: true,
                },
                ...[
                  ...new Set(
                    fatsByState
                      .filter(
                        (f) =>
                          f.municipality === $municipality &&
                          f.county === $county
                      )
                      .map(({ odn }) => odn)
                  ),
                ].map((o) => ({ value: o, label: o })),
              ]}
              value={$odn}
              onChange={handleChangeOdn}
            />
          )}
        </>
      )}

      {/* PDF Download Button */}
      <div className="mt-6 flex justify-center">
        <button
          type="button"
          onClick={handleDownloadPDF}
          disabled={isDownloading}
          className={`
            px-6 py-3 rounded-lg font-medium transition-all duration-300
            flex items-center gap-2
            ${isDownloading
              ? 'bg-blue-600 cursor-not-allowed'
              : 'bg-blue-500 hover:bg-blue-600 active:scale-95'
            }
            text-white shadow-lg hover:shadow-xl
          `}
        >
          {isDownloading ? (
            <>
              <svg
                className="animate-spin h-5 w-5 text-white"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle
                  className="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  strokeWidth="4"
                ></circle>
                <path
                  className="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
              Generando PDF...
            </>
          ) : (
            <>
              <svg
                className="h-5 w-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                ></path>
              </svg>
              Descargar PDF
            </>
          )}
        </button>
      </div>
    </form>
  );
}
