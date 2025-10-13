import { useState, useEffect } from "react";
import { useStore } from "@nanostores/react";
import {
  initDate,
  endDate,
  region,
  state,
  municipality,
  county,
  odn,
  ip,
  gpon,
} from "../../stores/traffic";
import Chart from "./Chart";

export default function ChartWrapper() {
  const [showChart, setShowChart] = useState(false);

  const $initDate = useStore(initDate);
  const $endDate = useStore(endDate);
  const $region = useStore(region);
  const $state = useStore(state);
  const $municipality = useStore(municipality);
  const $county = useStore(county);
  const $odn = useStore(odn);
  const $ip = useStore(ip);
  const $gpon = useStore(gpon);

  useEffect(() => {
    // Check if we have valid form data to show the chart
    const hasValidData =
      $initDate &&
      $endDate &&
      ($region || $state || $ip || $municipality || $county || $odn || $gpon);

    setShowChart(hasValidData);
  }, [
    $initDate,
    $endDate,
    $region,
    $state,
    $municipality,
    $county,
    $odn,
    $ip,
    $gpon,
  ]);

  return (
    <div className="flex-4 p-4 min-w-[500px] rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
      {showChart ? (
        <Chart />
      ) : (
        <div className="flex items-center justify-center h-full text-slate-400">
          <p>Selecciona filtros en el formulario para mostrar el gráfico</p>
        </div>
      )}
    </div>
  );
}

