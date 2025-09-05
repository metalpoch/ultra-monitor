export const isIpv4 = (ip) => {
  if (!/^(\d{1,3}\.){3}\d{1,3}$/.test(ip)) {
    return false;
  }

  const parts = ip.split('.');

  return parts.every(part => {
    const num = Number(part);
    return num >= 0 && num <= 255;
  });
}
