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

export default function BandwidthChart() {
  const { data, status } = useStore(trafficData);
  const traffic = data && dayField(data);

  if (traffic)
    return (
      <section className="w-full h-[300px] flex flex-col flex-6  px-6 py-3 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
        <h1 className="text-2xl font-semibold">Ancho de banda</h1>
        <p className="text-slate-400 text-sm">
          Monitoreo del ancho de banda, refleja la disponibilidad.
        </p>
        <ResponsiveContainer width="100%" height="100%">
          <LineChart width={500} height={300} data={traffic}>
            <XAxis dataKey="day" />
            <YAxis />
            <Tooltip />
            <Line
              type="monotone"
              dataKey="bandwidth_mbps_sec"
              name="Bandwidth"
              stroke={COLOR[4]}
              strokeWidth="3"
              dot={false}
            />
          </LineChart>
        </ResponsiveContainer>
      </section>
    );
}
