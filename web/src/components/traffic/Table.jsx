import { useEffect, useState } from "react";
import { useStore } from "@nanostores/react";
import useFetch from "../../hooks/useFetch";
import { initDate, endDate, loadingChart, region, state, ip, urlTableData } from "../../stores/traffic";

const BASE_URL = `${import.meta.env.PUBLIC_URL}/api/traffic`;
const TOKEN = sessionStorage.getItem("access_token").replace("Bearer ", "");

export default function Table() {
  const [header, setHeader] = useState(undefined);
  const $initDate = useStore(initDate);
  const $endDate = useStore(endDate);
  const $url = useStore(urlTableData);
  const $ip = useStore(ip);
  const $state = useStore(state);
  const $region = useStore(region);
  const $loadingChart = useStore(loadingChart);

  const { data, status } = useFetch($url, { headers: { Authorization: `Bearer ${TOKEN}` } });

  useEffect(() => {
    if ($region) {
      const url = new URL(`${BASE_URL}/stats/region/${$region}`)
      url.searchParams.append("initDate", $initDate);
      url.searchParams.append("finalDate", $endDate);
      urlTableData.set(url.href)
    }
  }, [$region, $initDate, $endDate])

  useEffect(() => {
    if ($state) {
      const url = new URL(`${BASE_URL}/stats/state/${$state}`)
      url.searchParams.append("initDate", $initDate);
      url.searchParams.append("finalDate", $endDate);
      urlTableData.set(url.href)
    }
  }, [$state, $initDate, $endDate])

  useEffect(() => {
    if ($ip) {
      const url = new URL(`${BASE_URL}/stats/ip/${$ip}`)
      url.searchParams.append("initDate", $initDate);
      url.searchParams.append("finalDate", $endDate);
      urlTableData.set(url.href)
    }
  }, [$ip, $initDate, $endDate])


  useEffect(() => {
    if (data) {
      if ($region) setHeader(<>
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
      </>)

      if ($state) setHeader(<>
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
      </>)

      if ($ip) setHeader(<>
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
      </>)

    }
  }, [data])

  if (status === 401) {
    sessionStorage.removeItem("access_token")
    window.location.href = "/";
  }

  if (data && $url && !$loadingChart) return (
    <section className="w-full h-[300px] overflow-auto px-6">
      <table className="min-w-full">
        <thead className="sticky top-0 bg-[#121b31] pb-1">
          {header}
        </thead>
        <tbody>
          {data.map((row, idx) => <tr key={row.port ? row.port : idx} className="text-center">
            <td>{$region && row.state || $state && row.sys_name || $ip && row.if_name}</td>
            <td>{(row.avg_in_bps / 1_000_000).toFixed(2)}</td>
            <td>{(row.max_in_bps / 1_000_000).toFixed(2)}</td>
            <td>{(row.avg_out_bps / 1_000_000).toFixed(2)}</td>
            <td>{(row.max_out_bps / 1_000_000).toFixed(2)}</td>
            <td>{(row.if_speed / 1_000_000).toFixed(2)}</td>
            <td>{row.usage_in.toFixed(2)}%</td>
            <td>{row.usage_out.toFixed(2)}%</td>
          </tr>)}
        </tbody>
      </table>
    </section>
  )
}
