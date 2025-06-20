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
    <div
      style={{
        width: "100%",
        height: 300,
      }}
    >
      <ResponsiveContainer width="100%" height="100%">
        <AreaChart data={filteredData}>
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
