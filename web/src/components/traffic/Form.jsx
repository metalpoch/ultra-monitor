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
  oltName,
  switchValue,
  // municipality,
  // odn,
  gpon,
} from "../../stores/traffic";
import { useStore } from "@nanostores/react";
import { isIpv4 } from "../../utils/validator";

const BASE_URL_TRAFFIC = `${import.meta.env.PUBLIC_URL || ""}/api/traffic`;
const BASE_URL_FATS = `${import.meta.env.PUBLIC_URL || ""}/api/fat`;
const TOKEN = sessionStorage.getItem("access_token")?.replace("Bearer ", "") || ""
endDate.set(dayjs().toJSON());
initDate.set(dayjs().subtract(1, "week").startOf("day").toJSON());

export default function Form() {

  // const [urlOdnStatsByMunicipality, setUrlOdnStatByMunicipality] = useState(undefined);
  // const [urlMunicipalityByState, setUrlunicipalityByState] = useState(undefined);
  const [urlOlt, setUrlOlt] = useState(undefined);
  const [regions, setRegions] = useState([])
  const [states, setStates] = useState([])
  const [selectionMethod, setSelectionMethod] = useState("");

  const $initDate = useStore(initDate);
  const $endDate = useStore(endDate);
  const $region = useStore(region);
  const $state = useStore(state);
  // const $municipality = useStore(municipality);
  const $ip = useStore(ip);
  const $gpon = useStore(gpon);
  // const $odn = useStore(odn);

  // Reset all form states when component mounts
  useEffect(() => {
    endDate.set(dayjs().toISOString());
    initDate.set(dayjs().subtract(1, "week").startOf("day").toJSON());
    region.set("");
    state.set("");
    // municipality.set("");
    ip.set("");
    oltName.set("");
    switchValue.set("")
    gpon.set("");
    // odn.set("");
    setSelectionMethod("");
  }, []);

  const headers = { headers: { Authorization: `Bearer ${TOKEN}` } }

  const { data, status, loading, error } = useFetch(`${BASE_URL_TRAFFIC}/hierarchy?initDate=${$initDate}&finalDate=${$endDate}`, headers);
  const { data: infoOlt } = useFetch(urlOlt, headers);
  // const { data: odnMunicipalities } = useFetch(urlMunicipalityByState, headers);
  // const { data: odnStatsByMunicipality } = useFetch(urlOdnStatsByMunicipality, headers);

  // Get regions array and ensure OLT data is properly loaded
  useEffect(() => {
    if (data) setRegions(data.regions)
  }, [data]);

  const handleDateChange = ({ init, end }) => {
    const hour = dayjs().hour()
    const minute = dayjs().minute()
    if (init) initDate.set(dayjs(init).toJSON());
    if (end) endDate.set(dayjs(end).hour(hour).minute(minute).toJSON());
  };

  const handleChangeRegion = ({ target }) => {
    setStates(data.states[target.value])
    region.set(target.value);
    state.set("");
    // municipality.set("");
    ip.set("");
    oltName.set("");
    gpon.set("");
    // odn.set("");
    setSelectionMethod("");
  };

  const handleChangeState = ({ target }) => {
    // const formatedState = removeAccentsAndToUpper(target.value);
    // setUrlunicipalityByState(`${BASE_URL_FATS}/municipalities/${formatedState}`);

    state.set(target.value);
    // municipality.set("");
    ip.set("");
    oltName.set("");
    switchValue.set("");
    gpon.set("");
    // odn.set("");
    setSelectionMethod("");
  };

  const handleChangeMethod = (method) => {
    setSelectionMethod(method);
    // municipality.set("");
    ip.set("");
    oltName.set("");
    switchValue.set("");
    gpon.set("");
    // odn.set("");
  };

  const handleChangeOlt = ({ target }) => {
    if (isIpv4(target.value)) {
      setUrlOlt(`${BASE_URL_TRAFFIC}/info/instance/${target.value}`);
    }

    const selectedObject = data.olts[$state].find(({ ip }) => ip === target.value);

    ip.set(target.value);
    oltName.set(selectedObject.sys_name);
    gpon.set("");
    // municipality.set("");
    // odn.set("");
  };

  // const handleChangeMunicipality = ({ target }) => {
  //   const formatedState = removeAccentsAndToUpper($state);
  //   setUrlOdnStatByMunicipality(`${BASE_URL_FATS}/stats/${formatedState}/${target.value}?finalDate=${$endDate}`);

  //   municipality.set(target.value);
  //   odn.set("");
  //   ip.set("");
  //   gpon.set("");
  // };

  // const handleChangeOdn = ({ target }) => {
  //   odn.set(target.value);
  // };

  const handleChangeGpon = ({ target }) => {
    gpon.set(target.value);
  };

  if (status === 401 || status === 403) {
    sessionStorage.removeItem("access_token");
    window.location.href = "/";
  }

  if (infoOlt) {
    switchValue.set(infoOlt[0]?.switch);
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

      {/* {$state && selectionMethod === "municipality" && odnMunicipalities && (
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
              ...odnMunicipalities.toSorted().map(m => ({ value: m, label: m }))
            ]}
            value={$municipality}
            onChange={handleChangeMunicipality}
          />

          {$municipality && odnStatsByMunicipality && (
            <SelectField
              id="odn"
              label="ODN"
              options={[
                {
                  value: "",
                  label: "Seleccionar ODN",
                  disabled: true,
                  hidden: true,
                },
                ...[
                  ...new Set(
                    odnStatsByMunicipality
                      .map(({ odn }) => odn)
                  ),
                ].toSorted().map((c) => ({ value: c, label: c })),
              ]}
              value={$odn}
              onChange={handleChangeOdn}
            />
          )}

        </>
      )} */}
    </form>
  );
}
