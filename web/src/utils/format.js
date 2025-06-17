export function formatDayField(dataArray) {
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
}
