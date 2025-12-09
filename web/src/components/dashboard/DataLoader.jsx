import { useEffect } from "react";
import { useStore } from "@nanostores/react";
import dayjs from "dayjs";
import {
  selectedLevel,
  selectedRegion,
  selectedState,
  dailyTrafficData,
  userData
} from "../../stores/dashboard";
import { removeAccentsAndToUpper } from "../../utils/formater";
import useFetch from "../../hooks/useFetch";

const BASE_URL_TRAFFIC = `${import.meta.env.PUBLIC_URL || ""}/api/traffic`;
const BASE_URL_FAT = `${import.meta.env.PUBLIC_URL || ""}/api/fat/trend/status`;
const BASE_URL_STATUS = `${
  import.meta.env.PUBLIC_URL || ""
}/api/prometheus/status`;

// Obtener datos del último año (365 días)
const getLastYearDate = () => {
  const minDate = dayjs("2025-07-01T00:00:00-04:00");
  const today = dayjs()
    .set("hour", 0)
    .set("minute", 0)
    .set("second", 0)
    .set("millisecond", 0);

  const lastYear = today.subtract(1, "year");
  return lastYear < minDate ? minDate : lastYear;
};

export default function DataLoader() {
  const $selectedLevel = useStore(selectedLevel);
  const $selectedRegion = useStore(selectedRegion);
  const $selectedState = useStore(selectedState);

  const token = sessionStorage.getItem("access_token")?.replace("Bearer ", "");

  // Configurar URL para datos diarios del último año
  useEffect(() => {
    if (!token) return;

    // Solo establecer loading: true si realmente necesitamos cargar nuevos datos
    // (cuando cambia la región o estado seleccionado)
    const prevTrafficData = dailyTrafficData.get();
    const prevUserData = userData.get();

    // Determinar si necesitamos cargar nuevos datos
    const needsNewData =
      ($selectedState && prevTrafficData.lastState !== $selectedState) ||
      ($selectedRegion && prevTrafficData.lastRegion !== $selectedRegion) ||
      (!$selectedState && !$selectedRegion && (prevTrafficData.lastState || prevTrafficData.lastRegion));

    if (needsNewData) {
      dailyTrafficData.set({
        ...prevTrafficData,
        loading: true
      });

      userData.set({
        ...prevUserData,
        loading: true
      });
    }

    const initDate = getLastYearDate();
    const finalDate = dayjs()
      .set("hour", 0)
      .set("minute", 0)
      .set("second", 0)
      .set("millisecond", 0);

    const params = new URLSearchParams({
      initDate: initDate.toISOString(),
      finalDate: finalDate.toISOString()
    });

    const baseUrl = $selectedState
      ? `${BASE_URL_TRAFFIC}/state/${$selectedState}`
      : $selectedRegion
      ? `${BASE_URL_TRAFFIC}/region/${$selectedRegion}`
      : `${BASE_URL_TRAFFIC}/total`;

    const trafficUrl = `${baseUrl}?${params.toString()}`;

    // Configurar URLs para usuarios
    let fatUrl, statusUrl;
    if ($selectedState) {
      fatUrl = `${BASE_URL_FAT}/state/${removeAccentsAndToUpper($selectedState)}`;
      statusUrl = `${BASE_URL_STATUS}/state/${$selectedState}`;
    } else if ($selectedRegion) {
      fatUrl = `${BASE_URL_FAT}/${$selectedRegion}`;
      statusUrl = `${BASE_URL_STATUS}/region/${$selectedRegion}`;
    } else {
      fatUrl = BASE_URL_FAT;
      statusUrl = BASE_URL_STATUS;
    }
  }, [$selectedLevel, $selectedRegion, $selectedState, token]);

  // Obtener datos de usuarios
  const { data: fatData } = useFetch(
    $selectedState
      ? `${BASE_URL_FAT}/state/${removeAccentsAndToUpper($selectedState)}`
      : $selectedRegion
      ? `${BASE_URL_FAT}/${$selectedRegion}`
      : BASE_URL_FAT,
    {
      headers: { Authorization: `Bearer ${token}` },
      skip: !token
    }
  );

  const { data: gponData } = useFetch(
    $selectedState
      ? `${BASE_URL_STATUS}/state/${$selectedState}`
      : $selectedRegion
      ? `${BASE_URL_STATUS}/region/${$selectedRegion}`
      : BASE_URL_STATUS,
    {
      headers: { Authorization: `Bearer ${token}` },
      skip: !token
    }
  );

  // Obtener datos de tráfico diarios del último año
  const getTrafficUrl = () => {
    const initDate = getLastYearDate();
    const finalDate = dayjs()
      .set("hour", 0)
      .set("minute", 0)
      .set("second", 0)
      .set("millisecond", 0);

    const params = new URLSearchParams({
      initDate: initDate.toISOString(),
      finalDate: finalDate.toISOString()
    });

    const baseUrl = $selectedState
      ? `${BASE_URL_TRAFFIC}/state/${$selectedState}`
      : $selectedRegion
      ? `${BASE_URL_TRAFFIC}/region/${$selectedRegion}`
      : `${BASE_URL_TRAFFIC}/total`;

    return `${baseUrl}?${params.toString()}`;
  };

  const { data: trafficData } = useFetch(getTrafficUrl(), {
    headers: { Authorization: `Bearer ${token}` },
    skip: !token
  });

  // Actualizar store de usuarios cuando lleguen los datos
  useEffect(() => {
    if (fatData && gponData) {
      const lastFatData = fatData[fatData.length - 1];
      const currOntAct = lastFatData.actives + lastFatData.provisioned_offline;

      userData.set({
        activeUsers: currOntAct,
        loading: false,
        error: null,
        lastState: $selectedState,
        lastRegion: $selectedRegion
      });
    }
  }, [fatData, gponData, $selectedState, $selectedRegion]);

  // Actualizar store de tráfico diario cuando lleguen los datos
  useEffect(() => {
    if (trafficData) {
      dailyTrafficData.set({
        data: trafficData || [],
        loading: false,
        error: null,
        lastState: $selectedState,
        lastRegion: $selectedRegion
      });
    }
  }, [trafficData, $selectedState, $selectedRegion]);

  // Este componente no renderiza nada visible
  return null;
}

