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

const BASE_URL_FAT = `${import.meta.env.PUBLIC_URL || ""}/api/fat/trend/status`;
const BASE_URL_STATUS = `${import.meta.env.PUBLIC_URL || ""}/api/prometheus/status`;

export default function ClientStatus() {
  const [urlFats, setUrlFats] = useState("");
  const [urlGponStatus, setUrlGponStatus] = useState("");
  const [dataTable, setDataTable] = useState(undefined);
  const $selectedLevel = useStore(selectedLevel);
  const $selectedRegion = useStore(selectedRegion);
  const $selectedState = useStore(selectedState);

  const token = sessionStorage.getItem("access_token").replace("Bearer ", "");

  useEffect(() => {
    let fURL, sURL

    if ($selectedState) {
      fURL = `${BASE_URL_FAT}/state/${removeAccentsAndToUpper($selectedState)}`
      sURL = `${BASE_URL_STATUS}/state/${$selectedState}`;
    } else if ($selectedRegion) {
      fURL = `${BASE_URL_FAT}/${$selectedRegion}`;
      sURL = `${BASE_URL_STATUS}/region/${$selectedRegion}`;
    } else {
      fURL = BASE_URL_FAT;
      sURL = BASE_URL_STATUS;
    }
    setUrlFats(fURL);
    setUrlGponStatus(sURL);
  }, [$selectedLevel, $selectedRegion, $selectedState]);

  const { data } = useFetch(urlFats, {
    headers: { Authorization: `Bearer ${token}` },
  });
  const { data: dataGpon } = useFetch(urlGponStatus, {
    headers: { Authorization: `Bearer ${token}` },
  });

  useEffect(() => {
    if (data && dataGpon) {
      const lastData = data[data.length - 1];
      const currOntAct = lastData.actives + lastData.provisioned_offline;
      const currOntInact = lastData.cut_off + lastData.in_progress;
      const currOntTotal = currOntAct + currOntInact;

      const ontPortInstalled =
        (currOntTotal / dataGpon.total_gpon).toFixed(2) + " %";
      const ontPortIluminated =
        (currOntTotal / dataGpon.gpon_actives).toFixed(2) + " %";

      const table = [
        { name: "OLT en medición con usuarios", value: dataGpon.olts },
        { name: "Cantidad de tarjetas", value: dataGpon.cards },
        { name: "GPON iluminados", value: dataGpon.gpon_actives },
        { name: "GPON instalados", value: dataGpon.total_gpon },
        { name: "ONTs definidos *", value: currOntTotal },
        { name: "ONTs activos *", value: currOntAct },
        { name: "ONTs inactivos *", value: currOntInact },
        {
          name: "ONTs/GPON instalados (densidad)",
          value: ontPortInstalled,
        },
        {
          name: "ONTs/GPON Iluminado (densidad)",
          value: ontPortIluminated,
        },
      ];

      setDataTable(table);
    }
  }, [data, dataGpon]);

  return (
    <>
      {data && (
        <section className="flex flex-col flex-1 sm:flex-2 px-6 py-3 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
          <h1 className="text-2xl font-semibold">Crecimiento de usuarios</h1>
          <>
            <p className="text-slate-400 text-sm">
              Monitoreo del comportamiento de los usuarios.
            </p>
            <ClientsChart data={data} client:load />
          </>
        </section>
      )}
      {dataTable && (
        <section className="overflow-x-auto rounded-md border-2 border-[hsl(217,33%,20%)] bg-[#121b31]">
          <h1 className="text-2xl px-4 pt-2 font-semibold">Inventario GPON</h1>
          <table className="min-w-full divide-y divide-[hsl(217,33%,20%)]">
            <thead>
              <tr>
                <th className="px-6 py-3 text-left text-sm font-semibold text-slate-300"></th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-slate-300">
                  Cantidad
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-[hsl(217,33%,20%)] text-slate-200">
              {dataTable.map(({ name, value }) => (
                <tr key={name} className="hover:bg-[rgba(255,255,255,0.05)]">
                  <td className="px-6 py-2">{name}</td>
                  <td className="px-6 py-2">{value}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>
      )}
    </>
  );
}
