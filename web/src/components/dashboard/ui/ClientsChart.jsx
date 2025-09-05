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
import dayjs from "dayjs";

export default function ClientsChart({ data }) {
  return (
    <div style={{ width: "100%", height: 400 }}>
      <ResponsiveContainer width="100%" height="100%">
        <LineChart width={500} height={400} data={data}>
          <XAxis
            dataKey="date"
            tickFormatter={(value) => dayjs(value).format("DD/MM/YY")}
            angle={-7}
            textAnchor="end"
          />
          <YAxis />
          <Tooltip
            contentStyle={{
              color: "#e0e6ed",
              backgroundColor: "#1a233a",
              border: "1px solid #2d3652",
            }}
          />
          <Legend />
          <Line
            type="monotone"
            dataKey="actives"
            name="Activos"
            stroke={COLOR[0]}
            fill={COLOR[0]}
            strokeWidth="1.5"
            dot={false}
          />
          <Line
            type="monotone"
            dataKey="provisioned_offline"
            name="Activos offline"
            stroke={COLOR[1]}
            fill={COLOR[1]}
            strokeWidth="1.5"
            dot={false}
          />
          <Line
            type="monotone"
            dataKey="cut_off"
            name="Cortados"
            stroke={COLOR[2]}
            fill={COLOR[2]}
            strokeWidth="1.5"
            dot={false}
          />
          <Line
            type="monotone"
            dataKey="in_progress"
            name="En proceso"
            stroke={COLOR[3]}
            fill={COLOR[3]}
            strokeWidth="1.5"
            dot={false}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
