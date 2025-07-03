import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import { useStore } from "@nanostores/react";
import { filterTrafficData } from "../../../utils/filterTrafficData";
import { COLOR } from "../../../constants/colors";
import {
  selectedLevel,
  selectedRegion,
  selectedState,
} from "../../../stores/dashboard";

export default function TrafficChartBps({ data }) {
  const $selectedLevel = useStore(selectedLevel);
  const $selectedRegion = useStore(selectedRegion);
  const $selectedState = useStore(selectedState);

  const filteredData = filterTrafficData(
    data,
    $selectedLevel,
    $selectedRegion,
    $selectedState
  );

  return (
    <div style={{ width: "100%", height: 300 }}>
      <ResponsiveContainer width="100%" height="100%">
        <LineChart width={500} height={300} data={filteredData}>
          <XAxis dataKey="day" />
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
            dataKey="mbps_in"
            name="Mbps In"
            stroke={COLOR[0]}
            fill={COLOR[0]}
            strokeWidth="3"
          />
          <Line
            type="monotone"
            dataKey="mbps_out"
            name="Mbps Out"
            stroke={COLOR[1]}
            fill={COLOR[1]}
            strokeWidth="3"
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
