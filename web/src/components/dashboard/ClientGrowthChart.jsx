import { useEffect, useState } from "react";
import { useStore } from "@nanostores/react";
import ClientsChart from "./ui/ClientsChartChartJS";
import {
  selectedLevel,
  selectedRegion,
  selectedState,
} from "../../stores/dashboard";
import { removeAccentsAndToUpper } from "../../utils/formater";
import useFetch from "../../hooks/useFetch";

const BASE_URL_FAT = `${import.meta.env.PUBLIC_URL || ""}/api/fat/trend/status`;

export default function ClientGrowthChart() {
  const [urlFats, setUrlFats] = useState("");
  const $selectedLevel = useStore(selectedLevel);
  const $selectedRegion = useStore(selectedRegion);
  const $selectedState = useStore(selectedState);

  const token = sessionStorage.getItem("access_token").replace("Bearer ", "");

  useEffect(() => {
    let fURL;
    if ($selectedState) {
      fURL = `${BASE_URL_FAT}/state/${removeAccentsAndToUpper($selectedState)}`;
    } else if ($selectedRegion) {
      fURL = `${BASE_URL_FAT}/${$selectedRegion}`;
    } else {
      fURL = BASE_URL_FAT;
    }
    setUrlFats(fURL);
  }, [$selectedLevel, $selectedRegion, $selectedState]);

  const { data } = useFetch(urlFats, {
    headers: { Authorization: `Bearer ${token}` },
  });

  if (!data) return null;

  return (
    <section className="flex flex-col flex-1 sm:flex-2 px-6 py-3 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
      <h1 className="text-2xl font-semibold">Crecimiento de usuarios</h1>
      <p className="text-slate-400 text-sm">
        Monitoreo del comportamiento de los usuarios.
      </p>
      <ClientsChart data={data} client:load />
    </section>
  );
}
