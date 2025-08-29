import { useEffect, useState } from "react";
import { useStore } from "@nanostores/react";
import useFetch from "../../hooks/useFetch";
import { initDate, endDate, ip, gponStats } from "../../stores/traffic";

const BASE_URL = `${import.meta.env.PUBLIC_URL}/api/traffic`;
const TOKEN = sessionStorage.getItem("access_token").replace("Bearer ", "");

export default function Table() {
  const [urlStats, setUrlStats] = useState(undefined);
  const $initDate = useStore(initDate);
  const $endDate = useStore(endDate);
  const $olt = useStore(ip);

  const { data, status, loading } = useFetch(urlStats, { headers: { Authorization: `Bearer ${TOKEN}` } });

  useEffect(() => {
    if ($olt) {
      const url = new URL(`${BASE_URL}/stats/${$olt}`)
      url.searchParams.append("initDate", $initDate);
      url.searchParams.append("finalDate", $endDate);
      setUrlStats(url.href)
    } else {
      setUrlStats(undefined)
    }
  }, [$olt, $initDate, $endDate])

  useEffect(() => {
    if (data) gponStats.set(data)
  }, [data])

  if (status === 401) {
    sessionStorage.removeItem("access_token")
    window.location.href = "/";
  }

  if (loading) {
    return (
      <section className="w-full">
        <span className="mx-auto py-20 loader"></span>
      </section>
    );
  }

  if (data && $olt) return (
    <section className="w-full h-[300px] overflow-auto px-6">
      <table className="min-w-full">
        <thead className="sticky top-0 bg-[#121b31] pb-1">
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
        </thead>
        <tbody>
          {data.map(row => <tr key={row.port} className="text-center">
            <td>{row.if_name}</td>
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
