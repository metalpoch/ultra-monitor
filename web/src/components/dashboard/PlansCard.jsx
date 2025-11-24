import { useEffect, useState } from "react";
import { useStore } from "@nanostores/react";
import {
  selectedLevel,
  selectedRegion,
  selectedState,
} from "../../stores/dashboard";
import useFetch from "../../hooks/useFetch";
import { removeAccentsAndToUpper } from "../../utils/formater";
import { STATES_BY_REGION, MAP_STATE_TRANSLATER } from "../../constants/regions";

const BASE_URL_FATS = `${import.meta.env.PUBLIC_URL || ""}/api/fat`;
const BASE_URL_TRAFFIC = `${import.meta.env.PUBLIC_URL || ""}/api/traffic`;

export default function PlansCard() {
  const $selectedLevel = useStore(selectedLevel);
  const $selectedRegion = useStore(selectedRegion);
  const $selectedState = useStore(selectedState);

  // Fetch traffic data to get devices info
  const token = sessionStorage.getItem("access_token")?.replace("Bearer ", "") || "";
  const { data: trafficData } = useFetch(`${BASE_URL_TRAFFIC}/info`, {
    headers: { Authorization: `Bearer ${token}` },
  });

  // Build URLs based on current selection
  const [urls, setUrls] = useState([]);

  useEffect(() => {
    if ($selectedState) {
      const formattedState = removeAccentsAndToUpper($selectedState);
      setUrls([`${BASE_URL_FATS}/location/${formattedState}?page=1&limit=65535`]);
    } else if ($selectedRegion) {
      // For regions, fetch data for all states in the region
      if (trafficData && Array.isArray(trafficData)) {
        const regionStates = [...new Set(
          trafficData
            .filter(device => device.region === $selectedRegion)
            .map(device => device.state)
        )];

        const regionUrls = regionStates.map(state => {
          const formattedState = removeAccentsAndToUpper(state);
          return `${BASE_URL_FATS}/location/${formattedState}?page=1&limit=65535`;
        });
        setUrls(regionUrls);
      } else {
        setUrls([]);
      }
    } else {
      // National level - fetch data for all states from all regions
      const allStates = Object.values(STATES_BY_REGION).flat();
      const nationalUrls = allStates.map(state => {
        const formattedState = MAP_STATE_TRANSLATER[state] || removeAccentsAndToUpper(state);
        return `${BASE_URL_FATS}/location/${formattedState}?page=1&limit=65535`;
      });
      setUrls(nationalUrls);
    }
  }, [$selectedState, $selectedRegion, $selectedLevel, trafficData]);

  // State to store combined FAT data
  const [allFatData, setAllFatData] = useState([]);

  // Fetch data for all URLs
  useEffect(() => {
    if (urls.length === 0) {
      setAllFatData([]);
      return;
    }

    const fetchAllData = async () => {
      try {
        const responses = await Promise.all(
          urls.map(url =>
            fetch(url, {
              headers: { Authorization: `Bearer ${token}` }
            }).then(res => {
              if (!res.ok) {
                throw new Error(`HTTP error! status: ${res.status}`);
              }
              return res.json();
            })
          )
        );

        const combinedData = responses.reduce((acc, data) => {
          if (data && Array.isArray(data)) {
            return [...acc, ...data];
          }
          return acc;
        }, []);

        setAllFatData(combinedData);
      } catch (error) {
        console.error("Error fetching FAT data:", error);
        setAllFatData([]);
      }
    };

    fetchAllData();
  }, [urls]);

  // Parse plans data from all FAT entries
  const parsePlansData = () => {
    if (!allFatData || !Array.isArray(allFatData)) return [];

    const planCounts = {};
    let totalUsers = 0;

    allFatData.forEach(fat => {
      if (fat.plans && typeof fat.plans === 'string') {
        const planEntries = fat.plans.split(';');
        planEntries.forEach(entry => {
          const [countStr, plan] = entry.split('x');
          const count = parseInt(countStr, 10);
          if (!isNaN(count) && plan) {
            const cleanPlan = plan.trim();
            planCounts[cleanPlan] = (planCounts[cleanPlan] || 0) + count;
            totalUsers += count;
          }
        });
      }
    });

    // Convert to array and sort by count (descending)
    const sortedPlans = Object.entries(planCounts)
      .map(([plan, count]) => ({ plan, count, percentage: ((count / totalUsers) * 100).toFixed(2) }))
      .sort((a, b) => b.count - a.count);

    return { plans: sortedPlans, totalUsers };
  };

  const plansData = parsePlansData();
  const { plans = [], totalUsers = 0 } = plansData;

  // Only show if we have plans data
  if (plans.length === 0) {
    return null;
  }

  return (
    <div className="w-full h-full p-6 rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
      <h3 className="text-lg font-semibold mb-4 text-white">
        Distribución de Usuarios por Plan
      </h3>

      <div className="mb-4">
        <div className="flex justify-between items-center text-sm text-gray-300">
          <span>Total de usuarios:</span>
          <span className="font-semibold text-white">{totalUsers.toLocaleString()}</span>
        </div>
      </div>

      <div className="space-y-3 max-h-72 overflow-y-auto">
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

