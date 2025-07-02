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

  let filteredData = filterTrafficData(
    data,
    $selectedLevel,
    $selectedRegion,
    $selectedState
  );

  // Filtrar volumen al momento de tener datos validos
  filteredData = filteredData.filter(
    ({ day }) => convert.parseDate(day) >= convert.parseDate("01/07/2025")
  );

  if (!filteredData || filteredData.length < 2) {
    return (
      <>
        <article className="flex flex-1/3 md:flex-1/5 flex-col gap-2 px-6 py-3 w-1/3 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
          <div className="w-full flex justify-between">
            <h2 className="text-slate-400">Volumen Entrante Total</h2>
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
        <article className="flex flex-1/3 md:flex-1/5 flex-col gap-2 px-6 py-3 w-1/3 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
          <div className="w-full flex justify-between">
            <h2 className="text-slate-400">Volumen Saliente Total</h2>
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
      </>
    );
  }

  const init = filteredData[0];
  const end = filteredData[filteredData.length - 1];

  const volumePercentIn =
    ((end.mbytes_in - init.mbytes_in) * 100) / init.mbytes_in;
  const volumePercentOut =
    ((end.mbytes_out - init.mbytes_out) * 100) / init.mbytes_out;
  const volumeIn = `${convert.traffic(end.mbytes_in)}bps`;
  const volumeOut = `${convert.traffic(end.mbytes_out)}bps`;

  const strVolumeIn = `${
    volumePercentIn >= 0 ? "+" : ""
  }${volumePercentIn.toFixed(2)}%`;
  const strVolumeOut = `${
    volumePercentOut >= 0 ? "+" : ""
  }${volumePercentOut.toFixed(2)}%`;

  return (
    <>
      <article className="flex flex-1/3 md:flex-1/5 flex-col gap-2 px-6 py-3 w-1/3 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
        <div className="w-full flex justify-between">
          <h2 className="text-slate-400">Volumen Entrante Total</h2>
          <img
            src={boxIcon.src ?? boxIcon}
            width={20}
            height={20}
            alt="icono de pulso"
          />
        </div>
        <h3 className="font-bold text-3xl">{volumeIn}</h3>
        <p className="text-slate-400 text-sm">{strVolumeIn} de crecimiento</p>
        {volumePercentIn > 0 ? (
          <p className="text-green-300 text-sm">↑ {strVolumeIn}</p>
        ) : (
          <p className="text-red-300 text-sm">↓ {strVolumeIn}</p>
        )}
      </article>
      <article className="flex flex-1/3 md:flex-1/5 flex-col gap-2 px-6 py-3 w-1/3 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
        <div className="w-full flex justify-between">
          <h2 className="text-slate-400">Volumen Saliente Total</h2>
          <img
            src={boxIcon.src ?? boxIcon}
            width={20}
            height={20}
            alt="icono de pulso"
          />
        </div>
        <h3 className="font-bold text-3xl">{volumeOut}</h3>
        <p className="text-slate-400 text-sm">{strVolumeOut} de crecimiento</p>
        {volumePercentOut > 0 ? (
          <p className="text-green-300 text-sm">↑ {strVolumeOut}</p>
        ) : (
          <p className="text-red-300 text-sm">↓ {strVolumeOut}</p>
        )}
      </article>
    </>
  );
}
