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

export default { traffic, qty, dayField };
