import { useEffect, useState } from "react";
import { useStore } from "@nanostores/react";
import useFetch from "../../hooks/useFetch";
import {
  initDate,
  endDate,
  region,
  state,
  ip,
  switchValue,
  oltName,
  gpon,
} from "../../stores/traffic";
import { pdfHeaderConfig, getColumnCount } from "../../stores/pdfHeader";
import { isIpv4 } from "../../utils/validator";
import { formatSpeed, removeAccentsAndToUpper } from "../../utils/formater";
import { generateTablePDF, extractTableData } from "../../utils/pdfTableGenerator";
import { chartImage } from "../../stores/chart";
import dayjs from "dayjs";

const BASE_URL = `${import.meta.env.PUBLIC_URL || ""}/api/traffic`;
const FAT_STATUS_URL = `${import.meta.env.PUBLIC_URL || ""}/api/fat/trend/detail`;
const TOKEN = sessionStorage.getItem("access_token").replace("Bearer ", "");

export default function Table() {
  const [header, setHeader] = useState(undefined);
  const [urlTraffic, setUrlTraffic] = useState(undefined)
  const [urlFat, setUrlFat] = useState(undefined)
  const [isDownloading, setIsDownloading] = useState(false);
  const $initDate = useStore(initDate);
  const $endDate = useStore(endDate);
  const $ip = useStore(ip);
  const $oltName = useStore(oltName);
  const $switchValue = useStore(switchValue);
  const $gpon = useStore(gpon);
  const $state = useStore(state);
  const $region = useStore(region);

  useEffect(() => {
    const params = new URLSearchParams();
    params.append("initDate", $initDate);
    params.append("finalDate", $endDate);

    if (isIpv4($ip)) {
      setUrlTraffic(`${BASE_URL}/stats/ip/${$ip}?${params.toString()}`);
      setUrlFat(`${FAT_STATUS_URL}/ip/${$ip}`);
    } else if ($state) {
      setUrlTraffic(`${BASE_URL}/stats/state/${$state}?${params.toString()}`);
      setUrlFat(`${FAT_STATUS_URL}/state/${removeAccentsAndToUpper($state)}`);
    } else if ($region) {
      setUrlTraffic(`${BASE_URL}/stats/region/${$region}?${params.toString()}`);
      setUrlFat(`${FAT_STATUS_URL}/region/${$region}`);
    } else {
      // Clear URLs when no valid filters are selected
      setUrlTraffic(undefined);
      setUrlFat(undefined);
    }
  }, [$region, $state, $ip, $initDate, $endDate]);

  const { data: dataTraffic, status } = useFetch(urlTraffic, {
    headers: { Authorization: `Bearer ${TOKEN}` },
  });

  const { data: dataFat } = useFetch(urlFat, {
    headers: { Authorization: `Bearer ${TOKEN}` },
  });

  useEffect(() => {
    if (!dataTraffic) return;
    if ($gpon) {
      const headerElement = (
        <>
          <tr>
            <th rowSpan="2">Puerto</th>
            <th colSpan="2">Entrante</th>
            <th colSpan="2">Saliente </th>
            <th rowSpan="2">Capacidad</th>
            <th rowSpan="2">Uso %</th>
          </tr>
          <tr>
            <th>Prom.</th>
            <th>Max.</th>
            <th>Prom.</th>
            <th>Max.</th>
          </tr>
        </>
      );
      setHeader(headerElement);
      pdfHeaderConfig.set({
        headers: ["Puerto", "Prom. Entrante", "Max. Entrante", "Prom. Saliente", "Max. Saliente", "Capacidad", "Uso %"],
        columnCount: getColumnCount('gpon')
      });
    } else if ($ip) {
      const headerElement = (
        <>
          <tr>
            <th rowSpan="2">Puerto</th>
            <th colSpan="2">Entrante</th>
            <th colSpan="2">Saliente</th>
            <th rowSpan="2">Capacidad</th>
            <th rowSpan="2">Uso</th>
            <th colSpan="3">Estatus ONT</th>
          </tr>
          <tr>
            <th>Prom.</th>
            <th>Max.</th>
            <th>Prom.</th>
            <th>Max.</th>
            <th className="bg-green-900">Activo</th>
            <th className="bg-red-900">Cortado</th>
            <th className="bg-amber-900">En progreso</th>
          </tr>
        </>
      );
      setHeader(headerElement);
      pdfHeaderConfig.set({
        headers: ["Puerto", "Prom. Entrante", "Max. Entrante", "Prom. Saliente", "Max. Saliente", "Capacidad", "Uso", "Activo", "Cortado", "En progreso"],
        columnCount: getColumnCount('ip')
      });
    } else if ($state) {
      const headerElement = (
        <>
          <tr>
            <th rowSpan="2">OLT</th>
            <th rowSpan="2">Switch</th>
            <th rowSpan="2">Agregador</th>
            <th colSpan="2">Entrante</th>
            <th colSpan="2">Saliente</th>
            <th rowSpan="2">Capacidad</th>
            <th rowSpan="2">Uso</th>
            <th colSpan="3">Estatus ONT</th>
          </tr>
          <tr>
            <th>Prom.</th>
            <th>Max.</th>
            <th>Prom.</th>
            <th>Max.</th>
            <th className="bg-green-900">Activo</th>
            <th className="bg-red-900">Cortado</th>
            <th className="bg-amber-900">En progreso</th>
          </tr>
        </>
      );
      setHeader(headerElement);
      pdfHeaderConfig.set({
        headers: ["OLT", "Switch", "Agregador", "Prom. Entrante", "Max. Entrante", "Prom. Saliente", "Max. Saliente", "Capacidad", "Uso", "Activo", "Cortado", "En progreso"],
        columnCount: getColumnCount('state')
      });
    } else if ($region) {
      const headerElement = (
        <>
          <tr>
            <th rowSpan="2">Estado</th>
            <th colSpan="2">Entrante</th>
            <th colSpan="2">Saliente</th>
            <th rowSpan="2">Capacidad</th>
            <th rowSpan="2">Uso</th>
            <th colSpan="3">Estatus ONT</th>
          </tr>
          <tr>
            <th>Prom.</th>
            <th>Max.</th>
            <th>Prom.</th>
            <th>Max.</th>
            <th className="bg-green-900">Activo</th>
            <th className="bg-red-900">Cortado</th>
            <th className="bg-amber-900">En progreso</th>
          </tr>
        </>
      );
      setHeader(headerElement);
      pdfHeaderConfig.set({
        headers: ["Estado", "Prom. Entrante", "Max. Entrante", "Prom. Saliente", "Max. Saliente", "Capacidad", "Uso", "Activo", "Cortado", "En progreso"],
        columnCount: getColumnCount('region')
      });
    }
  }, [dataTraffic, $region, $state, $ip, $gpon]);

  if (status === 401 || status === 403) {
    sessionStorage.removeItem("access_token");
    window.location.href = "/";
  }

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
        switchValue: $switchValue,
        oltName: $oltName,
        gpon: $gpon,
        initDate: $initDate,
        endDate: $endDate
      };

      // Get chart image
      const chartImg = chartImage.get();

      // Generate PDF
      const doc = generateTablePDF(data, headers, filters, chartImg);

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

  const filteredData =
    $gpon && dataTraffic
      ? dataTraffic.filter((row) => String(row.port) === String($gpon))
      : dataTraffic;

  if (filteredData && filteredData.length > 0) {
    return (
      <div className="w-full p-4 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
        {/* PDF Download Button */}
        <div className="mb-4 flex justify-end">
          <button
            type="button"
            onClick={handleDownloadPDF}
            disabled={isDownloading}
            className={`
              px-4 py-2 rounded-lg font-medium transition-all duration-300
              flex items-center gap-2 text-sm
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
                  className="animate-spin h-4 w-4 text-white"
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
                Generando...
              </>
            ) : (
              <>
                <svg
                  className="h-4 w-4"
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

        <section className="w-full h-[300px] overflow-auto px-6">
          <table className="min-w-full text-sm">
            <thead className="sticky top-0 bg-[#121b31] pb-1">{header}</thead>
            <tbody>
              {filteredData.map((row, idx) => {
                let fatStatus = null
                let title = ""
                if ($gpon) {
                  title = row.if_name
                } else if ($ip) {
                  title = row.if_name
                  fatStatus = dataFat && dataFat.find((r) => row.if_name === r.name)
                } else if ($state) {
                  title = row.sys_name
                  fatStatus = dataFat && dataFat.find((r) => row.ip === r.name)
                }
                else {
                  title = row.state
                  fatStatus = dataFat && dataFat.find((r) => r.name === (title ? removeAccentsAndToUpper(title) : ""))
                }
                return <tr key={row.port ? row.port : idx} className="text-center">
                  <td>{title}</td>
                  {$state && !$ip && <td>{row.switch || ""}</td>}
                  {$state && !$ip && <td>{fatStatus ? fatStatus.bras : ""}</td>}
                  <td>{formatSpeed(row.avg_in_bps)}</td>
                  <td>{formatSpeed(row.max_in_bps)}</td>
                  <td>{formatSpeed(row.avg_out_bps)}</td>
                  <td>{formatSpeed(row.max_out_bps)}</td>
                  <td>{formatSpeed(row.if_speed)}</td>
                  <td>{(row.usage_out + row.usage_in).toFixed(2)}%</td>
                  {fatStatus && <>
                    <td className="bg-green-801 w-[80px]">{fatStatus.actives + fatStatus.provisioned_offline}</td>
                    <td className="bg-red-801 w-[80px]">{fatStatus.cut_off}</td>
                    <td className="bg-amber-801 w-[80px]">{fatStatus.in_progress}</td>

                  </>}
                </tr>
              })}
            </tbody>
          </table>
        </section>
      </div>
    );
  }

  return null;
}
