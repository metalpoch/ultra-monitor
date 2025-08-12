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
import { trafficData, ontData } from "../../../stores/traffic";
import { dayField } from "../../../utils/delete.convert";

export default function TrafficChart() {
  const { data, status, loading, error, refetch } = useStore(trafficData);

  if (error) {
    alert(JSON.stringify(error));
    return;
  }
  const traffic = data && dayField(data);

  if (loading) return <span className="mx-auto py-20 loader"></span>;
  if (traffic)
    return (
      <section className="w-full h-[300px] flex flex-col flex-1 sm:flex-2 px-6 py-3 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
        <h1 className="text-2xl font-semibold">Tráfico de Red</h1>
        <p className="text-slate-400 text-sm">
          Monitoreo del tráfico de entrada y salida.
        </p>
        <ResponsiveContainer width="100%" height="100%">
          <LineChart width={500} height={300} data={traffic}>
            <XAxis dataKey="day" />
            <YAxis />
            <Tooltip />
            <Legend />
            <Line
              type="monotone"
              dataKey="mbps_in"
              name="Mbps In"
              stroke={COLOR[0]}
              fill={COLOR[0]}
              strokeWidth="3"
              dot={false}
            />
            <Line
              type="monotone"
              dataKey="mbps_out"
              name="Mbps Out"
              stroke={COLOR[1]}
              fill={COLOR[1]}
              strokeWidth="3"
              dot={false}
            />
          </LineChart>
        </ResponsiveContainer>
      </section>
    );
}
