import { useEffect, useState } from "react";
import { useStore } from "@nanostores/react";
import ClientsChart from "./ui/ClientsChart";
import {
  selectedLevel,
  selectedRegion,
  selectedState,
} from "../../stores/dashboard";
import { removeAccentsAndToUpper } from "../../utils/formater";
import useFetch from "../../hooks/useFetch";

export default function ClientStatus() {
  const [url, setUrl] = useState("");
  const $selectedLevel = useStore(selectedLevel);
  const $selectedRegion = useStore(selectedRegion);
  const $selectedState = useStore(selectedState);

  const baseURL = `${import.meta.env.PUBLIC_API_URL}/fat/trend/status`;
  const { data, loading, error } = useFetch(url);
  useEffect(() => {
    const u = $selectedState
      ? new URL(`${baseURL}/state/${removeAccentsAndToUpper($selectedState)}`)
      : $selectedRegion
      ? new URL(`${baseURL}/${$selectedRegion}`)
      : new URL(`${baseURL}`);

    setUrl(u.href);
  }, [$selectedLevel, $selectedRegion, $selectedState]);
  console.log({ data, loading, error, url });
  return (
    <>
      <section className="flex flex-col flex-1 sm:flex-2 px-6 py-3 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
        <h1 className="text-2xl font-semibold">Crecmiento de usuarios</h1>
        {loading && <span className="mx-auto py-20 loader"></span>}
        {data && (
          <>
            <p className="text-slate-400 text-sm">
              Monitoreo del comportamiento de los usuarios.
            </p>
            <ClientsChart data={data} client:load />
          </>
        )}
      </section>
    </>
  );
}
