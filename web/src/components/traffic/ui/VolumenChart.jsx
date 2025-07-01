import { useStore } from "@nanostores/react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import { COLOR } from "../../../constants/colors";
import { trafficData } from "../../../stores/traffic";
import { dayField } from "../../../utils/convert";

export default function VolumenChart() {
  const { data } = useStore(trafficData);
  const traffic = data && dayField(data);
  if (traffic)
    return (
      <section className="w-full h-[300px] flex flex-col flex-1 sm:flex-2 px-6 py-3 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
        <h1 className="text-2xl font-semibold">Volumen de Datos</h1>
        <p className="text-slate-400 text-sm">
          Monitoreo del volumen de datos transmitidos y recibidos.
        </p>
        <ResponsiveContainer width="100%" height="100%">
          <LineChart width={500} height={300} data={traffic}>
            <XAxis dataKey="day" />
            <YAxis />
            <Tooltip />
            <Legend />
            <Line
              type="monotone"
              dataKey="mbytes_in_sec"
              name="MBps In"
              stroke={COLOR[2]}
              fill={COLOR[2]}
              strokeWidth="3"
              dot={false}
            />
            <Line
              type="monotone"
              dataKey="mbytes_out_sec"
              name="MBps Out"
              stroke={COLOR[3]}
              fill={COLOR[3]}
              strokeWidth="3"
              dot={false}
            />
          </LineChart>
        </ResponsiveContainer>
      </section>
    );
}
