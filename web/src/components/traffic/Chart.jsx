import { useEffect, useState, useMemo } from "react";
import { useStore } from "@nanostores/react";
import TrafficChart from "./TrafficChart";
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
    ]
  );

  useEffect(() => {
    const params = new URLSearchParams();
    params.append("initDate", filters.initDate);
    params.append("finalDate", filters.endDate);

    if (filters.odn && filters.municipality && filters.state)
      setUrl(`${BASE_URL}/basic/odn/${MAP_STATE_TRANSLATER[filters.state]}/${filters.municipality}/${filters.odn}?${params.toString()}`);
    else if (filters.county && filters.municipality && filters.state)
      setUrl(`${BASE_URL}/basic/county/${MAP_STATE_TRANSLATER[filters.state]}/${filters.municipality}/${filters.county}?${params.toString()}`);
    else if (filters.municipality && filters.state)
      setUrl(`${BASE_URL}/basic/municipality/${MAP_STATE_TRANSLATER[filters.state]}/${filters.municipality}?${params.toString()}`);

    else if (isIpv4(filters.ip) && filters.gpon)
      setUrl(`${BASE_URL}/basic/index/${filters.ip}/${filters.gpon}?${params.toString()}`);
    else if (isIpv4(filters.ip))
      setUrl(`${BASE_URL}/basic/criteria/instance/${filters.ip}?${params.toString()}`);
    else if (filters.state)
      setUrl(`${BASE_URL}/basic/criteria/state/${filters.state}?${params.toString()}`);
    else if (filters.region)
      setUrl(`${BASE_URL}/basic/criteria/region/${filters.region}?${params.toString()}`);
    else setUrl(undefined)

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
      <section className="flex justify-center items-center flex-col flex-1 sm:flex-2 px-6 py-3 h-[100%] rounded-lg bg-[#121b31]">
        <div className="loader-wrapper">
          <span className="loader-letter">G</span>
          <span className="loader-letter">e</span>
          <span className="loader-letter">n</span>
          <span className="loader-letter">e</span>
          <span className="loader-letter">r</span>
          <span className="loader-letter">a</span>
          <span className="loader-letter">n</span>
          <span className="loader-letter">d</span>
          <span className="loader-letter">o</span>
          <span className="loader-letter">.</span>
          <div className="loader"></div>
        </div>
      </section>
    );
  }

  if (data && url)
    return (
      <section className="flex flex-col flex-1 sm:flex-2 px-6 py-3 rounded-lg bg-[#121b31]">
        <div>
          <div className="flex space-x-4 mb-4">
            <button
              className={`px-4 py-2 rounded-t-lg focus:outline-none ${activeTab === "traffic"
                ? "bg-[#1f2a48] font-semibold text-white"
                : "bg-[#121b31] text-slate-400 hover:text-white"
                }`}
              onClick={() => setActiveTab("traffic")}
            >
              Tráfico de Red
            </button>
            <button
              className={`px-4 py-2 rounded-t-lg focus:outline-none ${activeTab === "volume"
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
