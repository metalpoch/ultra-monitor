import { useEffect, useState, useMemo } from "react";
import { useStore } from "@nanostores/react";
import TrafficChart from "../ui/TrafficChart";
import useFetch from "../../hooks/useFetch";
import { MAP_STATE_TRANSLATER } from "../../constants/regions";
import {
  initDate,
  endDate,
  region,
  state,
  municipality,
  county,
  odn,
  ip,
  gpon,
  oltsPrometheus,
} from "../../stores/traffic";
import { isIpv4 } from "../../utils/validator";

const BASE_URL = `${import.meta.env.PUBLIC_URL || ""}/api/traffic`;
const TOKEN =
  sessionStorage.getItem("access_token")?.replace("Bearer ", "") || "";

export default function Chart() {
  const [url, setUrl] = useState(undefined);
  const [activeTab, setActiveTab] = useState("traffic");
  const $initDate = useStore(initDate);
  const $endDate = useStore(endDate);
  const $region = useStore(region);
  const $state = useStore(state);
  const $municipality = useStore(municipality);
  const $county = useStore(county);
  const $odn = useStore(odn);
  const $ip = useStore(ip);
  const $gpon = useStore(gpon);
  const $oltsPrometheus = useStore(oltsPrometheus);

  // Memoiza los filtros para evitar renders innecesarios
  const filters = useMemo(
    () => ({
      initDate: $initDate,
      endDate: $endDate,
      region: $region,
      state: $state,
      municipality: $municipality,
      county: $county,
      odn: $odn,
      ip: $ip,
      gpon: $gpon,
      oltsPrometheus: $oltsPrometheus,
    }),
    [
      $initDate,
      $endDate,
      $region,
      $state,
      $municipality,
      $county,
      $odn,
      $ip,
      $gpon,
      $oltsPrometheus,
    ]
  );

  useEffect(() => {
    const params = new URLSearchParams();
    params.append("initDate", filters.initDate);
    params.append("finalDate", filters.endDate);

    // Prioridad: county > municipality > odn > gpon > ip > state > region
    if (filters.county && filters.state && filters.municipality) {
      setUrl(
        `${BASE_URL}/county/${MAP_STATE_TRANSLATER[filters.state]}/${
          filters.municipality
        }/${filters.county}?${params.toString()}`
      );
    } else if (filters.odn && filters.state && filters.municipality) {
      setUrl(
        `${BASE_URL}/odn/${MAP_STATE_TRANSLATER[filters.state]}/${
          filters.municipality
        }/${filters.odn}?${params.toString()}`
      );
    } else if (filters.municipality && filters.state) {
      setUrl(
        `${BASE_URL}/municipality/${MAP_STATE_TRANSLATER[filters.state]}/${
          filters.municipality
        }?${params.toString()}`
      );
    } else if (filters.gpon && filters.ip) {
      setUrl(
        `${BASE_URL}/index/${filters.ip}/${filters.gpon}?${params.toString()}`
      );
    } else if (isIpv4(filters.ip)) {
      params.append("ip", filters.ip);
      setUrl(`${BASE_URL}/instances?${params.toString()}`);
    } else if (filters.state && filters.oltsPrometheus) {
      filters.oltsPrometheus
        .filter(({ state }) => state === filters.state)
        .forEach(({ ip }) => params.append("ip", ip));
      setUrl(`${BASE_URL}/instances?${params.toString()}`);
    } else if (filters.region && filters.oltsPrometheus) {
      filters.oltsPrometheus
        .filter(({ region }) => region === filters.region)
        .forEach(({ ip }) => params.append("ip", ip));
      setUrl(`${BASE_URL}/instances?${params.toString()}`);
    } else {
      setUrl(undefined);
    }
  }, [filters]);

  const { data, status, loading, error } = useFetch(url, {
    headers: { Authorization: `Bearer ${TOKEN}` },
  });

  if (status === 401 || status === 403) {
    sessionStorage.removeItem("access_token");
    window.location.href = "/";
  }

  if (loading) {
    return (
      <section className="flex justify-center items-center flex-col flex-1 sm:flex-2 px-6 py-3 h-[400px] rounded-lg bg-[#121b31]">
        <span className="loader-table"></span>
        <h1 className="text-3xl">Buscando...</h1>
      </section>
    );
  }

  if (data && url)
    return (
      <section className="flex flex-col flex-1 sm:flex-2 px-6 py-3 rounded-lg bg-[#121b31]">
        <div>
          <div className="flex space-x-4 mb-4">
            <button
              className={`px-4 py-2 rounded-t-lg focus:outline-none ${
                activeTab === "traffic"
                  ? "bg-[#1f2a48] font-semibold text-white"
                  : "bg-[#121b31] text-slate-400 hover:text-white"
              }`}
              onClick={() => setActiveTab("traffic")}
            >
              Tráfico de Red
            </button>
            <button
              className={`px-4 py-2 rounded-t-lg focus:outline-none ${
                activeTab === "volume"
                  ? "bg-[#1f2a48] font-semibold text-white"
                  : "bg-[#121b31] text-slate-400 hover:text-white"
              }`}
              onClick={() => setActiveTab("volume")}
            >
              Volumen de la Red
            </button>
          </div>

          {activeTab === "traffic" && (
            <>
              <p className="text-slate-400 text-sm">
                Monitoreo del tráfico de entrada y salida total.
              </p>
              <TrafficChart data={data} dataType="traffic" client:load />
            </>
          )}

          {activeTab === "volume" && (
            <>
              <p className="text-slate-400 text-sm">
                Monitoreo del volumen de datos de entrada y salida.
              </p>
              <TrafficChart data={data} dataType="volume" client:load />
            </>
          )}
        </div>
      </section>
    );
}
