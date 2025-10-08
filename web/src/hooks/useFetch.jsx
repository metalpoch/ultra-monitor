import { useState, useEffect, useCallback } from "react";

export default function useFetch(url, options = {}) {
  const [data, setData] = useState(null);
  const [status, setStatus] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchData = useCallback(() => {
    if (!url) return;
    const abortController = new AbortController();
    setLoading(true);
    setError(null);
    setData(null);

    fetch(url, { ...options, signal: abortController.signal })
      .then((res) => {
        setStatus(res.status);
        setLoading(false);
        return res.json();
      })
      .then((res) => setData(res))
      .catch((err) => {
        if (err.name === "AbortError") {
          console.warn("Petición Abortada");
        } else {
          setError(err);
          setLoading(false)
        }
      })

    return () => abortController.abort();
  }, [url, JSON.stringify(options)]);

  useEffect(() => {
    const abort = fetchData();
    return abort;
  }, [fetchData]);

  return { data, status, loading, error, refetch: fetchData };
}
