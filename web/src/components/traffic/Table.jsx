import { useEffect, useState } from "react";
import { useStore } from "@nanostores/react";
import useFetch from "../../hooks/useFetch";
import {
  initDate,
  endDate,
  region,
  state,
  ip,
  gpon,
} from "../../stores/traffic";
import { isIpv4 } from "../../utils/validator";
import { formatSpeed, removeAccentsAndToUpper } from "../../utils/formater";

const BASE_URL = `${import.meta.env.PUBLIC_URL || ""}/api/traffic`;
const FAT_STATUS_URL = `${import.meta.env.PUBLIC_URL || ""}/api/fat/trend/detail`;
const TOKEN = sessionStorage.getItem("access_token").replace("Bearer ", "");

export default function Table() {
  const [header, setHeader] = useState(undefined);
  const [urlTraffic, setUrlTraffic] = useState(undefined)
  const [urlFat, setUrlFat] = useState(undefined)
  const $initDate = useStore(initDate);
  const $endDate = useStore(endDate);
  const $ip = useStore(ip);
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
      setHeader(
        <>
          <tr>
            <th rowSpan="2">Puerto</th>
            <th colSpan="2">Entrante (Mbps)</th>
            <th colSpan="2">Saliente (Mbps)</th>
            <th rowSpan="2">Capacidad (Mbps)</th>
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
    } else if ($ip) {
      setHeader(
        <>
          <tr>
            <th rowSpan="2">Puerto</th>
            <th colSpan="2">Entrante (Mbps)</th>
            <th colSpan="2">Saliente (Mbps)</th>
            <th rowSpan="2">Capacidad (Mbps)</th>
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
    } else if ($state) {
      setHeader(
        <>
          <tr>
            <th rowSpan="2">OLT</th>
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
    } else if ($region) {
      setHeader(
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
    }
  }, [dataTraffic, $region, $state, $ip, $gpon]);

  if (status === 401 || status === 403) {
    sessionStorage.removeItem("access_token");
    window.location.href = "/";
  }

  const filteredData =
    $gpon && dataTraffic
      ? dataTraffic.filter((row) => String(row.port) === String($gpon))
      : dataTraffic;

  if (filteredData && filteredData.length > 0) {
    return (
      <div className="w-full p-4 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
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
