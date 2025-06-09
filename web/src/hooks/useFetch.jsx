import { useState, useEffect } from "react"

export default function useFetch(url, options) {
    const [data, setData] = useState(null);
    const [status, setStatus] = useState(null);
    const [loading, setLoading] = useState(null);
    const [error, setError] = useState(null);


    if (!url || !options) {
        return { data, status, loading, error }
    }

    useEffect(() => {
        const abortController = new AbortController();
        setLoading(true);
        fetch(url, { ...options, signal: abortController.signal })
            .then((res) => {
                setStatus(res.status)
                return res.json()
            })
            .then((res) => setData(res))
            .finally(() => setLoading(false))
            .catch((error) => {
                if (error.name === "AbortError") {
                    console.warn("Petición Abortada");
                } else if (error.name != "SyntaxError") {
                    setError(error);
                }
            });

        return () => abortController.abort();
    }, [url, options]);

    return { data, status, loading, error };
};