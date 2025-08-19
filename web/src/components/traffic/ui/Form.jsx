import { useEffect, useState } from "react";
import { useStore } from "@nanostores/react";
import useFetch from "../../../hooks/useFetch";
import SelectField from "../../ui/SelectField";
import DatalistField from "../../ui/DatalistField";
import {
  region,
  state,
  olt,
  odn,
  fat,
} from "../../../stores/traffic";

const BASE_URL = `${import.meta.env.PUBLIC_URL}/api`

export default function Form() {
  const [fatURL, setFatURL] = useState("")
  const [fats, setFats] = useState([]);
  const [odns, setOdns] = useState([]);
  const [olts, setOlts] = useState([]);
  const [states, setStates] = useState([]);
  const [regions, setRegions] = useState([]);
  const $region = useStore(region);
  const $state = useStore(state);
  const $olt = useStore(olt);
  const $odn = useStore(odn);
  const $fat = useStore(fat);

  const { data: deviceInfo, status } = useFetch(`${BASE_URL}/traffic/info`, {
    headers: { Authorization: `Bearer ${sessionStorage.getItem("access_token").replace("Bearer ", "")}` },
  });

  const { data: dataFats } = useFetch(fatURL, {
    headers: { Authorization: `Bearer ${sessionStorage.getItem("access_token").replace("Bearer ", "")}` },
  });

  useEffect(() => {
    if (deviceInfo)
      setRegions(
        [...new Set(deviceInfo.map(({ region }) => region))]
          .sort((a, b) => a.localeCompare(b))
          .map(reg => ({ label: reg, value: reg }))
      )
  }, [deviceInfo])

  useEffect(() => {
    if ($region)
      setStates(
        [...new Set(deviceInfo.filter(item => item.region === $region).map(({ state }) => state))]
          .sort((a, b) => a.localeCompare(b))
          .map(reg => ({ label: reg, value: reg }))
      )
  }, [$region])

  useEffect(() => {
    if ($state) {
      setOlts(
        deviceInfo.filter(item => item.state === $state).map(({ ip, sysName }) => ({ ip, sysName }))
          .sort((a, b) => a.sysName.localeCompare(b.sysName))
          .map(({ ip, sysName }) => ({ label: sysName, value: ip }))
      )
    }
  }, [$state]);

  useEffect(() => {
    if (dataFats && $olt) {
      setOdns(
        [...new Set(dataFats.map(({ odn, municipality, county }) => JSON.stringify({ odn, municipality, county })))]
          .map(strJson => JSON.parse(strJson))
          .map(({ odn, municipality, county }) => ({ label: `${municipality} - ${county}`, value: odn }))
      )
    }
    if (dataFats && $odn) {
      setFats(
        [...new Set(dataFats.map(({ id, fat, shell, card, port }) => JSON.stringify({ id, fat, shell, card, port })))]
          .map(strJson => JSON.parse(strJson))
          .map(({ id, fat, shell, card, port }) => ({ label: `${id} - GPON ${shell}/${card}/${port}`, value: fat }))
      )
    }
  }, [dataFats, $olt, $odn]);

  const handleChangeRegion = ({ target }) => {
    region.set(target.value)
    state.set("")
    olt.set("")
    odn.set("")
    fat.set("")
    setFatURL("")
  };

  const handleChangeState = ({ target }) => {
    state.set(target.value);
    olt.set("")
    odn.set("")
    fat.set("")
    setFatURL("")
  };

  const handleChangeOlt = ({ target }) => {
    olt.set(target.value)
    odn.set("")
    fat.set("")
    setFatURL(`${BASE_URL}/fat/ip/${target.value}`)
  };

  const handleChangeOdn = ({ target }) => {
    odn.set(target.value)
    fat.set("")
  };
  const handleChangeFat = ({ target }) => {
    fat.set(target.value)
  };


  if (status === 401) {
    sessionStorage.removeItem("access_token")
    window.location.href = "/";
  }
  return (
    <form className="w-full p-5 flex flex-wrap gap-5 content-center rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
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
          ...regions,
        ]}
        value={$region}
        onChange={handleChangeRegion}
      />

      <SelectField
        id="states"
        label="Estado"
        options={[
          {
            value: "",
            label: "Seleccionar estado",
            disabled: true,
            hidden: true,
          },
          ...states,
        ]}
        value={$state}
        onChange={handleChangeState}
        disabled={states.length === 0}
      />

      <DatalistField
        id="olts"
        label="OLTs"
        options={olts.map((o) => ({
          value: o.value,
          label: o.label,
        }))}
        value={$olt}
        onChange={handleChangeOlt}
        placeholder="Ingrese el OLT"
        disabled={olts.length === 0}
      />

      <DatalistField
        id="odn"
        label="ODN"
        options={[
          {
            value: "",
            label: "Seleccionar ODN",
            disabled: true,
            hidden: true,
          },
          ...odns,
        ]}
        value={$odn}
        onChange={handleChangeOdn}
        placeholder="Ingrese el ODN"
        disabled={odns.length === 0}
      />

      <DatalistField
        id="fat"
        label="FAT"
        options={[
          {
            value: "",
            label: "Seleccionar Fat",
            disabled: true,
            hidden: true,
          },
          ...fats,
        ]}
        value={$fat}
        onChange={handleChangeFat}
        placeholder="Ingrese el FAT"
        disabled={fats.length === 0}
      />
    </form>
  );
}
