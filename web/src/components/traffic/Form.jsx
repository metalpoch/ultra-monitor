import { useEffect, useState } from "react";
import dayjs from "dayjs";
import SelectField from "../ui/SelectField";
import DateField from "../ui/DateField";
import DatalistField from "../ui/DatalistField";
import RadioGroup from "../ui/RadioGroup";
import useFetch from "../../hooks/useFetch";
import { removeAccentsAndToUpper } from "../../utils/formater";
import { initDate, endDate, region, state, ip, municipality, county, odn, gpon, oltsPrometheus } from "../../stores/traffic";
import { useStore } from "@nanostores/react";

const URL_TRAFFIC = `${import.meta.env.PUBLIC_URL}/api/traffic`
const URL_FATS = `${import.meta.env.PUBLIC_URL}/api/fat`

endDate.set(dayjs().toISOString());
initDate.set(dayjs().subtract(1, "week").toISOString());

export default function Form() {
  const [urlFatState, serUrlFatState] = useState(undefined)
  const [urlOlt, serUrlOlt] = useState(undefined)

  const [regions, setRegions] = useState([])
  const [states, setStates] = useState([])

  const [selectionMethod, setSelectionMethod] = useState("")

  const $initDate = useStore(initDate)
  const $endDate = useStore(endDate)
  const $region = useStore(region)
  const $state = useStore(state)
  const $municipality = useStore(municipality)
  const $county = useStore(county)
  const $ip = useStore(ip)
  const $odn = useStore(odn)
  const $gpon = useStore(gpon)

  const headers = {
    headers: {
      Authorization: `Bearer ${sessionStorage.getItem("access_token").replace("Bearer ", "")}`
    }
  }

  const { data: infoAllOlt, status } = useFetch(`${URL_TRAFFIC}/info`, headers);
  const { data: fatsByState } = useFetch(urlFatState, headers);
  const { data: infoOlt } = useFetch(urlOlt, headers);

  // Get regions array
  useEffect(() => {
    oltsPrometheus.set(infoAllOlt)
    if (infoAllOlt)
      setRegions([...new Set(infoAllOlt.map(({ region }) => region))].map(r => ({ value: r, label: r })))
  }, [infoAllOlt])

  const handleDateChange = ({ init, end }) => {
    if (init) initDate.set(init);
    if (end) endDate.set(end);
  };

  const handleChangeRegion = ({ target }) => {
    setStates(
      [...new Set(
        infoAllOlt
          .filter(item => item.region === target.value)
          .map(({ state }) => state)
      )].map(item => ({ value: item, label: item })))

    region.set(target.value)
    state.set("")
    municipality.set("")
    county.set("")
    ip.set("")
    gpon.set("")
    odn.set("")
  }

  const handleChangeState = ({ target }) => {
    const formatedState = removeAccentsAndToUpper(target.value)
    serUrlFatState(`${URL_FATS}/location/${formatedState}?page=1&limit=65535`)
    state.set(target.value)
    municipality.set("")
    county.set("")
    ip.set("")
    gpon.set("")
    odn.set("")
  }

  const handleChangeMethod = (method) => {
    setSelectionMethod(method)
    municipality.set("")
    county.set("")
    ip.set("")
    gpon.set("")
    odn.set("")
  }

  const handleChangeMunicipality = ({ target }) => {
    municipality.set(target.value)
    county.set("")
    ip.set("")
    gpon.set("")
    odn.set("")
  }

  const handleChangeCounty = ({ target }) => {
    county.set(target.value)
    ip.set("")
    gpon.set("")
    odn.set("")
  }

  const handleChangeOlt = ({ target }) => {
    serUrlOlt(`${URL_TRAFFIC}/info/instance/${target.value}`)
    ip.set(target.value)
    gpon.set("")
    odn.set("")
  }

  const handleChangeOdn = ({ target }) => {
    odn.set(target.value)
  }
  const handleChangeGpon = ({ target }) => {
    gpon.set(target.value)
  }

  if (status === 401) {
    sessionStorage.removeItem("access_token")
    window.location.href = "/";
  }

  return (
    <form>
      {/* Fecha inicial */}
      <DateField
        id="initDate"
        label="Fecha inicial *"
        value={$initDate}
        onChange={(init) => handleDateChange({ init })}
      />

      {/* Fecha final */}
      <DateField
        id="endDate"
        label="Fecha final *"
        endOfDay={true}
        value={$endDate}
        onChange={(end) => handleDateChange({ end })}
      />

      {/* Región */}
      <SelectField
        id="region"
        label="Región *"
        options={[
          {
            value: "",
            label: "Seleccionar región",
            disabled: true,
            hidden: true
          },
          ...regions,
        ]}
        value={$region}
        onChange={handleChangeRegion}
      />

      {/* Estado: se habilita solo si los tres anteriores están completos */}
      {$initDate && $endDate && $region && (
        <SelectField
          id="state"
          label="Estado *"
          options={[
            { value: "", label: "Seleccionar estado", disabled: true, hidden: true },
            ...states,
          ]}
          value={$state}
          onChange={handleChangeState}
        />
      )}

      {/* Método de selección: Municipio/Parroquia o OLT directo */}
      {$state && (
        <RadioGroup
          id="selectionMethod"
          label="Método de selección"
          options={[
            { value: "olt", label: "Seleccionar OLT directamente" },
            { value: "municipality", label: "Seleccionar por Municipio" },
          ]}
          value={selectionMethod}
          onChange={handleChangeMethod}
        />
      )}

      {/* Opción 1: OLT directo */}
      {$state && selectionMethod === "olt" && (
        <>
          <DatalistField
            id="olt"
            label="OLT"
            options={
              infoAllOlt
                .filter(item => item.region === $region && item.state === $state)
                .map(({ ip, sysName }) => ({ value: ip, label: sysName }))}
            value={$ip}
            onChange={handleChangeOlt}
            placeholder="Ingrese el OLT"
          />
          {$ip && infoOlt && infoOlt.length > 0 && (
            <SelectField
              id="gpon"
              label="GPON *"
              options={[
                { value: "", label: "Seleccionar Puerto GPON", disabled: true, hidden: true },
                ...infoOlt.map(({ if_name, if_index }) => ({ value: if_index, label: if_name }))
              ]}
              value={$gpon}
              onChange={handleChangeGpon}
            />
          )}
        </>
      )}

      {/* Opción 2: Municipio → Parroquia → OLT */}
      {$state && selectionMethod === "municipality" && fatsByState &&
        <>
          <SelectField
            id="municipality"
            label="Municipio *"
            options={[
              { value: "", label: "Seleccionar Municipio", disabled: true, hidden: true },
              ...[...new Set(fatsByState.map(f => f.municipality))].map(m => ({ value: m, label: m }))
            ]}
            value={$municipality}
            onChange={handleChangeMunicipality}
          />

          {$municipality && (
            <SelectField
              id="county"
              label="Parroquia *"
              options={[
                { value: "", label: "Seleccionar Parroquia", disabled: true, hidden: true },
                ...[...new Set(fatsByState.filter((f) => f.municipality === $municipality).map(({ county }) => county))].map(c => ({ value: c, label: c }))
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
                { value: "", label: "Seleccionar ODN", disabled: true, hidden: true },
                ...[...new Set(fatsByState.filter((f) => f.municipality === $municipality && f.county === $county).map(({ odn }) => odn))]
                  .map(o => ({ value: o, label: o }))
              ]}
              value={$odn}
              onChange={handleChangeOdn}
            />
          )}
        </>
      }
    </form>
  )
}
