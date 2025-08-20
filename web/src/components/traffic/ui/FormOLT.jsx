import { useEffect, useState } from "react";
import { useStore } from "@nanostores/react";
import useFetch from "../../../hooks/useFetch";
import SelectField from "../../ui/SelectField";
import DateField from "../../ui/DateField";
import DatalistField from "../../ui/DatalistField";
import {
  initDate,
  endDate,
  olt,
  odn,
  fat
} from "../../../stores/traffic";

const BASE_URL = `${import.meta.env.PUBLIC_URL}/api`

export default function Form() {
  const [state, setState] = useState("");
  const [olts, setOlts] = useState([]);
  const [states, setStates] = useState([]);
  const $initDate = useStore(initDate);
  const $endDate = useStore(endDate);
  const $olt = useStore(olt);

  const { data, status } = useFetch(`${BASE_URL}/traffic/info`, {
    headers: { Authorization: `Bearer ${sessionStorage.getItem("access_token").replace("Bearer ", "")}` },
  });

  useEffect(() => {
    if (data)
      setStates(
        [...new Set(data.map(({ state }) => state))]
          .sort((a, b) => a.localeCompare(b))
          .map(reg => ({ label: reg, value: reg }))
      )
  }, [data])

  useEffect(() => {
    if (state) {
      setOlts(
        data.filter(item => item.state === state).map(({ ip, sysName }) => ({ ip, sysName }))
          .sort((a, b) => a.sysName.localeCompare(b.sysName))
          .map(({ ip, sysName }) => ({ label: sysName, value: ip }))
      )
    }
  }, [state]);

  const handleChangeState = ({ target }) => {
    setState(target.value)
    olt.set("")
    odn.set("")
    fat.set("")
  };

  const handleChangeOlt = ({ target }) => {
    olt.set(target.value)
    odn.set("")
    fat.set("")
  };

  const handleDateChange = ({ init, end }) => {
    if (init) initDate.set(init);
    if (end) endDate.set(end);
  };

  if (status === 401) {
    sessionStorage.removeItem("access_token")
    window.location.href = "/";
  }

  return (
    <form className="w-full p-5 flex flex-wrap gap-5 content-center rounded-lg rounded-b-none bg-[#121b31] border-2 border-[hsl(217,33%,20%)] border-b-0">
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
        value={state}
        onChange={handleChangeState}
        disabled={states.length === 0}
      />

      {state && <DatalistField
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
      />}

      {$olt &&
        <DateField
          id="initDate"
          label="Fecha inicial"
          value={$initDate}
          onChange={(init) => handleDateChange({ init })}
        />
      }
      {$initDate &&
        < DateField
          id="endDate"
          label="Fecha final"
          endOfDay={true}
          value={$endDate}
          onChange={(end) => handleDateChange({ end })}
        />
      }

    </form>
  );
}
