import { useStore } from "@nanostores/react";
import { groupFilteredStatusData } from "../../../utils/filterTrafficData";
import { selectedRegion, selectedState } from "../../../stores/dashboard";

export default function UserGrowthTable({ data }) {
  const $selectedRegion = useStore(selectedRegion);
  const $selectedState = useStore(selectedState);

  const filteredData = groupFilteredStatusData(
    data,
    $selectedRegion,
    $selectedState
  );

  return (
    <div className="w-5/12 text-slate-300  rounded-xl shadow p-4">
      <div className="mb-2 text-sm text-slate-300">
        Estado diario de usuarios por <span className="font-bold">estado</span>
      </div>
      <div className="w-full  overflow-auto h-80">
        <table className="min-w-full text-sm h-[300px] ">
          <thead>
            <tr className="bg-slate-700 text-slate-300">
              <th className="px-4 py-2 text-left font-semibold  rounded-tl-lg">
                Fecha
              </th>
              <th className="px-4 py-2 text-left font-semibold text-blue-400">
                {$selectedRegion === "" ? "Estado" : "Municipio"}
              </th>
              <th className="px-4 py-2 text-left font-semibold text-blue-400">
                Usuarios Activos
              </th>
              <th className="px-4 py-2 text-left font-semibold text-yellow-400">
                Usuarios Inactivos
              </th>
              <th className="px-4 py-2 text-left font-semibold text-red-400 rounded-tr-lg">
                Usuarios Fallas
              </th>
            </tr>
          </thead>
          <tbody className="border-b-2 border-slate-700 ">
            {filteredData.map((item, index) => (
              <tr key={index}>
                <td className="px-4 py-2 rounded-l-lg">{item.day}</td>
                <td className="px-4 py-2 rounded-l-lg">{item.description}</td>
                <td className="px-4 py-2 ">{item.actives}</td>
                <td className="px-4 py-2 ">{item.inactives}</td>
                <td className="px-4 py-2  rounded-r-lg">{item.unknowns}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
