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
    <div style={{ width: "100%", height: 400 }}>
      <ResponsiveContainer width="100%" height="100%">
        <LineChart width={500} height={500} data={data}>
          <XAxis
            tick={{ fill: "#bbb" }}
            dataKey="time"
            tickFormatter={(value) => dayjs(value).format("DD/MM/YY HH:mm")}
            angle={-7}
            textAnchor="end"
          />
          <YAxis
            tick={{ fill: "#bbb" }}
            tickFormatter={(value) =>
              value / 1e12 >= 1
                ? `${value / 1e12} ${
                    dataType === "traffic" ? "Tbps" : "Tbytes"
                  }`
                : value / 1e9 >= 1
                ? `${value / 1e9} ${dataType === "traffic" ? "Gbps" : "Gbytes"}`
                : value / 1e6 >= 1
                ? `${value / 1e6} ${dataType === "traffic" ? "Mbps" : "Mbytes"}`
                : value / 1e3 >= 1
                ? `${value / 1e3} ${dataType === "traffic" ? "Kbps" : "Kbytes"}`
                : `${value} ${dataType === "traffic" ? "bps" : "bytes"}`
            }
          />
          <Tooltip
            formatter={(value) =>
              value / 1e12 >= 1
                ? `${(value / 1e12).toFixed(2)} ${
                    dataType === "traffic" ? "Tbps" : "Tbytes"
                  }`
                : value / 1e9 >= 1
                ? `${(value / 1e9).toFixed(2)} ${
                    dataType === "traffic" ? "Gbps" : "Gbytes"
                  }`
                : value / 1e6 >= 1
                ? `${(value / 1e6).toFixed(2)} ${
                    dataType === "traffic" ? "Mbps" : "Mbytes"
                  }`
                : value / 1e3 >= 1
                ? `${(value / 1e3).toFixed(2)} ${
                    dataType === "traffic" ? "Kbps" : "Kbytes"
                  }`
                : `${value.toFixed(2)} ${
                    dataType === "traffic" ? "bps" : "bytes"
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
            stroke={dataType === "traffic" ? COLOR[9] : COLOR[1]}
            fill={dataType === "traffic" ? COLOR[9] : COLOR[1]}
            strokeWidth="1.5"
            dot={false}
          />
          <Line
            type="monotone"
            dataKey={dataType === "traffic" ? "bps_out" : "bytes_out"}
            name="Saliente"
            stroke={dataType === "traffic" ? COLOR[5] : COLOR[3]}
            fill={dataType === "traffic" ? COLOR[5] : COLOR[3]}
            strokeWidth="1.5"
            dot={false}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
