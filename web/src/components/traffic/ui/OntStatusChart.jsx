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
import { ontData, trafficData } from "../../../stores/traffic";
import { dayField } from "../../../utils/delete.convert";

export default function OntStatusChart() {
  const { data } = useStore(ontData);
  const { loading } = useStore(trafficData);
  const ontStatus = data && dayField(data);
  if (ontStatus && !loading)
    return (
      <section className="w-full h-[300px] flex flex-col flex-1 sm:flex-2 px-6 py-3 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
        <h1 className="text-2xl font-semibold">Crecimiento de ONT</h1>
        <p className="text-slate-400 text-sm">
          Evolución de los ONT, activos, inactivos y fallas.
        </p>
        <ResponsiveContainer width="100%" height="100%">
          <LineChart width={500} height={300} data={ontStatus}>
            <XAxis dataKey="day" />
            <YAxis />
            <Tooltip />
            <Legend />
            <Line
              type="monotone"
              dataKey="actives"
              name="Activos"
              stroke={COLOR[5]}
              fill={COLOR[5]}
              strokeWidth="3"
              dot={false}
            />
            <Line
              type="monotone"
              dataKey="inactives"
              name="Inactivos"
              stroke={COLOR[6]}
              fill={COLOR[6]}
              strokeWidth="3"
              dot={false}
            />
            <Line
              type="monotone"
              dataKey="unknowns"
              name="Fallas"
              stroke={COLOR[7]}
              fill={COLOR[7]}
              strokeWidth="3"
              dot={false}
            />
          </LineChart>
        </ResponsiveContainer>
      </section>
    );
}
