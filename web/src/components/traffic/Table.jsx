import { useEffect, useState } from "react";
import { useStore } from "@nanostores/react";
import useFetch from "../../hooks/useFetch";
import {
  initDate,
  endDate,
  oltsPrometheus,
  region,
  state,
  ip,
  gpon,
  urlTableData,
} from "../../stores/traffic";
import { isIpv4 } from "../../utils/validator";

const BASE_URL = `${import.meta.env.PUBLIC_URL || ""}/api/traffic`;
const TOKEN = sessionStorage.getItem("access_token").replace("Bearer ", "");

export default function Table() {
  const [header, setHeader] = useState(undefined);
  const $initDate = useStore(initDate);
  const $endDate = useStore(endDate);
  const $ip = useStore(ip);
  const $gpon = useStore(gpon);
  const $state = useStore(state);
  const $region = useStore(region);
  const $dataChart = useStore(oltsPrometheus);

  useEffect(() => {
    const params = new URLSearchParams();
    params.append("initDate", $initDate);
    params.append("finalDate", $endDate);

    if (isIpv4($ip)) {
      urlTableData.set(`${BASE_URL}/stats/ip/${$ip}?${params.toString()}`);
    } else if ($state) {
      urlTableData.set(
        `${BASE_URL}/stats/state/${$state}?${params.toString()}`
      );
    } else if ($region) {
      urlTableData.set(
        `${BASE_URL}/stats/region/${$region}?${params.toString()}`
      );
    }
  }, [$region, $state, $ip, $initDate, $endDate]);

  const $url = useStore(urlTableData);
  const { data, status } = useFetch($url, {
    headers: { Authorization: `Bearer ${TOKEN}` },
  });

  useEffect(() => {
    if (!data) return;
    if ($gpon) {
      setHeader(
        <>
          <tr>
            <th rowSpan="2">Puerto</th>
            <th colSpan="2">Entrante (Mbps)</th>
            <th colSpan="2">Saliente (Mbps)</th>
            <th rowSpan="2">Capacidad (Mbps)</th>
            <th colSpan="2">Uso %</th>
          </tr>
          <tr>
            <th>Prom.</th>
            <th>Max.</th>
            <th>Prom.</th>
            <th>Max.</th>
            <th>Entrante</th>
            <th>Saliente</th>
          </tr>
        </>
      );
    } else if ($state) {
      setHeader(
        <>
          <tr>
            <th rowSpan="2">OLT</th>
            <th colSpan="2">Entrante (Mbps)</th>
            <th colSpan="2">Saliente (Mbps)</th>
            <th rowSpan="2">Capacidad (Mbps)</th>
            <th colSpan="2">Uso %</th>
          </tr>
          <tr>
            <th>Prom.</th>
            <th>Max.</th>
            <th>Prom.</th>
            <th>Max.</th>
            <th>Entrante</th>
            <th>Saliente</th>
          </tr>
        </>
      );
    } else if ($region) {
      setHeader(
        <>
          <tr>
            <th rowSpan="2">Estado</th>
            <th colSpan="2">Entrante (Mbps)</th>
            <th colSpan="2">Saliente (Mbps)</th>
            <th rowSpan="2">Capacidad (Mbps)</th>
            <th colSpan="2">Uso %</th>
          </tr>
          <tr>
            <th>Prom.</th>
            <th>Max.</th>
            <th>Prom.</th>
            <th>Max.</th>
            <th>Entrante</th>
            <th>Saliente</th>
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
            <th colSpan="2">Uso %</th>
          </tr>
          <tr>
            <th>Prom.</th>
            <th>Max.</th>
            <th>Prom.</th>
            <th>Max.</th>
            <th>Entrante</th>
            <th>Saliente</th>
          </tr>
        </>
      );
    }
  }, [data, $region, $state, $ip, $gpon]);

  if (status === 401) {
    sessionStorage.removeItem("access_token");
    window.location.href = "/";
  }

  // Filtra por GPON si corresponde
  const filteredData =
    $gpon && data
      ? data.filter((row) => String(row.port) === String($gpon))
      : data;

  if (filteredData && filteredData.length > 0 && $dataChart.length > 0) {
    return (
      <section className="w-full h-[300px] overflow-auto px-6">
        <table className="min-w-full">
          <thead className="sticky top-0 bg-[#121b31] pb-1">{header}</thead>
          <tbody>
            {filteredData.map((row, idx) => (
              <tr key={row.port ? row.port : idx} className="text-center">
                <td>
                  {($gpon && row.if_name) ||
                    ($region && row.state) ||
                    ($state && row.sys_name) ||
                    ($ip && row.if_name)}
                </td>
                <td>{(row.avg_in_bps / 1_000_000).toFixed(2)}</td>
                <td>{(row.max_in_bps / 1_000_000).toFixed(2)}</td>
                <td>{(row.avg_out_bps / 1_000_000).toFixed(2)}</td>
                <td>{(row.max_out_bps / 1_000_000).toFixed(2)}</td>
                <td>{(row.if_speed / 1_000_000).toFixed(2)}</td>
                <td>{row.usage_in.toFixed(2)}%</td>
                <td>{row.usage_out.toFixed(2)}%</td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>
    );
  }

  return null;
}
