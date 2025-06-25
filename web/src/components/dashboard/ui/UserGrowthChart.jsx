import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import { useStore } from "@nanostores/react";
import { filterStatusData } from "../../../utils/filterTrafficData";
import { COLOR } from "../../../constants/colors";
import {
  selectedLevel,
  selectedRegion,
  selectedState,
} from "../../../stores/dashboard";

export default function UserGrowthChart({ data }) {
  const $selectedLevel = useStore(selectedLevel);
  const $selectedRegion = useStore(selectedRegion);
  const $selectedState = useStore(selectedState);

  const filteredData = filterStatusData(
    data,
    $selectedLevel,
    $selectedRegion,
    $selectedState
  );

  return (
    <div style={{ width: "100%", height: 300 }}>
      <ResponsiveContainer width="100%" height="100%">
        <AreaChart data={filteredData}>
          <XAxis dataKey="day" />
          <YAxis />
          <Tooltip />
          <Legend />
          <Area
            type="monotone"
            dataKey="actives"
            name="Activos"
            stroke={COLOR[0]}
            fill={COLOR[0]}
            strokeWidth="3"
          />
          <Area
            type="monotone"
            dataKey="inactives"
            name="Inactivos"
            stroke={COLOR[1]}
            fill={COLOR[1]}
            strokeWidth="3"
          />
          <Area
            type="monotone"
            dataKey="unknowns"
            name="Fallas"
            stroke={COLOR[2]}
            fill={COLOR[2]}
            strokeWidth="3"
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
}
