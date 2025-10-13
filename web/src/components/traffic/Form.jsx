import { useEffect, useState } from "react";
import dayjs from "dayjs";
import SelectField from "../ui/SelectField";
import DateField from "../ui/DateField";
import DatalistField from "../ui/DatalistField";
import RadioGroup from "../ui/RadioGroup";
import useFetch from "../../hooks/useFetch";
import { removeAccentsAndToUpper } from "../../utils/formater";
import {
  initDate,
  endDate,
  region,
  state,
  ip,
  municipality,
  county,
  odn,
  gpon,
} from "../../stores/traffic";
import { useStore } from "@nanostores/react";
import { isIpv4 } from "../../utils/validator";

const BASE_URL_TRAFFIC = `${import.meta.env.PUBLIC_URL || ""}/api/traffic`;
const BASE_URL_FATS = `${import.meta.env.PUBLIC_URL || ""}/api/fat`;
const TOKEN = sessionStorage.getItem("access_token")?.replace("Bearer ", "") || ""
endDate.set(dayjs().toISOString());
initDate.set(dayjs().subtract(1, "week").toISOString());

export default function Form() {
  const [urlFatState, setUrlFatState] = useState(undefined);
  const [urlOlt, setUrlOlt] = useState(undefined);
  const [regions, setRegions] = useState([])
  const [states, setStates] = useState([])
  const [selectionMethod, setSelectionMethod] = useState("");

  const $initDate = useStore(initDate);
  const $endDate = useStore(endDate);
  const $region = useStore(region);
  const $state = useStore(state);
  const $municipality = useStore(municipality);
  const $county = useStore(county);
  const $ip = useStore(ip);
  const $gpon = useStore(gpon);
  const $odn = useStore(odn);

  // Reset all form states when component mounts
  useEffect(() => {
    endDate.set(dayjs().toISOString());
    initDate.set(dayjs().subtract(1, "week").toISOString());
    region.set("");
    state.set("");
    municipality.set("");
    county.set("");
    ip.set("");
    gpon.set("");
    odn.set("");
    setSelectionMethod("");
  }, []);

  const headers = { headers: { Authorization: `Bearer ${TOKEN}` } }

  // Load all OLT from prometheus data
  const { data, status, loading, error } = useFetch(`${BASE_URL_TRAFFIC}/hierarchy?initDate=${$initDate}&finalDate=${$endDate}`, headers);

  // Load specific OLT info when URL is set
  const { data: infoOlt } = useFetch(urlOlt, headers);

  // Load specific OLT info when URL is set
  const { data: fatsByState } = useFetch(urlFatState, headers);

  // Get regions array and ensure OLT data is properly loaded
  useEffect(() => {
    if (data) setRegions(data.regions)
  }, [data]);

  const handleDateChange = ({ init, end }) => {
    if (init) initDate.set(init);
    if (end) endDate.set(end);
  };

  const handleChangeRegion = ({ target }) => {
    setStates(data.states[target.value])
    region.set(target.value);
    state.set("");
    municipality.set("");
    county.set("");
    ip.set("");
    gpon.set("");
    odn.set("");
    setSelectionMethod("");
  };

  const handleChangeState = ({ target }) => {
    const formatedState = removeAccentsAndToUpper(target.value);
    setUrlFatState(
      `${BASE_URL_FATS}/location/${formatedState}?page=1&limit=65535`
    );

    state.set(target.value);
    municipality.set("");
    county.set("");
    ip.set("");
    gpon.set("");
    odn.set("");
    setSelectionMethod("");
  };

  const handleChangeMethod = (method) => {
    setSelectionMethod(method);
    municipality.set("");
    county.set("");
    ip.set("");
    gpon.set("");
    odn.set("");
  };

  const handleChangeOlt = ({ target }) => {
    if (isIpv4(target.value)) {
      setUrlOlt(`${BASE_URL_TRAFFIC}/info/instance/${target.value}`);
    }
    ip.set(target.value);
    gpon.set("");
    municipality.set("");
    county.set("");
    odn.set("");
  };

  const handleChangeMunicipality = ({ target }) => {
    municipality.set(target.value);
    county.set("");
    odn.set("");
    ip.set("");
    gpon.set("");
  };

  const handleChangeCounty = ({ target }) => {
    county.set(target.value);
    ip.set("");
    gpon.set("");
    odn.set("");
  };

  const handleChangeOdn = ({ target }) => {
    odn.set(target.value);
  };

  const handleChangeGpon = ({ target }) => {
    gpon.set(target.value);
  };

  if (status === 401 || status === 403) {
    sessionStorage.removeItem("access_token");
    window.location.href = "/";
  }

  return (
    <form>
      <DateField
        id="initDate"
        label="Fecha inicial *"
        value={$initDate}
        max={$endDate ? $endDate.split("T")[0] : undefined}
        onChange={(init) => handleDateChange({ init })}
      />

      <DateField
        id="endDate"
        label="Fecha final *"
        endOfDay={true}
        value={$endDate}
        min={$initDate ? $initDate.split("T")[0] : undefined}
        onChange={(end) => handleDateChange({ end })}
      />

      <SelectField
        id="region"
        label="Región *"
        options={[
          {
            value: "",
            label: "Seleccionar región",
            disabled: true,
            hidden: true,
          },
          ...regions.map(region => ({ value: region, label: region })),
        ]}
        value={$region}
        onChange={handleChangeRegion}
      />

      {$initDate && $endDate && $region && (
        <SelectField
          id="state"
          label="Estado *"
          options={[
            {
              value: "",
              label: "Seleccionar estado",
              disabled: true,
              hidden: true,
            },
            ...states.map(state => ({ value: state, label: state })), ,
          ]}
          value={$state}
          onChange={handleChangeState}
        />
      )}

      {$state && (
        <RadioGroup
          id="selectionMethod"
          label="Método de busqueda"
          options={[
            { value: "olt", label: "OLT" },
            { value: "municipality", label: "Ubicación" },
          ]}
          value={selectionMethod}
          onChange={handleChangeMethod}
        />
      )}

      {$state && selectionMethod === "olt" && data && data.olts && data.olts[$state] && (
        <>
          <DatalistField
            id="olt"
            label="OLT"
            options={data.olts[$state].map(({ ip: value, sys_name: label }) => ({ value, label }))}
            value={$ip}
            onChange={handleChangeOlt}
            placeholder="Ingrese el OLT"
          />
          {$ip && infoOlt && infoOlt.length > 0 && (
            <SelectField
              id="gpon"
              label="GPON *"
              options={[
                {
                  value: "",
                  label: "Seleccionar Puerto GPON",
                  disabled: true,
                  hidden: true,
                },
                ...infoOlt.map(({ if_name, if_index }) => ({
                  value: if_index,
                  label: if_name,
                })),
              ]}
              value={$gpon}
              onChange={handleChangeGpon}
            />
          )}
        </>
      )}

      {$state && selectionMethod === "municipality" && fatsByState && (
        <>
          <SelectField
            id="municipality"
            label="Municipio *"
            options={[
              {
                value: "",
                label: "Seleccionar Municipio",
                disabled: true,
                hidden: true,
              },
              ...[...new Set(fatsByState.map((f) => f.municipality))].map(
                (m) => ({ value: m, label: m })
              ),
            ]}
            value={$municipality}
            onChange={handleChangeMunicipality}
          />

          {$municipality && (
            <SelectField
              id="county"
              label="Parroquia *"
              options={[
                {
                  value: "",
                  label: "Seleccionar Parroquia",
                  disabled: true,
                  hidden: true,
                },
                ...[
                  ...new Set(
                    fatsByState
                      .filter((f) => f.municipality === $municipality)
                      .map(({ county }) => county)
                  ),
                ].map((c) => ({ value: c, label: c })),
              ]}
              value={$county}
              onChange={handleChangeCounty}
            />
          )}

          {$county && (
            <SelectField
              id="odn"
              label="ODN *"
              options={[
                {
                  value: "",
                  label: "Seleccionar ODN",
                  disabled: true,
                  hidden: true,
                },
                ...[
                  ...new Set(
                    fatsByState
                      .filter(
                        (f) =>
                          f.municipality === $municipality &&
                          f.county === $county
                      )
                      .map(({ odn }) => odn)
                  ),
                ].map((o) => ({ value: o, label: o })),
              ]}
              value={$odn}
              onChange={handleChangeOdn}
            />
          )}
        </>
      )}
    </form>
  );
}
