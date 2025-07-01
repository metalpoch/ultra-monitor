import { useEffect, useState, useMemo } from "react";
import SelectField from "../../ui/SelectField";
import DatalistField from "../../ui/DatalistField";
import { trafficData, ontData } from "../../../stores/traffic";
import { REGIONS, STATES_BY_REGION } from "../../../constants/regions";
import useFetch from "../../../hooks/useFetch";
import { getDateRange } from "../../../utils/convert";

const BASE_URL = import.meta.env.PUBLIC_API_URL;

export default function Form() {
  const [trafficURL, setTrafficURL] = useState("");
  const [ontStatusURL, setOntStatusURL] = useState("");
  const [selectedDate, setSelectedDate] = useState("1d");
  const [selectedRegion, setSelectedRegion] = useState("");
  const [selectedState, setSelectedState] = useState("");
  const [selectedMunicipality, setSelectedMunicipality] = useState("");
  const [selectedCounty, setSelectedCounty] = useState("");
  const [selectedOdn, setSelectedOdn] = useState("");
  const [selectedOlt, setSelectedOlt] = useState("");
  const [selectedFat, setSelectedFat] = useState("");

  const dataTraffic = useFetch(trafficURL);
  const dataOnt = useFetch(ontStatusURL);

  useEffect(() => trafficData.set(dataTraffic), [dataTraffic]);
  useEffect(() => ontData.set(dataOnt), [dataOnt]);

  const url = selectedState
    ? `${BASE_URL}/olt/location/${encodeURIComponent(selectedState)}`
    : null;

  const { data: rawData, loading, error } = useFetch(url);
  const data = Array.isArray(rawData) ? rawData : [];

  const municipalities = useMemo(
    () =>
      Array.from(
        new Set(data.map((item) => item.municipality).filter(Boolean))
      ),
    [data]
  );
  const counties = useMemo(
    () =>
      Array.from(
        new Set(
          data
            .filter((item) =>
              selectedMunicipality
                ? item.municipality === selectedMunicipality
                : true
            )
            .map((item) => item.county)
            .filter(Boolean)
        )
      ),
    [data, selectedMunicipality]
  );

  const filteredOlts = useMemo(
    () =>
      data.filter(
        (item) =>
          (!selectedMunicipality ||
            item.municipality === selectedMunicipality) &&
          (!selectedCounty || item.county === selectedCounty) &&
          (!selectedOdn || item.odn === selectedOdn)
      ),
    [data, selectedMunicipality, selectedCounty, selectedOdn]
  );

  const oltOptions = useMemo(() => {
    const seen = new Set();
    return filteredOlts
      .filter(({ ip }) => {
        if (seen.has(ip)) return false;
        seen.add(ip);
        return true;
      })
      .map(({ ip, sys_name, sys_location }) => ({
        value: sys_name,
        label: `${ip} / ${sys_location}`,
      }))
      .filter((item) => item.value);
  }, [filteredOlts]);

  const fatOptions = useMemo(() => {
    const seen = new Set();
    return filteredOlts
      .filter(({ fat }) => {
        if (seen.has(fat)) return false;
        seen.add(fat);
        return true;
      })
      .map(({ ip, fat, sys_location }) => ({
        value: fat,
        label: `${ip} / ${sys_location}`,
      }));
  }, [filteredOlts]);

  const odnOptions = useMemo(() => {
    const seen = new Set();
    return filteredOlts
      .filter((item) => {
        if (seen.has(item.odn)) return false;
        seen.add(item.odn);
        return true;
      })
      .map(({ ip, odn, sys_location }) => ({
        value: odn,
        label: `${ip} / ${sys_location}`,
      }))
      .filter((item) => item.value);
  }, [filteredOlts]);

  useEffect(() => {
    setSelectedState("");
    setSelectedMunicipality("");
    setSelectedCounty("");
    setSelectedOdn("");
    setSelectedOlt("");
    setSelectedFat("");
  }, [selectedRegion]);

  useEffect(() => {
    setSelectedMunicipality("");
    setSelectedCounty("");
    setSelectedOdn("");
    setSelectedOlt("");
    setSelectedFat("");
  }, [selectedState]);

  useEffect(() => {
    setSelectedCounty("");
    setSelectedOdn("");
    setSelectedOlt("");
    setSelectedFat("");
  }, [selectedMunicipality]);

  useEffect(() => {
    setSelectedOdn("");
    setSelectedOlt("");
    setSelectedFat("");
  }, [selectedCounty]);

  useEffect(() => {
    if (selectedOlt) {
      setSelectedFat("");
      setSelectedOdn("");

      const sysname = encodeURIComponent(selectedOlt);
      const { initDate, endDate } = getDateRange(selectedDate);
      const rangeDate = `initDate=${initDate}&endDate=${endDate}`;
      setTrafficURL(`${BASE_URL}/pon/traffic/${sysname}?${rangeDate}`);
      setOntStatusURL(`${BASE_URL}/ont/status/sysname/${sysname}?${rangeDate}`);
    }
  }, [selectedOlt, selectedDate]);

  useEffect(() => {
    if (selectedFat) {
      setSelectedOlt("");
      setSelectedOdn("");

      const param = encodeURIComponent(selectedFat);
      const { initDate, endDate } = getDateRange(selectedDate);
      const rangeDate = `initDate=${initDate}&endDate=${endDate}`;
      setTrafficURL(`${BASE_URL}/fat/traffic/${param}?${rangeDate}`);
    }
  }, [selectedFat, selectedDate]);

  useEffect(() => {
    if (selectedOdn) {
      setSelectedOlt("");
      setSelectedFat("");

      const param = encodeURIComponent(selectedOdn);
      const { initDate, endDate } = getDateRange(selectedDate);
      const rangeDate = `initDate=${initDate}&endDate=${endDate}`;
      setTrafficURL(`${BASE_URL}/odn/traffic/${param}?${rangeDate}`);
    }
  }, [selectedOdn, selectedDate]);

  return (
    <form className="w-full p-5 flex flex-wrap gap-5 content-center rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
      <section className="w-full flex flex-wrap gap-2">
        <SelectField
          id="date"
          label="Visualizar"
          options={[
            { value: "1d", label: "24 horas" },
            { value: "7d", label: "1 semana" },
            { value: "1m", label: "1 mes" },
          ]}
          value={selectedDate}
          onChange={({ target }) => setSelectedDate(target.value)}
        />
        <SelectField
          id="region"
          label="Región"
          options={[
            {
              value: "",
              label: "Seleccionar región",
              disabled: true,
              hidden: true,
            },
            ...REGIONS.sort((a, b) => a.label.localeCompare(b.label)),
          ]}
          value={selectedRegion}
          onChange={({ target }) => setSelectedRegion(target.value)}
        />
        <SelectField
          id="state"
          label="Estado*"
          options={[
            {
              value: "",
              label: "Seleccionar estado",
              disabled: true,
              hidden: true,
            },
            ...(STATES_BY_REGION[selectedRegion] ||
              Object.values(STATES_BY_REGION)
                .flat()
                .toSorted((a, b) => a.label.localeCompare(b.label))),
          ]}
          value={selectedState}
          onChange={({ target }) => setSelectedState(target.value)}
        />
        {selectedState && (
          <SelectField
            id="municipality"
            label="Municipio"
            options={[
              {
                value: "",
                label: "Seleccionar municipio",
                disabled: true,
                hidden: true,
              },
              ...municipalities.map((m) => ({ value: m, label: m })),
            ]}
            value={selectedMunicipality}
            onChange={({ target }) => setSelectedMunicipality(target.value)}
            disabled={!selectedState || municipalities.length === 0}
          />
        )}
        {selectedMunicipality && (
          <SelectField
            id="county"
            label="Parroquia"
            options={[
              {
                value: "",
                label: "Seleccionar parroquia",
                disabled: true,
                hidden: true,
              },
              ...counties.map((c) => ({ value: c, label: c })),
            ]}
            value={selectedCounty}
            onChange={({ target }) => setSelectedCounty(target.value)}
            disabled={!selectedMunicipality || counties.length === 0}
          />
        )}
      </section>

      <section className="flex flex-wrap gap-2">
        {oltOptions.length > 0 && (
          <DatalistField
            id="olt"
            label="OLT"
            options={oltOptions.map((o) => ({
              value: o.value,
              label: o.label,
            }))}
            value={selectedOlt || ""}
            onChange={({ target }) => setSelectedOlt(target.value)}
            placeholder="Ingrese el OLT"
            disabled={oltOptions.length === 0}
          />
        )}

        {fatOptions.length > 0 && (
          <DatalistField
            id="fat"
            label="FAT"
            options={fatOptions.map((f) => ({
              value: f.value,
              label: f.label,
            }))}
            value={selectedFat || ""}
            onChange={({ target }) => setSelectedFat(target.value)}
            placeholder="Ingrese el FAT"
            disabled={fatOptions.length === 0}
          />
        )}

        {odnOptions.length > 0 && (
          <DatalistField
            id="odn"
            label="ODN"
            options={odnOptions.map((f) => ({
              value: f.value,
              label: f.label,
            }))}
            value={selectedOdn || ""}
            onChange={({ target }) => setSelectedOdn(target.value)}
            placeholder="Ingrese el ODN"
            disabled={odnOptions.length === 0}
          />
        )}
      </section>
      {loading && (
        <div className="w-full text-center text-cyan-400">
          Cargando datos...
        </div>
      )}
      {error && (
        <div className="w-full text-center text-red-400">
          Error al cargar datos
        </div>
      )}
    </form>
  );
}
