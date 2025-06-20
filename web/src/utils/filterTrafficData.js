import { STATES_BY_REGION } from "../constants/regions";

/**
 * Agrupa y suma los datos de tráfico por día.
 * @param {Array} data - Datos originales.
 * @returns {Array} - Datos agrupados por día.
 */
function sumTrafficByDay(data) {
  const grouped = {};
  data.forEach((item) => {
    if (!grouped[item.day]) {
      grouped[item.day] = {
        day: item.day,
        mbps_in: 0,
        mbps_out: 0,
        mbytes_in_sec: 0,
        mbytes_out_sec: 0,
      };
    }
    grouped[item.day].mbps_in += item.mbps_in;
    grouped[item.day].mbps_out += item.mbps_out;
    grouped[item.day].mbytes_in_sec += item.mbytes_in_sec;
    grouped[item.day].mbytes_out_sec += item.mbytes_out_sec;
  });
  return Object.values(grouped);
}

/**
 * Agrupa y suma los datos de ont status por día.
 * @param {Array} data - Datos originales.
 * @returns {Array} - Datos agrupados por día.
 */
function sumStatusByDay(data) {
  const grouped = {};
  data.forEach((item) => {
    if (!grouped[item.day]) {
      grouped[item.day] = {
        day: item.day,
        ports_pon: 0,
        actives: 0,
        inactives: 0,
        unknowns: 0,
      };
    }
    grouped[item.day].ports_pon += item.ports_pon;
    grouped[item.day].actives += item.actives;
    grouped[item.day].inactives += item.inactives;
    grouped[item.day].unknowns += item.unknowns;
  });
  return Object.values(grouped);
}

/**
 * Filtra y agrupa los datos según los filtros seleccionados.
 * @param {Array} data - Datos originales.
 * @param {string} selectedLevel - "nacional" o "regional".
 * @param {string} selectedRegion - nombre de la región.
 * @param {string} selectedState - nombre del estado.
 * @returns {Array} - Datos filtrados y agrupados.
 */
export function filterTrafficData(
  data,
  selectedLevel,
  selectedRegion,
  selectedState
) {
  if (!selectedLevel || selectedLevel === "") {
    return sumTrafficByDay(data);
  }

  if (selectedLevel === "regional" && selectedRegion && !selectedState) {
    const states = STATES_BY_REGION[selectedRegion]?.map((s) => s.value) || [];
    const filtered = data.filter((item) => states.includes(item.description));
    return sumTrafficByDay(filtered);
  }

  if (selectedState) {
    const filtered = data.filter((item) => item.description === selectedState);
    return sumTrafficByDay(filtered);
  }

  return sumTrafficByDay(data);
}

/**
 * Filtra y agrupa los datos según los filtros seleccionados.
 * @param {Array} data - Datos originales.
 * @param {string} selectedLevel - "nacional" o "regional".
 * @param {string} selectedRegion - nombre de la región.
 * @param {string} selectedState - nombre del estado.
 * @returns {Array} - Datos filtrados y agrupados.
 */
export function filterStatusData(
  data,
  selectedLevel,
  selectedRegion,
  selectedState
) {
  if (!selectedLevel || selectedLevel === "") {
    return sumStatusByDay(data);
  }

  if (selectedLevel === "regional" && selectedRegion && !selectedState) {
    const states = STATES_BY_REGION[selectedRegion]?.map((s) => s.value) || [];
    const filtered = data.filter((item) => states.includes(item.description));
    return sumStatusByDay(filtered);
  }

  if (selectedState) {
    const filtered = data.filter((item) => item.description === selectedState);
    return sumStatusByDay(filtered);
  }

  return sumStatusByDay(data);
}
