import { useStore } from "@nanostores/react";
import dayjs from "dayjs";
import {
  dailyTrafficData,
  userData
} from "../../stores/dashboard";

// Función para formatear bytes
const formatBytes = (bytes) => {
  if (!bytes || bytes === 0 || isNaN(bytes)) return "0 B";
  const k = 1000;
  const sizes = ["B", "KB", "MB", "GB", "TB", "PB", "EB"];
  const i = Math.floor(Math.log(Math.abs(bytes)) / Math.log(k));
  const maxIndex = sizes.length - 1;
  const index = i > maxIndex ? maxIndex : i;
  const value = bytes / Math.pow(k, index);
  return `${value.toFixed(2)} ${sizes[index]}`;
};

// Función para calcular volúmenes por período
const calculateVolumeForPeriod = (trafficData, days) => {
  if (!trafficData || trafficData.length === 0) {
    return { incoming: 0, outgoing: 0 };
  }

  // Tomar los últimos N días
  const recentData = trafficData.slice(-days);

  // SUMAR los volume_in y volume_out
  const totalIncoming = recentData.reduce((sum, item) => sum + (item.volume_in || 0), 0);
  const totalOutgoing = recentData.reduce((sum, item) => sum + (item.volume_out || 0), 0);

  return { incoming: totalIncoming, outgoing: totalOutgoing };
};

// Función para calcular promedio por usuario
const calculateAveragePerUser = (volume, activeUsers) => {
  if (activeUsers === 0) return { incoming: "0 B", outgoing: "0 B", total: "0 B" };

  const incomingPerUser = volume.incoming / activeUsers;
  const outgoingPerUser = volume.outgoing / activeUsers;
  const totalPerUser = incomingPerUser + outgoingPerUser;

  return {
    incoming: formatBytes(incomingPerUser),
    outgoing: formatBytes(outgoingPerUser),
    total: formatBytes(totalPerUser)
  };
};

export default function VolumePerUserCards() {
  const $dailyTrafficData = useStore(dailyTrafficData);
  const $userData = useStore(userData);

  // Calcular volúmenes para cada período
  const annualVolume = calculateVolumeForPeriod($dailyTrafficData.data, 365); // Último año
  const monthlyVolume = calculateVolumeForPeriod($dailyTrafficData.data, 30);  // Último mes
  const weeklyVolume = calculateVolumeForPeriod($dailyTrafficData.data, 7);    // Última semana
  const dailyVolume = calculateVolumeForPeriod($dailyTrafficData.data, 1);     // Último día

  // Calcular promedios por usuario
  const annualAverage = calculateAveragePerUser(annualVolume, $userData.activeUsers);
  const monthlyAverage = calculateAveragePerUser(monthlyVolume, $userData.activeUsers);
  const weeklyAverage = calculateAveragePerUser(weeklyVolume, $userData.activeUsers);
  const dailyAverage = calculateAveragePerUser(dailyVolume, $userData.activeUsers);

  const cards = [
    {
      title: "Volumen Anual",
      period: "Último año (365 días)",
      volume: annualVolume,
      average: annualAverage
    },
    {
      title: "Volumen Mensual",
      period: "Último mes (30 días)",
      volume: monthlyVolume,
      average: monthlyAverage
    },
    {
      title: "Volumen Semanal",
      period: "Última semana (7 días)",
      volume: weeklyVolume,
      average: weeklyAverage
    },
    {
      title: "Volumen Diario",
      period: "Último día",
      volume: dailyVolume,
      average: dailyAverage
    }
  ];

  if ($dailyTrafficData.loading || $userData.loading) {
    return (
      <section className="w-full flex flex-wrap gap-4 p-4 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
        <div className="w-full text-center py-8">
          <div className="text-slate-400">Cargando datos de volumen...</div>
        </div>
      </section>
    );
  }

  return (
    <section className="w-full flex flex-wrap gap-4 p-4 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
      <div className="w-full mb-4">
        <h2 className="text-xl font-semibold text-white">Volumen por Usuario</h2>
      </div>

      <div className="w-full grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {cards.map((card, index) => (
          <div
            key={index}
            className="bg-[#1a233a] rounded-lg p-4 border border-[hsl(217,33%,20%)] hover:border-[hsl(217,33%,30%)] transition-colors"
          >
            <div className="mb-3">
              <h3 className="text-lg font-semibold text-white">{card.title}</h3>
              <p className="text-slate-400 text-sm">{card.period}</p>
            </div>

            <div className="space-y-3">
              <div>
                <div className="flex justify-between items-center mb-1">
                  <span className="text-slate-300 text-sm">Entrante Total</span>
                  <span className="text-white font-medium">{formatBytes(card.volume.incoming)}</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-slate-300 text-sm">Saliente Total</span>
                  <span className="text-white font-medium">{formatBytes(card.volume.outgoing)}</span>
                </div>
              </div>

              <div className="pt-2 border-t border-[hsl(217,33%,20%)]">
                <div className="flex justify-between items-center mb-1">
                  <span className="text-slate-300 text-sm">Promedio Entrante/Usuario</span>
                  <span className="text-green-400 font-medium">{card.average.incoming}</span>
                </div>
                <div className="flex justify-between items-center mb-1">
                  <span className="text-slate-300 text-sm">Promedio Saliente/Usuario</span>
                  <span className="text-blue-400 font-medium">{card.average.outgoing}</span>
                </div>
                <div className="flex justify-between items-center pt-1 border-t border-[hsl(217,33%,20%)]">
                  <span className="text-slate-300 text-sm font-medium">Total/Usuario</span>
                  <span className="text-white font-bold">{card.average.total}</span>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="w-full mt-4 text-xs text-slate-500">
        <p>* El promedio se calcula dividiendo el volumen total entre {$userData.activeUsers.toLocaleString()} usuarios activos</p>
      </div>
    </section>
  );
}

