import { useEffect, useState } from "react";
import { useStore } from "@nanostores/react";
import TrafficChart from "./ui/TrafficChart";
import {
  selectedLevel,
  selectedRegion,
  selectedState,
} from "../../stores/dashboard";
import useFetch from "../../hooks/useFetch";

export default function Traffic({ today, lastYear }) {
  const [url, setUrl] = useState("");
  const $selectedLevel = useStore(selectedLevel);
  const $selectedRegion = useStore(selectedRegion);
  const $selectedState = useStore(selectedState);

  const baseURL = `${import.meta.env.PUBLIC_API_URL}/traffic`;
  const { data, loading, error } = useFetch(url);
  useEffect(() => {
    const u = $selectedState
      ? new URL(`${baseURL}/state/${$selectedState}`)
      : $selectedRegion
      ? new URL(`${baseURL}/region/${$selectedRegion}`)
      : new URL(`${baseURL}/total`);

    u.searchParams.append("initDate", lastYear);
    u.searchParams.append("finalDate", today);
    setUrl(u.href);
  }, [$selectedLevel, $selectedRegion, $selectedState]);

  return (
    <>
      <section className="flex flex-col flex-1 sm:flex-2 px-6 py-3 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
        <h1 className="text-2xl font-semibold">Tráfico de Red</h1>
        {loading && <span className="mx-auto py-20 loader"></span>}
        {data && (
          <>
            <p className="text-slate-400 text-sm">
              Monitoreo del tráfico de entrada y salida.
            </p>
            <TrafficChart data={data} dataType="traffic" client:load />
          </>
        )}
      </section>
      <section className="flex flex-col flex-1 sm:flex-2 px-6 py-3 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
        <h1 className="text-2xl font-semibold">Volumen de la Red</h1>
        {loading && <span className="mx-auto py-20 loader"></span>}
        {data && (
          <>
            <p className="text-slate-400 text-sm">
              Monitoreo del volumen de datos de entrada y salida.
            </p>
            <TrafficChart data={data} dataType="volume" client:load />
          </>
        )}
      </section>
    </>
  );
}
