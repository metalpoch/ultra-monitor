import { useStore } from "@nanostores/react";
import boxIcon from "../../../assets/icons/box.svg";
import {
  selectedLevel,
  selectedRegion,
  selectedState,
} from "../../../stores/dashboard";
import convert from "../../../utils/convert";

import { filterTrafficData } from "../../../utils/filterTrafficData";

export default function SummaryVolume({ data }) {
  const $selectedLevel = useStore(selectedLevel);
  const $selectedRegion = useStore(selectedRegion);
  const $selectedState = useStore(selectedState);

  const filteredData = filterTrafficData(
    data,
    $selectedLevel,
    $selectedRegion,
    $selectedState
  );

  if (!filteredData || filteredData.length < 2) {
    return (
      <article className="flex flex-1/3 md:flex-1/5 flex-col gap-2 px-6 py-3 w-1/3 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
        <div className="w-full flex justify-between">
          <h2 className="text-slate-400">Tráfico Total</h2>
          <img
            src={boxIcon.src ?? boxIcon}
            width={20}
            height={20}
            alt="icono de caja"
          />
        </div>
        <h3 className="font-bold text-3xl">-</h3>
        <p className="text-slate-400 text-sm">Sin datos suficientes</p>
      </article>
    );
  }

  const init = filteredData[0];
  const end = filteredData[filteredData.length - 1];

  const totalInit = init.mbytes_in_sec + init.mbytes_out_sec;
  const totalEnd = end.mbytes_in_sec + end.mbytes_out_sec;

  const volumePercent = ((totalEnd - totalInit) * 100) / totalInit;
  const volume = `${convert.traffic(totalEnd)}bps`;
  const strVolume = `${volumePercent >= 0 ? "+" : ""}${volumePercent.toFixed(
    2
  )}%`;

  return (
    <>
      <article className="flex flex-1/3 md:flex-1/5 flex-col gap-2 px-6 py-3 w-1/3 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
        <div className="w-full flex justify-between">
          <h2 className="text-slate-400">Volumen Total</h2>
          <img
            src={boxIcon.src ?? boxIcon}
            width={20}
            height={20}
            alt="icono de pulso"
          />
        </div>
        <h3 className="font-bold text-3xl">{volume}</h3>
        <p className="text-slate-400 text-sm">{strVolume} de crecimiento</p>
        {volumePercent > 0 ? (
          <p className="text-green-300 text-sm">↑ {strVolume}</p>
        ) : (
          <p className="text-red-300 text-sm">↓ {strVolume}</p>
        )}
      </article>
    </>
  );
}
