export function removeAccentsAndToUpper(text) {
  const textWithoutAccents = text
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "");
  return textWithoutAccents.toUpperCase().replaceAll("LA GUAIRA", "VARGAS");
}


export function formatSpeed(bps) {
  const suffixes = ["bps", "Kbps", "Mbps", "Gbps", "Tbps", "Pbps", "Ebps"];
  let i = 0;
  let value = bps;

  while (value >= 1000 && i < suffixes.length - 1) {
    value /= 1000;
    i++;
  }

  return `${value.toFixed(2)} ${suffixes[i]}`;
}
