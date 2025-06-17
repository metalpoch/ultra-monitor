import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";

export default function TrafficChartMbytes({ data }) {
  return (
    <div
      style={{
        width: "100%",
        height: 300 /* o cualquier altura fija o relativa */,
      }}
    >
      <ResponsiveContainer width="100%" height="100%">
        <LineChart width={500} height={300} data={data}>
          <XAxis dataKey="day" />
          <YAxis />
          <Tooltip />
          <Legend />
          <Line
            type="monotone"
            dataKey="mbytes_in_sec"
            name="Mbyte In"
            stroke="#3b82f6"
            fill="#3b82f6"
            strokeWidth="3"
          />
          <Line
            type="monotone"
            dataKey="mbytes_out_sec"
            name="Mbyte Out"
            stroke="#f59e0b"
            fill="#f59e0b"
            strokeWidth="3"
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
