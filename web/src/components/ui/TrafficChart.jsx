import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import { COLOR } from "../../constants/colors";
import dayjs from "dayjs";

export default function TrafficChart({ data, dataType }) {
  return (
    <div style={{ width: "100%", height: 300 }}>
      <ResponsiveContainer width="100%" height="100%">
        <LineChart width={500} height={300} data={data}>
          <XAxis
            dataKey="time"
            tickFormatter={(value) => dayjs(value).format("DD/MM/YY HH:mm")}
            angle={-7}
            textAnchor="end"
          />
          <YAxis
            tickFormatter={(value) =>
              `${value / 1e9} ${dataType === "traffic" ? "Gbps" : "Gbytes"}`
            }
          />
          <Tooltip
            formatter={(value) =>
              `${(value / 1e9).toFixed(2)} ${dataType === "traffic" ? "Gbps" : "Gbytes"
              }`
            }
            contentStyle={{
              color: "#e0e6ed",
              backgroundColor: "#1a233a",
              border: "1px solid #2d3652",
            }}
          />
          <Legend />
          <Line
            type="monotone"
            dataKey={dataType === "traffic" ? "bps_in" : "bytes_in"}
            name="Entrante"
            stroke={dataType === "traffic" ? COLOR[0] : COLOR[1]}
            fill={dataType === "traffic" ? COLOR[0] : COLOR[1]}
            strokeWidth="1.5"
            dot={false}
          />
          <Line
            type="monotone"
            dataKey={dataType === "traffic" ? "bps_out" : "bytes_out"}
            name="Saliente"
            stroke={dataType === "traffic" ? COLOR[2] : COLOR[3]}
            fill={dataType === "traffic" ? COLOR[2] : COLOR[3]}
            strokeWidth="1.5"
            dot={false}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
