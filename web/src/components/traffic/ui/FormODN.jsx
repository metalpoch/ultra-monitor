
import { useEffect, useState } from "react";
import { useStore } from "@nanostores/react";
import useFetch from "../../../hooks/useFetch";
import SelectField from "../../ui/SelectField";
import {
  endDate,
  olt,
  odn,
  fat,
  fats
} from "../../../stores/traffic";

const BASE_URL = `${import.meta.env.PUBLIC_URL}/api`

export default function FormODN() {
  const [fatURL, setFatURL] = useState("")
  const [odns, setOdns] = useState([]);
  const $endDate = useStore(endDate);
  const $olt = useStore(olt);
  const $odn = useStore(odn);

  const { data, status } = useFetch(fatURL, {
    headers: { Authorization: `Bearer ${sessionStorage.getItem("access_token").replace("Bearer ", "")}` },
  });

  useEffect(() => {
    if ($olt) {
      setOdns([])
      setFatURL(`${BASE_URL}/fat/ip/${$olt}`)
    }
  }, [$olt])

  useEffect(() => {
    if (data) {
      setOdns(
        [...new Set(data.map(({ odn }) => JSON.stringify({ odn })))]
          .map(strJson => JSON.parse(strJson))
          .map(({ odn }) => ({ label: odn, value: odn }))
      )
    }
  }, [data])

  const handleChangeOdn = ({ target }) => {
    odn.set(target.value)
    fats.set([...new Set(data.filter(({ odn }) => odn === target.value).map(({ fat }) => fat))].map(fat => ({ label: fat, value: fat }))
    )
    fat.set("")
  };

  if (status === 401) {
    sessionStorage.removeItem("access_token")
    window.location.href = "/";
  }

  if (odns.length > 0 && $endDate) return (
    <form className="w-full p-5 flex flex-wrap gap-5 content-center rounded-lg rounded-b-none bg-[#121b31] border-2 border-[hsl(217,33%,20%)] border-b-0">
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
          ...odns,
        ]}
        value={$odn}
        onChange={handleChangeOdn}
        disabled={odns.length === 0}
      />
    </form>
  );
}

