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

export default function TrafficChartMbytes({ data }) {
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
            dataKey="mbytes_in_sec"
            name="Mbyte In"
            stroke={COLOR[0]}
            fill={COLOR[0]}
            strokeWidth="3"
          />
          <Line
            type="monotone"
            dataKey="mbytes_out_sec"
            name="Mbyte Out"
            stroke={COLOR[1]}
            fill={COLOR[1]}
            strokeWidth="3"
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
