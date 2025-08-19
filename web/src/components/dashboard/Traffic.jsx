import { useEffect, useState } from "react";
import { useStore } from "@nanostores/react";
import dayjs from "dayjs";
import TrafficChart from "../ui/TrafficChart";
import {
  selectedLevel,
  selectedRegion,
  selectedState,
} from "../../stores/dashboard";
import useFetch from "../../hooks/useFetch";

const BASE_URL = `${import.meta.env.PUBLIC_URL}/api/traffic`;

const minDate = dayjs("2025-07-01T00:00:00-04:00");
const today = dayjs()
  .set("hour", 0)
  .set("minute", 0)
  .set("second", 0)
  .set("millisecond", 0);
const lastYear =
  today.subtract(1, "year") < minDate ? minDate : today.subtract(1, "year");

export default function Traffic() {
  const [url, setUrl] = useState("");
  const [activeTab, setActiveTab] = useState("traffic");
  const $selectedLevel = useStore(selectedLevel);
  const $selectedRegion = useStore(selectedRegion);
  const $selectedState = useStore(selectedState);

  const token = sessionStorage.getItem("access_token").replace("Bearer ", "");
  const { data, status, loading } = useFetch(url, {
    headers: { Authorization: `Bearer ${token}` },
  });

  useEffect(() => {
    const u = $selectedState
      ? new URL(`${BASE_URL}/state/${$selectedState}`)
      : $selectedRegion
        ? new URL(`${BASE_URL}/region/${$selectedRegion}`)
        : new URL(`${BASE_URL}/total`);

    u.searchParams.append("initDate", lastYear.toISOString());
    u.searchParams.append("finalDate", today.toISOString());
    setUrl(u.href);
  }, [$selectedLevel, $selectedRegion, $selectedState]);

  if (status === 401) {
    sessionStorage.removeItem("access_token")
    window.location.href = "/";
  }

  if (loading) {
    return (
      <section className="flex flex-col flex-1 sm:flex-2 px-6 py-3 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
        <span className="mx-auto py-20 loader"></span>
      </section>
    );
  }

  return (
    <section className="flex flex-col flex-1 sm:flex-2 px-6 py-3 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
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

      {data && (
        <>
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
        </>
      )}
    </section>
  );
}
