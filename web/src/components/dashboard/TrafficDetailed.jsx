import { useEffect, useState } from "react";
import { useStore } from "@nanostores/react";
import TrafficChartDetailed from "./ui/TrafficChartDetailed";
import {
  selectedLevel,
  selectedRegion,
  selectedState,
} from "../../stores/dashboard";
import useFetch from "../../hooks/useFetch";

export default function TrafficRegions({ today, lastYear }) {
  const [url, setUrl] = useState("");
  const [activeTab, setActiveTab] = useState("traffic");
  const $selectedLevel = useStore(selectedLevel);
  const $selectedRegion = useStore(selectedRegion);
  const $selectedState = useStore(selectedState);

  const baseURL = `${import.meta.env.PUBLIC_API_URL}/traffic`;
  const { data, loading, error } = useFetch(url);
  useEffect(() => {
    const u = $selectedState
      ? new URL(`${baseURL}/state/${$selectedState}`)
      : $selectedRegion
      ? new URL(`${baseURL}/states/${$selectedRegion}`)
      : new URL(`${baseURL}/regions`);

    u.searchParams.append("initDate", lastYear);
    u.searchParams.append("finalDate", today);
    setUrl(u.href);
  }, [$selectedLevel, $selectedRegion, $selectedState]);

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

      {data && (
        <>
          {activeTab === "traffic" && (
            <>
              <p className="text-slate-400 text-sm">
                Monitoreo del tráfico de entrada y salida.
              </p>
              <TrafficChartDetailed
                data={data}
                dataType="traffic"
                client:load
              />
            </>
          )}

          {activeTab === "volume" && (
            <>
              <p className="text-slate-400 text-sm">
                Monitoreo del volumen de datos de entrada y salida.
              </p>
              <TrafficChartDetailed data={data} dataType="volume" client:load />
            </>
          )}
        </>
      )}
    </section>
  );
}
