import { useStore } from "@nanostores/react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { COLOR } from "../../../constants/colors";
import { trafficData } from "../../../stores/traffic";
import { dayField } from "../../../utils/convert";
import availability from "../../../assets/icons/availability.svg";

export default function BandwidthPercent() {
  const { data } = useStore(trafficData);
  const traffic = data && dayField(data);

  if (!data) {
    return;
  }

  const max = Math.max(
    ...traffic.map(({ bandwidth_mbps_sec }) => bandwidth_mbps_sec)
  );
  const avg =
    traffic.reduce((prev, curr) => curr && prev + curr.bandwidth_mbps_sec, 0) /
    traffic.length;
  const percent = (avg * 100) / max;

  return (
    <article className="flex flex-1 flex-col gap-2 px-6 py-3 w-1/3 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
      <div className="w-full flex justify-between ">
        <h2 className="text-slate-400 ">Disponibilidad</h2>
        <img
          src={availability.src ?? availability}
          width={20}
          height={20}
          alt="icono de pulso"
        />
      </div>
      <div className="h-full flex flex-col items-center justify-center">
        <h3 className="font-bold text-3xl">{percent.toFixed(2)}%</h3>
      </div>
    </article>
  );
}
