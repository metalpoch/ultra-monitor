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
import { filterTrafficData } from "../../utils/filterTrafficData";
import {
  selectedLevel,
  selectedRegion,
  selectedState,
} from "../../stores/dashboard";

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
          <Tooltip />
          <Legend />
          <Line
            type="monotone"
            dataKey="mbps_in"
            name="Mbps In"
            stroke="#3b82f6"
            fill="#3b82f6"
            strokeWidth="3"
          />
          <Line
            type="monotone"
            dataKey="mbps_out"
            name="Mbps Out"
            stroke="#f59e0b"
            fill="#f59e0b"
            strokeWidth="3"
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
