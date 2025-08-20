import { useEffect, useState } from "react";
import { useStore } from "@nanostores/react";
import TrafficChart from "../ui/TrafficChart";
import useFetch from "../../hooks/useFetch";
import {
  initDate,
  endDate,
  olt,
  odn,
  fat,
} from "../../stores/traffic";

const BASE_URL = `${import.meta.env.PUBLIC_URL}/api/traffic/fat`;

export default function TrafficState() {
  const [urlTraffic, setUrlTraffic] = useState("");
  const [activeTab, setActiveTab] = useState("traffic");
  const $initDate = useStore(initDate);
  const $endDate = useStore(endDate);
  const $olt = useStore(olt);
  const $odn = useStore(odn);
  const $fat = useStore(fat);

  const token = sessionStorage.getItem("access_token").replace("Bearer ", "");
  const { data, status, loading } = useFetch(urlTraffic, {
    headers: { Authorization: `Bearer ${token}` },
  });
  useEffect(() => {
    if ($initDate && $endDate && $fat) {
      const url = new URL(`${BASE_URL}/${$olt}/${$odn}/${$fat}`)
      url.searchParams.append("initDate", $initDate);
      url.searchParams.append("finalDate", $endDate);
      setUrlTraffic(url.href)
    }
  }, [$fat, $initDate, $endDate]);

  if (status === 401) {
    sessionStorage.removeItem("access_token")
    window.location.href = "/";
  }

  if (loading) {
    return (
      <section className="flex flex-col flex-1 sm:flex-2 px-6 py-3 h-[400px] rounded-lg rounded-t-none border-t-0 bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
        <span className="mx-auto py-20 loader"></span>
      </section>
    );
  }

  if ($endDate) return (
    <section className="flex flex-col flex-1 sm:flex-2 px-6 py-3 rounded-lg rounded-t-none border-t-0 bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
      {data && (
        <>
          <div className="flex space-x-4 mb-4">
            <button
              className={`px-4 py-2 rounded-t-lg focus:outline-none ${activeTab === "traffic"
                ? "bg-[#1f2a48] font-semibold text-white"
                : "bg-[#121b31] text-slate-400 hover:text-white"
                }`}
              onClick={() => setActiveTab("traffic")}
            >
              Tráfico de Red
            </button>
            <button
              className={`px-4 py-2 rounded-t-lg focus:outline-none ${activeTab === "volume"
                ? "bg-[#1f2a48] font-semibold text-white"
                : "bg-[#121b31] text-slate-400 hover:text-white"
                }`}
              onClick={() => setActiveTab("volume")}
            >
              Volumen de la Red
            </button>
          </div>

          {activeTab === "traffic" && (
            <>
              <p className="text-slate-400 text-sm">
                Monitoreo del tráfico de entrada y salida total.
              </p>
              <TrafficChart data={data} dataType="traffic" client:load />
            </>
          )}

          {activeTab === "volume" && (
            <>
              <p className="text-slate-400 text-sm">
                Monitoreo del volumen de datos de entrada y salida.
              </p>
              <TrafficChart data={data} dataType="volume" client:load />
            </>
          )}
        </>
      )}
    </section>
  );
}
