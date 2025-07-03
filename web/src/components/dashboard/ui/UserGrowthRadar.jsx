import {
  Radar,
  Legend,
  RadarChart,
  Tooltip,
  PolarGrid,
  PolarAngleAxis,
  PolarRadiusAxis,
  ResponsiveContainer,
} from "recharts";
import { useStore } from "@nanostores/react";
import { groupFilteredStatusData } from "../../../utils/filterTrafficData";
import { selectedRegion } from "../../../stores/dashboard";
import { COLOR } from "../../../constants/colors";

export default function UserGrowthChart({ data }) {
  const $selectedRegion = useStore(selectedRegion);
  const groupedData = groupFilteredStatusData(data, $selectedRegion)
    .filter((item, _, arr) => item.day === arr[arr.length - 1].day)
    .sort((a, b) => a.actives - b.actives)
    .slice(0, 5);

  let radarChart = [];
  if (groupedData.length === 1) {
    const radarData = [
      { description: "Activos", value: groupedData[0].actives },
      { description: "Inactivos", value: groupedData[0].inactives },
      { description: "Fallas", value: groupedData[0].unknowns },
    ];
    radarChart = (
      <RadarChart cx="50%" cy="50%" outerRadius="80%" data={radarData}>
        <PolarGrid />
        <PolarAngleAxis dataKey="description" />
        <PolarRadiusAxis angle={90} tick={false} />
        <Radar
          name={groupedData[0].description}
          dataKey="value"
          stroke={COLOR[0]}
          fill={COLOR[0]}
          fillOpacity={0.6}
        />
        <Legend />
      </RadarChart>
    );
  } else {
    radarChart = (
      <RadarChart cx="50%" cy="50%" outerRadius="80%" data={groupedData}>
        <PolarGrid />
        <PolarAngleAxis dataKey="description" />
        <PolarRadiusAxis angle={90} tick={false} />
        <Tooltip
          contentStyle={{
            color: "#e0e6ed",
            backgroundColor: "#1a233a",
            border: "1px solid #2d3652",
          }}
        />
        <Radar
          name="Activos"
          dataKey="actives"
          stroke={COLOR[0]}
          fill={COLOR[0]}
          fillOpacity={0.6}
        />
        <Radar
          name="Inactivos"
          dataKey="inactives"
          stroke={COLOR[1]}
          fill={COLOR[1]}
          fillOpacity={0.6}
        />
        <Radar
          name="Fallas"
          dataKey="unknowns"
          stroke={COLOR[2]}
          fill={COLOR[2]}
          fillOpacity={0.6}
        />
        <Legend />
      </RadarChart>
    );
  }

  return (
    <div style={{ width: "100%", height: 175 }}>
      <ResponsiveContainer width="100%" height="100%">
        {radarChart}
      </ResponsiveContainer>
    </div>
  );
}

//
