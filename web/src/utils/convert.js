export function traffic(mbps, fracction = 2) {
  const units = [
    { suffix: "E", divisor: 1_000_000_000_000 },
    { suffix: "P", divisor: 1_000_000_000 },
    { suffix: "T", divisor: 1_000_000 },
    { suffix: "G", divisor: 1_000 },
    { suffix: "M", divisor: 1 },
    { suffix: "k", divisor: 0.001 },
    { suffix: "", divisor: 0.000001 },
  ];
  for (const unit of units) {
    if (mbps >= unit.divisor) {
      const value = mbps / unit.divisor;
      return `${value.toFixed(fracction).replace(".", ",")} ${unit.suffix}`;
    }
  }
  return `${mbps.toFixed(fracction).replace(".", ",")} M`;
}

export const qty = (number) =>
  new Intl.NumberFormat("es", {
    notation: "compact",
    compactDisplay: "long",
  }).format(number);

export const dayField = (dataArray) => {
  return dataArray.map((obj) => {
    const dateStr = obj.date || obj.day;
    const date = new Date(dateStr);
    const day = String(date.getUTCDate()).padStart(2, "0");
    const month = String(date.getUTCMonth() + 1).padStart(2, "0");
    const year = String(date.getUTCFullYear()).slice(-2);
    return {
      ...obj,
      day: `${day}/${month}/${year}`,
    };
  });
};

export const getDateRange = (period) => {
  const endDate = new Date();
  const initDate = new Date(endDate);

  if (period === "1d") {
    initDate.setDate(initDate.getDate() - 1);
  } else if (period === "7d") {
    initDate.setDate(initDate.getDate() - 7);
  } else if (period === "1m") {
    initDate.setMonth(initDate.getMonth() - 1);
  } else {
    throw new Error('Período no válido. Use: "1d", "7d" o "1m"');
  }

  const formatToISO = (date) => {
    const pad = (num) => num.toString().padStart(2, "0");

    const year = date.getFullYear();
    const month = pad(date.getMonth() + 1);
    const day = pad(date.getDate());
    const hours = pad(date.getHours());
    const minutes = pad(date.getMinutes());
    const seconds = pad(date.getSeconds());

    return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}-04:00`;
  };

  return {
    initDate: formatToISO(initDate),
    endDate: formatToISO(endDate),
  };
};

export const parseDate = (str) => {
  const [day, month, year] = str.split("/");
  const fullYear = 2000 + parseInt(year, 10);
  return new Date(fullYear, parseInt(month, 10) - 1, parseInt(day, 10));
};

export default { traffic, qty, dayField, getDateRange, parseDate };
