import { useStore } from "@nanostores/react";
import { trafficForm, trafficFormDate } from "../../../stores/traffic";
import { getDateRange } from "../../../utils/convert";
import useFetch from "../../../hooks/useFetch";
import { useEffect, useState } from "react";

const BASE_URL = import.meta.env.PUBLIC_API_URL;

export default function TrafficChart() {
  const [url, setUrl] = useState("");
  const $trafficForm = useStore(trafficForm);
  const $trafficFormDate = useStore(trafficFormDate);

  const { data, loading, error } = useFetch(url);

  useEffect(() => {
    if (!$trafficFormDate || !$trafficForm.type) {
      return;
    }
    const { initDate, endDate } = getDateRange($trafficFormDate);
    const { type, value } = $trafficForm;

    if (type === "olt")
      setUrl(
        `${BASE_URL}/pon/traffic/${encodeURIComponent(
          value
        )}?initDate=${initDate}&endDate=${endDate}`
      );
  }, [$trafficForm, $trafficFormDate]);

  return (
    <>
      {error && <p>{JSON.stringify(error)}</p>}
      {loading && <p>loading</p>}
      {data && <p>{JSON.stringify(data)}</p>}
    </>
  );
}
