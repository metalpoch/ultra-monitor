import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";

export default function UserGrowthChart({ data }) {
  return (
    <div
      style={{
        width: "100%",
        height: 300 /* o cualquier altura fija o relativa */,
      }}
    >
      <ResponsiveContainer width="100%" height="100%">
        <AreaChart data={data}>
          <XAxis dataKey="date" />
          <YAxis />
          <Tooltip />
          <Legend />
          <Area
            type="monotone"
            dataKey="actives"
            name="Activos"
            stroke="#3b82f6"
            fill="#3b82f6"
            strokeWidth="3"
          />
          <Area
            type="monotone"
            dataKey="inactives"
            name="Inactivos"
            stroke="#f59e0b"
            fill="#f59e0b"
            strokeWidth="3"
          />
          <Area
            type="monotone"
            dataKey="unknowns"
            name="Fallas"
            stroke="red"
            fill="red"
            strokeWidth="3"
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
}
