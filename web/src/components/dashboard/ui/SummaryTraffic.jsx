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

  const trafficPercentIn =
    ((end.mbps_in - init.mbps_in) * 100) / (init.mbps_in || 1);
  const trafficPercentOut =
    ((end.mbps_out - init.mbps_out) * 100) / (init.mbps_out || 1);

  const trafficIn = `${convert.traffic(end.mbps_in)}bps`;
  const trafficOut = `${convert.traffic(end.mbps_out)}bps`;

  const strTrafficIn = `${
    trafficPercentIn >= 0 ? "+" : ""
  }${trafficPercentIn.toFixed(2)}% `;

  const strTrafficOut = `${
    trafficPercentOut >= 0 ? "+" : ""
  }${trafficPercentOut.toFixed(2)}% `;

  return (
    <>
      <article className="flex flex-1/3 md:flex-1/5 flex-col gap-2 px-6 py-3 w-1/3 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
        <div className="w-full flex justify-between">
          <h2 className="text-slate-400">Tráfico Entrante Total</h2>
          <img
            src={pulseIcon.src ?? pulseIcon}
            width={20}
            height={20}
            alt="icono de pulso"
          />
        </div>
        <h3 className="font-bold text-3xl">{trafficIn}</h3>
        <p className="text-slate-400 text-sm">{strTrafficIn} de crecimiento</p>
        {trafficPercentIn > 0 ? (
          <p className="text-green-300 text-sm">↑ {strTrafficIn}</p>
        ) : (
          <p className="text-red-300 text-sm">↓ {strTrafficIn}</p>
        )}
      </article>

      <article className="flex flex-1/3 md:flex-1/5 flex-col gap-2 px-6 py-3 w-1/3 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
        <div className="w-full flex justify-between">
          <h2 className="text-slate-400">Tráfico Saliente Total</h2>
          <img
            src={pulseIcon.src ?? pulseIcon}
            width={20}
            height={20}
            alt="icono de pulso"
          />
        </div>
        <h3 className="font-bold text-3xl">{trafficOut}</h3>
        <p className="text-slate-400 text-sm">{strTrafficOut} de crecimiento</p>
        {trafficPercentOut > 0 ? (
          <p className="text-green-300 text-sm">↑ {strTrafficOut}</p>
        ) : (
          <p className="text-red-300 text-sm">↓ {strTrafficOut}</p>
        )}
      </article>
    </>
  );
}
