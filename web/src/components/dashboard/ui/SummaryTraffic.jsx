import { useStore } from "@nanostores/react";
import pulseIcon from "../../../assets/icons/pulse.svg";
import {
  selectedLevel,
  selectedRegion,
  selectedState,
} from "../../../stores/dashboard";
import convert from "../../../utils/convert";
import { filterTrafficData } from "../../../utils/filterTrafficData";

export default function SummaryTraffic({ data }) {
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
            src={pulseIcon.src ?? pulseIcon}
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
  const end = filteredData[filteredData.length - 1];

  const totalInit = init.mbps_in + init.mbps_out;
  const totalEnd = end.mbps_in + end.mbps_out;

  const trafficPercent = ((totalEnd - totalInit) * 100) / (totalInit || 1);
  const traffic = `${convert.traffic(totalEnd)}bps`;
  const strTraffic = `${trafficPercent >= 0 ? "+" : ""}${trafficPercent.toFixed(
    2
  )}%`;

  return (
    <article className="flex flex-1/3 md:flex-1/5 flex-col gap-2 px-6 py-3 w-1/3 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
      <div className="w-full flex justify-between">
        <h2 className="text-slate-400">Tráfico Total</h2>
        <img
          src={pulseIcon.src ?? pulseIcon}
          width={20}
          height={20}
          alt="icono de pulso"
        />
      </div>
      <h3 className="font-bold text-3xl">{traffic}</h3>
      <p className="text-slate-400 text-sm">{strTraffic} de crecimiento</p>
      {trafficPercent > 0 ? (
        <p className="text-green-300 text-sm">↑ {strTraffic}</p>
      ) : (
        <p className="text-red-300 text-sm">↓ {strTraffic}</p>
      )}
    </article>
  );
}
