import { useStore } from "@nanostores/react";
import userIcon from "../../assets/icons/users.svg";
import {
  selectedLevel,
  selectedRegion,
  selectedState,
} from "../../stores/dashboard";
import convert from "../../utils/convert";
import { filterStatusData } from "../../utils/filterTrafficData";

export default function SummaryUsers({ data }) {
  const $selectedLevel = useStore(selectedLevel);
  const $selectedRegion = useStore(selectedRegion);
  const $selectedState = useStore(selectedState);

  const filteredData = filterStatusData(
    data,
    $selectedLevel,
    $selectedRegion,
    $selectedState
  );

  if (!filteredData || filteredData.length < 2) {
    return (
      <article className="flex flex-1/3 md:flex-1/5 flex-col gap-2 px-6 py-3 w-1/3 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
        <div className="w-full flex justify-between">
          <h2 className="text-slate-400">Usuarios Activos</h2>
          <img
            src={userIcon.src ?? userIcon}
            width={20}
            height={20}
            alt="icono de pulso"
          />
        </div>
        <h3 className="font-bold text-3xl">-</h3>
        <p className="text-slate-400 text-sm">Sin datos suficientes</p>
      </article>
    );
  }

  const init = filteredData[0];
  const totalInit = init.actives + init.actives;

  const end = filteredData[filteredData.length - 1];
  const totalEnd = end.actives + end.actives;

  const usersDiff = ((totalEnd - totalInit) * 100) / (totalInit || 1);
  const users = `${convert.qty(totalEnd)}`;
  const strUsers = `${usersDiff >= 0 ? "+" : ""}${usersDiff.toFixed(2)}%`;

  return (
    <article className="flex flex-1/3 md:flex-1/5 flex-col gap-2 px-6 py-3 w-1/3 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
      <div className="w-full flex justify-between">
        <h2 className="text-slate-400">Usuarios Activos</h2>
        <img
          src={userIcon.src ?? userIcon}
          width={20}
          height={20}
          alt="icono de pulso"
        />
      </div>
      <h3 className="font-bold text-3xl">{users}</h3>
      <p className="text-slate-400 text-sm">{strUsers} de crecimiento</p>
      {usersDiff > 0 ? (
        <p className="text-green-300 text-sm">↑ {strUsers}</p>
      ) : (
        <p className="text-red-300 text-sm">↓ {strUsers}</p>
      )}
    </article>
  );
}
