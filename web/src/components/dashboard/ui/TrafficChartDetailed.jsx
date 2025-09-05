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

export default function TrafficChartDetailed({ data, dataType }) {
  const allTimes = [
    ...new Set(
      Object.values(data)
        .flat()
        .map((d) => d.time)
    ),
  ].sort();

  const chartData = allTimes.map((time) => {
    const entry = { time };
    for (const name of Object.keys(data)) {
      const point = data[name].find((d) => d.time === time);
      if (point) {
        if (dataType === "traffic") {
          entry[name] = Math.max(point.bps_in || 0, point.bps_out || 0);
        } else {
          entry[name] = Math.max(point.bytes_in || 0, point.bytes_out || 0);
        }
      } else {
        entry[name] = null;
      }
    }
    return entry;
  });

  return (
    <div style={{ width: "100%", height: 400 }}>
      <ResponsiveContainer width="100%" height="100%">
        <LineChart width={500} height={400} data={chartData}>
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
          />{" "}
          <Legend />
          {Object.keys(data).map((name, index) => (
            <Line
              key={name}
              type="monotone"
              dataKey={name}
              stroke={COLOR[index % COLOR.length]}
              strokeWidth="1.5"
              dot={false}
            />
          ))}
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
