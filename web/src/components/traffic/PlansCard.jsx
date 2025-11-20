import { useStore } from "@nanostores/react";
import useFetch from "../../hooks/useFetch";
import { removeAccentsAndToUpper } from "../../utils/formater";
import { state, region, ip } from "../../stores/traffic";

const BASE_URL_FATS = `${import.meta.env.PUBLIC_URL || ""}/api/fat`;
const TOKEN = sessionStorage.getItem("access_token")?.replace("Bearer ", "") || "";

export default function PlansCard() {
  const $state = useStore(state);
  const $region = useStore(region);
  const $ip = useStore(ip);

  // Build URL based on current selection
  const buildFatUrl = () => {
    if ($ip) {
      return `${BASE_URL_FATS}/ip/${$ip}?page=1&limit=65535`;
    } else if ($state) {
      const formattedState = removeAccentsAndToUpper($state);
      return `${BASE_URL_FATS}/location/${formattedState}?page=1&limit=65535`;
    } else if ($region) {
      // Use the base endpoint with region filter to get plans data
      // Apply same formatting as states for database consistency
      const formattedRegion = removeAccentsAndToUpper($region);
      return `${BASE_URL_FATS}/?field=region&value=${formattedRegion}&page=1&limit=65535`;
    }
    return null;
  };

  const url = buildFatUrl();
  const { data: fatData } = useFetch(url, {
    headers: { Authorization: `Bearer ${TOKEN}` },
  });

  // Parse plans data from all FAT entries
  const parsePlansData = () => {
    if (!fatData || !Array.isArray(fatData)) return [];

    const planCounts = {};
    let totalUsers = 0;

    fatData.forEach(fat => {
      if (fat.plans && typeof fat.plans === 'string') {
        const planEntries = fat.plans.split(';');
        planEntries.forEach(entry => {
          const [countStr, plan] = entry.split('x');
          const count = parseInt(countStr, 10);
          if (!isNaN(count) && plan) {
            planCounts[plan] = (planCounts[plan] || 0) + count;
            totalUsers += count;
          }
        });
      }
    });

    // Convert to array and sort by count (descending)
    const sortedPlans = Object.entries(planCounts)
      .map(([plan, count]) => ({ plan, count, percentage: ((count / totalUsers) * 100).toFixed(1) }))
      .sort((a, b) => b.count - a.count);

    return { plans: sortedPlans, totalUsers };
  };

  const plansData = parsePlansData();
  const { plans = [], totalUsers = 0 } = plansData;

  if (!url || plans.length === 0) {
    return null;
  }

  return (
    <div className="w-full p-4 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
      <h3 className="text-lg font-semibold mb-4 text-white">
        Distribución de Usuarios por Plan
      </h3>

      <div className="mb-4">
        <div className="flex justify-between items-center text-sm text-gray-300">
          <span>Total de usuarios:</span>
          <span className="font-semibold text-white">{totalUsers.toLocaleString()}</span>
        </div>
      </div>

      <div className="space-y-3 max-h-64 overflow-y-auto">
        {plans.map(({ plan, count, percentage }) => (
          <div key={plan} className="flex items-center justify-between p-3 bg-[#1a233a] rounded-lg border border-[hsl(217,33%,25%)]">
            <div className="flex-1">
              <div className="flex justify-between items-center mb-1">
                <span className="text-sm font-medium text-white">{plan}</span>
                <span className="text-sm text-gray-300">{percentage}%</span>
              </div>
              <div className="w-full bg-[#121b31] rounded-full h-2">
                <div
                  className="bg-blue-500 h-2 rounded-full transition-all duration-300"
                  style={{ width: `${percentage}%` }}
                />
              </div>
            </div>
            <div className="ml-4 text-right">
              <span className="text-lg font-bold text-white">{count.toLocaleString()}</span>
              <div className="text-xs text-gray-400">usuarios</div>
            </div>
          </div>
        ))}
      </div>

      {plans.length > 5 && (
        <div className="mt-3 text-xs text-gray-400 text-center">
          Mostrando {plans.length} planes diferentes
        </div>
      )}
    </div>
  );
}

